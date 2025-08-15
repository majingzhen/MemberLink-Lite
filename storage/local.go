package storage

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalAdapter 本地存储适配器
type LocalAdapter struct {
	basePath string // 本地存储基础路径
	baseURL  string // 访问基础URL
}

// NewLocalAdapter 创建本地存储适配器
func NewLocalAdapter(basePath, baseURL string) *LocalAdapter {
	// 确保基础路径存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		panic(fmt.Sprintf("创建存储目录失败: %v", err))
	}

	// 确保baseURL以/结尾
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return &LocalAdapter{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

// Upload 上传文件到本地存储
func (l *LocalAdapter) Upload(ctx context.Context, path string, reader io.Reader, size int64, contentType string) error {
	// 构建完整的本地文件路径
	fullPath := filepath.Join(l.basePath, path)

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &StorageError{
			Code:    "CREATE_DIR_FAILED",
			Message: "创建目录失败",
			Err:     err,
		}
	}

	// 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return &StorageError{
			Code:    "CREATE_FILE_FAILED",
			Message: "创建文件失败",
			Err:     err,
		}
	}
	defer file.Close()

	// 复制数据
	_, err = io.Copy(file, reader)
	if err != nil {
		// 删除部分上传的文件
		os.Remove(fullPath)
		return &StorageError{
			Code:    "WRITE_FILE_FAILED",
			Message: "写入文件失败",
			Err:     err,
		}
	}

	return nil
}

// Download 从本地存储下载文件
func (l *LocalAdapter) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(l.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, &StorageError{
			Code:    "OPEN_FILE_FAILED",
			Message: "打开文件失败",
			Err:     err,
		}
	}

	return file, nil
}

// Delete 删除本地存储文件
func (l *LocalAdapter) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(l.basePath, path)

	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return &StorageError{
			Code:    "DELETE_FILE_FAILED",
			Message: "删除文件失败",
			Err:     err,
		}
	}

	return nil
}

// Exists 检查文件是否存在
func (l *LocalAdapter) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(l.basePath, path)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, &StorageError{
			Code:    "CHECK_FILE_FAILED",
			Message: "检查文件失败",
			Err:     err,
		}
	}

	return true, nil
}

// GetURL 获取文件访问URL
func (l *LocalAdapter) GetURL(ctx context.Context, path string) (string, error) {
	// 对路径进行URL编码
	encodedPath := url.PathEscape(strings.ReplaceAll(path, "\\", "/"))
	return l.baseURL + encodedPath, nil
}

// GetSignedURL 获取签名URL（本地存储直接返回普通URL）
func (l *LocalAdapter) GetSignedURL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	// 本地存储不需要签名，直接返回普通URL
	return l.GetURL(ctx, path)
}

// GetStorageType 获取存储类型
func (l *LocalAdapter) GetStorageType() string {
	return "local"
}

// GetFileHash 计算文件MD5哈希
func (l *LocalAdapter) GetFileHash(ctx context.Context, path string) (string, error) {
	fullPath := filepath.Join(l.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		return "", &StorageError{
			Code:    "OPEN_FILE_FAILED",
			Message: "打开文件失败",
			Err:     err,
		}
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", &StorageError{
			Code:    "CALCULATE_HASH_FAILED",
			Message: "计算文件哈希失败",
			Err:     err,
		}
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// GetFileSize 获取文件大小
func (l *LocalAdapter) GetFileSize(ctx context.Context, path string) (int64, error) {
	fullPath := filepath.Join(l.basePath, path)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, ErrFileNotFound
		}
		return 0, &StorageError{
			Code:    "GET_FILE_INFO_FAILED",
			Message: "获取文件信息失败",
			Err:     err,
		}
	}

	return info.Size(), nil
}
