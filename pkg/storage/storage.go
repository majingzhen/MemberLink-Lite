package storage

import (
	"context"
	"io"
	"time"
)

// StorageAdapter 存储适配器接口
type StorageAdapter interface {
	// Upload 上传文件
	Upload(ctx context.Context, path string, reader io.Reader, size int64, contentType string) error

	// Download 下载文件
	Download(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete 删除文件
	Delete(ctx context.Context, path string) error

	// Exists 检查文件是否存在
	Exists(ctx context.Context, path string) (bool, error)

	// GetURL 获取文件访问URL
	GetURL(ctx context.Context, path string) (string, error)

	// GetSignedURL 获取签名URL（临时访问）
	GetSignedURL(ctx context.Context, path string, expiry time.Duration) (string, error)

	// GetStorageType 获取存储类型
	GetStorageType() string
}

// UploadResult 上传结果
type UploadResult struct {
	Path        string `json:"path"`         // 存储路径
	URL         string `json:"url"`          // 访问URL
	Size        int64  `json:"size"`         // 文件大小
	ContentType string `json:"content_type"` // 内容类型
	Hash        string `json:"hash"`         // 文件哈希
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type     string                 `json:"type"`      // 存储类型: local, aliyun, tencent
	Config   map[string]interface{} `json:"config"`    // 具体配置
	BasePath string                 `json:"base_path"` // 基础路径
	BaseURL  string                 `json:"base_url"`  // 基础URL
}

// StorageError 存储错误
type StorageError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *StorageError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *StorageError) Unwrap() error {
	return e.Err
}

// 预定义错误
var (
	ErrFileNotFound    = &StorageError{Code: "FILE_NOT_FOUND", Message: "文件不存在"}
	ErrUploadFailed    = &StorageError{Code: "UPLOAD_FAILED", Message: "文件上传失败"}
	ErrDownloadFailed  = &StorageError{Code: "DOWNLOAD_FAILED", Message: "文件下载失败"}
	ErrDeleteFailed    = &StorageError{Code: "DELETE_FAILED", Message: "文件删除失败"}
	ErrInvalidConfig   = &StorageError{Code: "INVALID_CONFIG", Message: "存储配置无效"}
	ErrUnsupportedType = &StorageError{Code: "UNSUPPORTED_TYPE", Message: "不支持的存储类型"}
)
