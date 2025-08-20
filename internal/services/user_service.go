package services

import (
	"context"
	"fmt"
	"member-link-lite/internal/database"
	"member-link-lite/internal/models"
	"member-link-lite/pkg/common"
	"member-link-lite/pkg/utils"
	"mime/multipart"
	"regexp"
	"time"

	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	// 注册
	Register(ctx context.Context, req *RegisterRequest) (*models.User, error)
	// 登录
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	// 根据用户名查找用户
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	// 根据手机号查找用户
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	// 根据邮箱查找用户
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	// 根据微信OpenID查找用户
	GetByWeChatOpenID(ctx context.Context, openID string) (*models.User, error)
	// 根据ID查找用户
	GetByID(ctx context.Context, id uint64) (*models.User, error)
	// 检查用户名是否存在
	IsUsernameExists(ctx context.Context, username string) (bool, error)
	// 检查手机号是否存在
	IsPhoneExists(ctx context.Context, phone string) (bool, error)
	// 检查邮箱是否存在
	IsEmailExists(ctx context.Context, email string) (bool, error)
	// 检查微信OpenID是否存在
	IsWeChatOpenIDExists(ctx context.Context, openID string) (bool, error)
	// 刷新令牌
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	// 更新最后登录信息
	UpdateLastLogin(ctx context.Context, userID uint64, ip string) error
	// 更新用户信息
	UpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileRequest) error
	// 修改密码
	ChangePassword(ctx context.Context, userID uint64, req *ChangePasswordRequest) error
	// 上传头像
	UploadAvatar(ctx context.Context, userID uint64, file *multipart.FileHeader) (string, error)
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username      string `json:"username" binding:"required,min=3,max=20" example:"testuser"`
	Password      string `json:"password" binding:"required,min=6,max=20" example:"password123"`
	Phone         string `json:"phone" binding:"required" example:"13800138000"`
	Email         string `json:"email" binding:"required,email" example:"test@example.com"`
	Nickname      string `json:"nickname" binding:"max=20" example:"测试用户"`
	WeChatOpenID  string `json:"wechat_openid" example:"wx_openid_123"`
	WeChatUnionID string `json:"wechat_unionid" example:"wx_unionid_123"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"testuser"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	User   *models.User   `json:"user"`
	Tokens *TokenResponse `json:"tokens"`
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname      string `json:"nickname" binding:"max=20" example:"新昵称"`
	Email         string `json:"email" binding:"omitempty,email" example:"newemail@example.com"`
	Phone         string `json:"phone" binding:"omitempty" example:"13800138001"`
	Avatar        string `json:"avatar" binding:"omitempty" example:"http://example.com/avatar.jpg"`
	WeChatOpenID  string `json:"wechat_openid" example:"wx_openid_123"`
	WeChatUnionID string `json:"wechat_unionid" example:"wx_unionid_123"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"oldpassword123"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20" example:"newpassword123"`
}

// userServiceImpl 用户服务实现
type userServiceImpl struct {
	db         *gorm.DB
	jwtService JWTService
}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userServiceImpl{
		db:         database.GetDB(),
		jwtService: NewJWTService(),
	}
}

// Register 用户注册
func (s *userServiceImpl) Register(ctx context.Context, req *RegisterRequest) (*models.User, error) {
	// 验证请求参数
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	tenantID := database.GetTenantIDFromContext(ctx)

	// 检查用户名是否已存在（限定租户）
	exists, err := s.IsUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrUserExists
	}

	// 检查手机号是否已存在（限定租户）
	exists, err = s.IsPhoneExists(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrPhoneExists
	}

	// 检查邮箱是否已存在（限定租户）
	exists, err = s.IsEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrEmailExists
	}

	// 检查微信OpenID是否已存在（如果提供）
	if req.WeChatOpenID != "" {
		exists, err = s.IsWeChatOpenIDExists(ctx, req.WeChatOpenID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, common.ErrUserExists
		}
	}

	// 创建用户
	user := &models.User{
		Username:      req.Username,
		Phone:         req.Phone,
		Email:         req.Email,
		Nickname:      req.Nickname,
		WeChatOpenID:  req.WeChatOpenID,
		WeChatUnionID: req.WeChatUnionID,
	}
	// 设置租户ID
	user.TenantID = tenantID

	// 设置默认昵称
	if user.Nickname == "" {
		user.Nickname = user.Username
	}

	// 加密密码
	if err := user.HashPassword(req.Password); err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 保存到数据库
	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user, nil
}

// Login 用户登录
func (s *userServiceImpl) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	user, err := s.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, common.ErrInvalidPassword // 不暴露用户是否存在的信息
	}

	// 检查用户状态
	if !user.IsActive() {
		if user.Status == models.UserStatusDisabled {
			return nil, common.ErrUserDisabled
		}
		return nil, common.ErrInvalidPassword
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, common.ErrInvalidPassword
	}

	// 生成令牌
	tokens, err := GenerateTokenPair(s.jwtService, user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	return &LoginResponse{
		User:   user,
		Tokens: tokens,
	}, nil
}

// RefreshToken 刷新令牌
func (s *userServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// 验证刷新令牌
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 检查令牌类型
	if claims.Type != "refresh" {
		return nil, common.ErrInvalidToken
	}

	// 验证用户是否仍然有效
	user, err := s.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, common.ErrInvalidToken
	}

	if !user.IsActive() {
		return nil, common.ErrUserDisabled
	}

	// 生成新的令牌对
	tokens, err := GenerateTokenPair(s.jwtService, user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	return tokens, nil
}

// UpdateLastLogin 更新最后登录信息
func (s *userServiceImpl) UpdateLastLogin(ctx context.Context, userID uint64, ip string) error {
	return s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_ip":   ip,
			"last_time": time.Now(),
		}).Error
}

// UpdateProfile 更新用户信息
func (s *userServiceImpl) UpdateProfile(ctx context.Context, userID uint64, req *UpdateProfileRequest) error {
	// 验证请求参数
	if err := s.validateUpdateProfileRequest(ctx, userID, req); err != nil {
		return err
	}

	tenantID := database.GetTenantIDFromContext(ctx)

	// 构建更新数据
	updates := make(map[string]interface{})

	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}

	if req.Email != "" {
		updates["email"] = req.Email
	}

	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if req.WeChatOpenID != "" {
		updates["wechat_openid"] = req.WeChatOpenID
	}

	if req.WeChatUnionID != "" {
		updates["wechat_unionid"] = req.WeChatUnionID
	}

	// 如果没有要更新的字段，直接返回
	if len(updates) == 0 {
		return nil
	}

	// 执行更新（限定租户）
	result := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ? AND status = ? AND tenant_id = ?", userID, models.StatusActive, tenantID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("更新用户信息失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return common.ErrUserNotFound
	}

	return nil
}

// ChangePassword 修改密码
func (s *userServiceImpl) ChangePassword(ctx context.Context, userID uint64, req *ChangePasswordRequest) error {
	// 验证新密码强度
	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		return err
	}

	// 获取用户信息
	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		return common.ErrInvalidPassword
	}

	// 加密新密码
	newUser := &models.User{}
	if err := newUser.HashPassword(req.NewPassword); err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	tenantID := database.GetTenantIDFromContext(ctx)

	// 更新密码（限定租户）
	result := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ? AND status = ? AND tenant_id = ?", userID, models.StatusActive, tenantID).
		Update("password", newUser.Password)

	if result.Error != nil {
		return fmt.Errorf("修改密码失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return common.ErrUserNotFound
	}

	return nil
}

// UploadAvatar 上传头像
func (s *userServiceImpl) UploadAvatar(ctx context.Context, userID uint64, file *multipart.FileHeader) (string, error) {
	// 验证文件类型
	if err := validateAvatarFile(file); err != nil {
		return "", err
	}

	// 验证用户是否存在
	_, err := s.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	// TODO: 实现文件上传到OSS的逻辑
	// 这里暂时返回一个示例URL，实际应该上传到OSS并返回真实URL
	avatarURL := fmt.Sprintf("https://example.com/avatars/%d_%d.jpg", userID, time.Now().Unix())

	tenantID := database.GetTenantIDFromContext(ctx)

	// 更新用户头像URL（限定租户）
	result := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ? AND status = ? AND tenant_id = ?", userID, models.StatusActive, tenantID).
		Update("avatar", avatarURL)

	if result.Error != nil {
		return "", fmt.Errorf("更新头像失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return "", common.ErrUserNotFound
	}

	// 记录文件信息到数据库（如果有文件表的话）
	// TODO: 实现文件记录逻辑

	return avatarURL, nil
}

// validateAvatarFile 验证头像文件
func validateAvatarFile(file *multipart.FileHeader) error {
	// 检查文件大小（5MB限制）
	const maxSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxSize {
		return &common.CustomError{
			Code:    413,
			Message: "文件大小不能超过5MB",
		}
	}

	// 检查文件类型
	contentType := file.Header.Get("Content-Type")
	switch contentType {
	case "image/jpeg", "image/jpg", "image/png":
		// 支持的格式
	default:
		return &common.CustomError{
			Code:    415,
			Message: "仅支持jpg、png格式的图片",
		}
	}

	return nil
}

// validateUpdateProfileRequest 验证更新用户信息请求
func (s *userServiceImpl) validateUpdateProfileRequest(ctx context.Context, userID uint64, req *UpdateProfileRequest) error {
	// 验证邮箱格式和唯一性
	if req.Email != "" {
		if err := validateEmail(req.Email); err != nil {
			return err
		}

		// 检查邮箱是否被其他用户使用（限定租户）
		var count int64
		err := s.db.WithContext(ctx).
			Model(&models.User{}).
			Where("email = ? AND id != ? AND tenant_id = ?", req.Email, userID, database.GetTenantIDFromContext(ctx)).
			Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return common.ErrEmailExists
		}
	}

	// 验证手机号格式和唯一性
	if req.Phone != "" {
		if err := validatePhone(req.Phone); err != nil {
			return err
		}

		// 检查手机号是否被其他用户使用（限定租户）
		var count int64
		err := s.db.WithContext(ctx).
			Model(&models.User{}).
			Where("phone = ? AND id != ? AND tenant_id = ?", req.Phone, userID, database.GetTenantIDFromContext(ctx)).
			Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return common.ErrPhoneExists
		}
	}

	return nil
}

// GetByUsername 根据用户名查找用户
func (s *userServiceImpl) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).
		Scopes(models.ScopeActive).
		Where("username = ? AND tenant_id = ?", username, database.GetTenantIDFromContext(ctx)).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByPhone 根据手机号查找用户
func (s *userServiceImpl) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).
		Scopes(models.ScopeActive).
		Where("phone = ? AND tenant_id = ?", phone, database.GetTenantIDFromContext(ctx)).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail 根据邮箱查找用户
func (s *userServiceImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).
		Scopes(models.ScopeActive).
		Where("email = ? AND tenant_id = ?", email, database.GetTenantIDFromContext(ctx)).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByID 根据ID查找用户
func (s *userServiceImpl) GetByID(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).
		Scopes(models.ScopeActive).
		Where("id = ? AND tenant_id = ?", id, database.GetTenantIDFromContext(ctx)).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// IsUsernameExists 检查用户名是否存在
func (s *userServiceImpl) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ? AND tenant_id = ?", username, database.GetTenantIDFromContext(ctx)).
		Count(&count).Error

	return count > 0, err
}

// IsPhoneExists 检查手机号是否存在
func (s *userServiceImpl) IsPhoneExists(ctx context.Context, phone string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("phone = ? AND tenant_id = ?", phone, database.GetTenantIDFromContext(ctx)).
		Count(&count).Error

	return count > 0, err
}

// IsEmailExists 检查邮箱是否存在
func (s *userServiceImpl) IsEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ? AND tenant_id = ?", email, database.GetTenantIDFromContext(ctx)).
		Count(&count).Error

	return count > 0, err
}

// validateRegisterRequest 验证注册请求
func (s *userServiceImpl) validateRegisterRequest(req *RegisterRequest) error {
	// 验证用户名格式
	if err := validateUsername(req.Username); err != nil {
		return err
	}

	// 验证密码强度
	if err := utils.ValidatePassword(req.Password); err != nil {
		return err
	}

	// 验证手机号格式
	if err := validatePhone(req.Phone); err != nil {
		return err
	}

	// 验证邮箱格式
	if err := validateEmail(req.Email); err != nil {
		return err
	}

	return nil
}

// validateUsername 验证用户名格式
func validateUsername(username string) error {
	if len(username) < 3 {
		return fmt.Errorf("用户名长度不能少于3位")
	}
	if len(username) > 20 {
		return fmt.Errorf("用户名长度不能超过20位")
	}

	// 用户名只能包含字母、数字和下划线
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return fmt.Errorf("用户名只能包含字母、数字和下划线")
	}

	return nil
}

// validatePhone 验证手机号格式
func validatePhone(phone string) error {
	if phone == "" {
		return fmt.Errorf("手机号不能为空")
	}

	// 简单的中国手机号验证
	matched, _ := regexp.MatchString("^1[3-9]\\d{9}$", phone)
	if !matched {
		return fmt.Errorf("手机号格式不正确")
	}

	return nil
}

// validateEmail 验证邮箱格式
func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("邮箱不能为空")
	}

	// 简单的邮箱格式验证
	matched, _ := regexp.MatchString("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", email)
	if !matched {
		return fmt.Errorf("邮箱格式不正确")
	}

	return nil
}

// GetByWeChatOpenID 根据微信OpenID查找用户
func (s *userServiceImpl) GetByWeChatOpenID(ctx context.Context, openID string) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).
		Where("wechat_openid = ? AND tenant_id = ?", openID, database.GetTenantIDFromContext(ctx)).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// IsWeChatOpenIDExists 检查微信OpenID是否存在
func (s *userServiceImpl) IsWeChatOpenIDExists(ctx context.Context, openID string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("wechat_openid = ? AND tenant_id = ?", openID, database.GetTenantIDFromContext(ctx)).
		Count(&count).Error

	return count > 0, err
}
