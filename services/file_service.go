package services

import (
	"MemberLink-Lite/common"
	"MemberLink-Lite/models"
	"MemberLink-Lite/storage"
	"MemberLink-Lite/utils"
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"gorm.io/gorm"
)

// FileService 文件服务
type FileService struct {
	db      *gorm.DB
	storage storage.StorageAdapter
}

// NewFileService 创建文件服务
func NewFileService(db *gorm.DB) *FileService {
	return &FileService{
		db:      db,
		storage: storage.GetCurrentAdapter(),
	}
}

// UploadFileRequest 文件上传请求
type UploadFileRequest struct {
	UserID   uint64                `json:"user_id"`
	Category string                `json:"category"`
	File     *multipart.FileHeader `json:"file"`
	TenantID string                `json:"tenant_id"`
}

// UploadFileResponse 文件上传响应
type UploadFileResponse struct {
	ID       uint64 `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	Hash     string `json:"hash"`
}

// UploadAvatar 上传头像
func (s *FileService) UploadAvatar(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error) {
	// 验证头像文件
	validator := utils.CreateAvatarValidator(models.MaxAvatarSize)
	if err := validator.Validate(req.File); err != nil {
		return nil, common.NewCustomError(common.CodeBadRequest, err.Error())
	}

	// 设置文件分类为头像
	req.Category = models.FileCategoryAvatar

	return s.uploadFile(ctx, req)
}

// UploadImage 上传图片
func (s *FileService) UploadImage(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error) {
	// 验证图片文件
	validator := utils.CreateImageValidator(models.MaxImageSize)
	if err := validator.Validate(req.File); err != nil {
		return nil, common.NewCustomError(common.CodeBadRequest, err.Error())
	}

	// 设置文件分类为图片
	if req.Category == "" {
		req.Category = models.FileCategoryImage
	}

	return s.uploadFile(ctx, req)
}

// UploadGeneral 上传通用文件
func (s *FileService) UploadGeneral(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error) {
	// 验证通用文件
	if err := models.ValidateGeneralFile(req.File.Filename, req.File.Size); err != nil {
		return nil, common.NewCustomError(common.CodeBadRequest, err.Error())
	}

	// 设置文件分类为通用
	if req.Category == "" {
		req.Category = models.FileCategoryGeneral
	}

	return s.uploadFile(ctx, req)
}

// uploadFile 上传文件的通用方法
func (s *FileService) uploadFile(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error) {
	// 设置默认租户ID
	if req.TenantID == "" {
		req.TenantID = "default"
	}

	// 清理文件名
	filename := utils.SanitizeFilename(req.File.Filename)

	// 生成存储路径
	storagePath := models.GenerateFilePath(req.TenantID, req.Category, filename)

	// 打开文件
	file, err := req.File.Open()
	if err != nil {
		return nil, common.ErrOpenFileFailed
	}
	defer file.Close()

	// 计算文件哈希
	hash, err := utils.CalculateFileHashFromMultipart(req.File)
	if err != nil {
		return nil, common.ErrCalculateFileHashFailed
	}

	// 检查文件是否已存在（基于哈希去重）
	var existingFile models.File
	if err := s.db.Where("hash = ? AND tenant_id = ?", hash, req.TenantID).First(&existingFile).Error; err == nil {
		// 文件已存在，返回现有文件信息
		return &UploadFileResponse{
			ID:       existingFile.ID,
			Filename: existingFile.Filename,
			URL:      existingFile.URL,
			Size:     existingFile.Size,
			Hash:     existingFile.Hash,
		}, nil
	}

	// 重新打开文件用于上传
	file, err = req.File.Open()
	if err != nil {
		return nil, common.ErrOpenFileFailed
	}
	defer file.Close()

	// 检测MIME类型
	buffer := make([]byte, 512)
	file.Read(buffer)
	file.Seek(0, 0) // 重置文件指针
	mimeType := utils.DetectMimeType(buffer, filename)

	// 上传到存储
	if err := s.storage.Upload(ctx, storagePath, file, req.File.Size, mimeType); err != nil {
		return nil, common.ErrUploadFailed
	}

	// 获取文件访问URL
	url, err := s.storage.GetURL(ctx, storagePath)
	if err != nil {
		return nil, common.ErrGetFileUrlFailed
	}

	// 创建文件记录
	fileRecord := &models.File{
		UserID:   req.UserID,
		Filename: filename,
		Path:     storagePath,
		URL:      url,
		Size:     req.File.Size,
		MimeType: mimeType,
		Hash:     hash,
		Category: req.Category,
	}
	fileRecord.TenantID = req.TenantID
	fileRecord.Status = models.StatusActive

	// 保存到数据库
	if err := s.db.Create(fileRecord).Error; err != nil {
		// 如果数据库保存失败，尝试删除已上传的文件
		s.storage.Delete(ctx, storagePath)
		return nil, common.ErrSaveFileFailed
	}

	return &UploadFileResponse{
		ID:       fileRecord.ID,
		Filename: fileRecord.Filename,
		URL:      fileRecord.URL,
		Size:     fileRecord.Size,
		Hash:     fileRecord.Hash,
	}, nil
}

// GetFileByID 根据ID获取文件信息
func (s *FileService) GetFileByID(ctx context.Context, fileID uint64, tenantID string) (*models.File, error) {
	if tenantID == "" {
		tenantID = "default"
	}

	var file models.File
	err := s.db.Scopes(
		models.ScopeActiveByTenant(tenantID),
	).First(&file, fileID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrFileNotFound
		}
		return nil, common.ErrSearchFileFailed
	}

	return &file, nil
}

// GetSignedURL 获取文件签名URL
func (s *FileService) GetSignedURL(ctx context.Context, fileID uint64, tenantID string, expiry time.Duration) (string, error) {
	// 获取文件信息
	file, err := s.GetFileByID(ctx, fileID, tenantID)
	if err != nil {
		return "", err
	}

	// 生成签名URL
	signedURL, err := s.storage.GetSignedURL(ctx, file.Path, expiry)
	if err != nil {
		return "", common.ErrGetSignedURLFailed
	}

	return signedURL, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(ctx context.Context, fileID uint64, userID uint64, tenantID string) error {
	if tenantID == "" {
		tenantID = "default"
	}

	// 查找文件
	var file models.File
	err := s.db.Scopes(
		models.ScopeActiveByTenant(tenantID),
		models.ScopeByUser(userID),
	).First(&file, fileID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return common.ErrFileNotFound
		}
		return common.ErrSearchFileFailed
	}

	// 软删除文件记录
	file.SetDeleted()
	if err := s.db.Save(&file).Error; err != nil {
		return common.ErrDeleteFileFailed
	}

	// 从存储中删除文件（异步执行，不影响响应）
	go func() {
		if err := s.storage.Delete(context.Background(), file.Path); err != nil {
			// 记录日志，但不影响主流程
			fmt.Printf("删除存储文件失败: %v\n", err)
		}
	}()

	return nil
}

// GetUserFiles 获取用户文件列表
func (s *FileService) GetUserFiles(ctx context.Context, userID uint64, tenantID string, req *common.PageRequest) (*common.PageResponse, error) {
	if tenantID == "" {
		tenantID = "default"
	}

	var files []models.File
	var total int64

	// 构建查询
	query := s.db.Model(&models.File{}).Scopes(
		models.ScopeActiveByTenant(tenantID),
		models.ScopeByUser(userID),
	)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, common.ErrSearchFileCountFailed
	}

	// 分页查询
	if err := query.Scopes(common.Paginate(req.Page, req.PageSize)).
		Order("created_at DESC").
		Find(&files).Error; err != nil {
		return nil, common.ErrSearchFileListFailed
	}

	return &common.PageResponse{
		List:     files,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Pages:    int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
	}, nil
}

// GetUserFilesByCategory 根据分类获取用户文件
func (s *FileService) GetUserFilesByCategory(ctx context.Context, userID uint64, category string, tenantID string, req *common.PageRequest) (*common.PageResponse, error) {
	if tenantID == "" {
		tenantID = "default"
	}

	var files []models.File
	var total int64

	// 构建查询
	query := s.db.Model(&models.File{}).Scopes(
		models.ScopeActiveByTenant(tenantID),
		models.ScopeByUser(userID),
		models.ScopeByCategory(category),
	)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, common.ErrSearchFileCountFailed
	}

	// 分页查询
	if err := query.Scopes(common.Paginate(req.Page, req.PageSize)).
		Order("created_at DESC").
		Find(&files).Error; err != nil {
		return nil, common.ErrSearchFileListFailed
	}

	return &common.PageResponse{
		List:     files,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Pages:    int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
	}, nil
}

// UpdateUserAvatar 更新用户头像
func (s *FileService) UpdateUserAvatar(ctx context.Context, userID uint64, avatarURL string, tenantID string) error {
	if tenantID == "" {
		tenantID = "default"
	}

	// 更新用户头像URL
	err := s.db.Model(&models.User{}).
		Where("id = ? AND tenant_id = ?", userID, tenantID).
		Update("avatar", avatarURL).Error

	if err != nil {
		return common.ErrEditUserAvatar
	}

	return nil
}
