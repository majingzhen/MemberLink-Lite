package models

import (
	"fmt"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

// File 文件模型
type File struct {
	BaseModel
	UserID   uint64 `json:"user_id" gorm:"index;comment:上传用户ID"`
	Filename string `json:"filename" gorm:"size:255;not null;comment:原始文件名"`
	Path     string `json:"path" gorm:"size:500;not null;comment:存储路径"`
	URL      string `json:"url" gorm:"size:500;not null;comment:访问URL"`
	Size     int64  `json:"size" gorm:"not null;comment:文件大小(字节)"`
	MimeType string `json:"mime_type" gorm:"size:100;comment:MIME类型"`
	Hash     string `json:"hash" gorm:"size:64;comment:文件哈希值"`
	Category string `json:"category" gorm:"size:50;default:'general';comment:文件分类"`
	User     *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// FileStatus 文件状态常量
const (
	FileStatusUploading = 0 // 上传中
	FileStatusNormal    = 1 // 正常
	FileStatusDeleted   = 2 // 已删除
	FileStatusAuditing  = 3 // 审核中
	FileStatusRejected  = 4 // 审核拒绝
)

// FileCategory 文件分类常量
const (
	FileCategoryGeneral = "general" // 通用文件
	FileCategoryAvatar  = "avatar"  // 头像
	FileCategoryDoc     = "doc"     // 文档
	FileCategoryImage   = "image"   // 图片
)

// 支持的图片格式
var SupportedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// 支持的文件扩展名
var SupportedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// 文件大小限制（字节）
const (
	MaxAvatarSize  = 5 * 1024 * 1024  // 5MB
	MaxImageSize   = 10 * 1024 * 1024 // 10MB
	MaxGeneralSize = 50 * 1024 * 1024 // 50MB
)

// TableName 指定表名
func (File) TableName() string {
	return "files"
}

// BeforeCreate GORM钩子：创建前处理
func (f *File) BeforeCreate(tx *gorm.DB) error {
	// 调用基础模型的BeforeCreate
	if err := f.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// 设置默认分类
	if f.Category == "" {
		f.Category = FileCategoryGeneral
	}

	// 根据文件名推断MIME类型
	if f.MimeType == "" {
		f.MimeType = mime.TypeByExtension(filepath.Ext(f.Filename))
	}

	return nil
}

// IsImage 检查是否为图片文件
func (f *File) IsImage() bool {
	return SupportedImageTypes[f.MimeType] ||
		SupportedImageExtensions[strings.ToLower(filepath.Ext(f.Filename))]
}

// IsAvatar 检查是否为头像文件
func (f *File) IsAvatar() bool {
	return f.Category == FileCategoryAvatar
}

// GetSizeString 获取文件大小的可读字符串
func (f *File) GetSizeString() string {
	const unit = 1024
	if f.Size < unit {
		return fmt.Sprintf("%d B", f.Size)
	}
	div, exp := int64(unit), 0
	for n := f.Size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(f.Size)/float64(div), "KMGTPE"[exp])
}

// ValidateImageFile 验证图片文件
func ValidateImageFile(filename string, size int64, mimeType string) error {
	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))
	if !SupportedImageExtensions[ext] {
		return fmt.Errorf("不支持的图片格式，仅支持: jpg, jpeg, png, gif, webp")
	}

	// 检查MIME类型
	if mimeType != "" && !SupportedImageTypes[mimeType] {
		return fmt.Errorf("不支持的图片类型: %s", mimeType)
	}

	// 检查文件大小
	if size > MaxImageSize {
		return fmt.Errorf("图片文件大小不能超过 %s", getSizeString(MaxImageSize))
	}

	return nil
}

// ValidateAvatarFile 验证头像文件
func ValidateAvatarFile(filename string, size int64, mimeType string) error {
	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return fmt.Errorf("头像仅支持 JPG 和 PNG 格式")
	}

	// 检查MIME类型
	if mimeType != "" {
		if mimeType != "image/jpeg" && mimeType != "image/jpg" && mimeType != "image/png" {
			return fmt.Errorf("头像仅支持 JPEG 和 PNG 类型")
		}
	}

	// 检查文件大小
	if size > MaxAvatarSize {
		return fmt.Errorf("头像文件大小不能超过 %s", getSizeString(MaxAvatarSize))
	}

	return nil
}

// ValidateGeneralFile 验证通用文件
func ValidateGeneralFile(filename string, size int64) error {
	// 检查文件大小
	if size > MaxGeneralSize {
		return fmt.Errorf("文件大小不能超过 %s", getSizeString(MaxGeneralSize))
	}

	// 检查文件名
	if len(filename) > 255 {
		return fmt.Errorf("文件名长度不能超过255个字符")
	}

	return nil
}

// getSizeString 获取大小的可读字符串
func getSizeString(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// GenerateFilePath 生成文件存储路径
func GenerateFilePath(tenantID, category, filename string) string {
	if tenantID == "" {
		tenantID = "default"
	}

	// 获取文件扩展名
	ext := filepath.Ext(filename)

	// 生成时间路径
	now := time.Now()
	datePath := now.Format("2006/01/02")

	// 生成唯一文件名
	timestamp := now.Unix()
	uniqueName := fmt.Sprintf("%d_%s%s", timestamp, generateRandomString(8), ext)

	// 组合完整路径: tenant_id/category/2006/01/02/timestamp_random.ext
	return fmt.Sprintf("%s/%s/%s/%s", tenantID, category, datePath, uniqueName)
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// ScopeByUser 按用户查询
func ScopeByUser(userID uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}

// ScopeByCategory 按分类查询
func ScopeByCategory(category string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category = ?", category)
	}
}

// ScopeByMimeType 按MIME类型查询
func ScopeByMimeType(mimeType string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("mime_type = ?", mimeType)
	}
}

// ScopeImages 查询图片文件
func ScopeImages(db *gorm.DB) *gorm.DB {
	return db.Where("mime_type LIKE ?", "image/%")
}
