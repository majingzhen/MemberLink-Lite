package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// FileValidator 文件验证器
type FileValidator struct {
	MaxSize           int64           // 最大文件大小（字节）
	AllowedExtensions map[string]bool // 允许的文件扩展名
	AllowedMimeTypes  map[string]bool // 允许的MIME类型
}

// NewFileValidator 创建文件验证器
func NewFileValidator() *FileValidator {
	return &FileValidator{
		AllowedExtensions: make(map[string]bool),
		AllowedMimeTypes:  make(map[string]bool),
	}
}

// SetMaxSize 设置最大文件大小
func (v *FileValidator) SetMaxSize(size int64) *FileValidator {
	v.MaxSize = size
	return v
}

// AddAllowedExtension 添加允许的文件扩展名
func (v *FileValidator) AddAllowedExtension(ext string) *FileValidator {
	v.AllowedExtensions[strings.ToLower(ext)] = true
	return v
}

// AddAllowedMimeType 添加允许的MIME类型
func (v *FileValidator) AddAllowedMimeType(mimeType string) *FileValidator {
	v.AllowedMimeTypes[strings.ToLower(mimeType)] = true
	return v
}

// Validate 验证文件
func (v *FileValidator) Validate(header *multipart.FileHeader) error {
	// 检查文件大小
	if v.MaxSize > 0 && header.Size > v.MaxSize {
		return fmt.Errorf("文件大小超过限制，最大允许 %s", FormatFileSize(v.MaxSize))
	}

	// 检查文件扩展名
	if len(v.AllowedExtensions) > 0 {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if !v.AllowedExtensions[ext] {
			return fmt.Errorf("不支持的文件格式: %s", ext)
		}
	}

	// 检查MIME类型
	if len(v.AllowedMimeTypes) > 0 {
		// 从文件头获取MIME类型
		file, err := header.Open()
		if err != nil {
			return fmt.Errorf("无法打开文件: %v", err)
		}
		defer file.Close()

		// 读取文件头部分来检测MIME类型
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("无法读取文件内容: %v", err)
		}

		// 重置文件指针
		file.Seek(0, 0)

		// 检测MIME类型
		mimeType := DetectMimeType(buffer, header.Filename)
		if !v.AllowedMimeTypes[strings.ToLower(mimeType)] {
			return fmt.Errorf("不支持的文件类型: %s", mimeType)
		}
	}

	return nil
}

// CreateImageValidator 创建图片文件验证器
func CreateImageValidator(maxSize int64) *FileValidator {
	return NewFileValidator().
		SetMaxSize(maxSize).
		AddAllowedExtension(".jpg").
		AddAllowedExtension(".jpeg").
		AddAllowedExtension(".png").
		AddAllowedExtension(".gif").
		AddAllowedExtension(".webp").
		AddAllowedMimeType("image/jpeg").
		AddAllowedMimeType("image/jpg").
		AddAllowedMimeType("image/png").
		AddAllowedMimeType("image/gif").
		AddAllowedMimeType("image/webp")
}

// CreateAvatarValidator 创建头像文件验证器
func CreateAvatarValidator(maxSize int64) *FileValidator {
	return NewFileValidator().
		SetMaxSize(maxSize).
		AddAllowedExtension(".jpg").
		AddAllowedExtension(".jpeg").
		AddAllowedExtension(".png").
		AddAllowedMimeType("image/jpeg").
		AddAllowedMimeType("image/jpg").
		AddAllowedMimeType("image/png")
}

// FormatFileSize 格式化文件大小
func FormatFileSize(size int64) string {
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

// DetectMimeType 检测文件MIME类型
func DetectMimeType(data []byte, filename string) string {
	// 基于文件内容检测MIME类型
	if len(data) >= 8 {
		// JPEG
		if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
			return "image/jpeg"
		}
		// PNG
		if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
			return "image/png"
		}
		// GIF
		if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
			return "image/gif"
		}
		// WebP
		if data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
			data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
			return "image/webp"
		}
	}

	// 基于文件扩展名推断
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// CalculateFileHash 计算文件MD5哈希
func CalculateFileHash(reader io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// CalculateFileHashFromMultipart 从multipart文件计算哈希
func CalculateFileHashFromMultipart(header *multipart.FileHeader) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	return CalculateFileHash(file)
}

// IsImageFile 检查是否为图片文件
func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return imageExts[ext]
}

// SanitizeFilename 清理文件名
func SanitizeFilename(filename string) string {
	// 移除路径分隔符和其他危险字符
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, "..", "_")
	filename = strings.ReplaceAll(filename, "<", "_")
	filename = strings.ReplaceAll(filename, ">", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "\"", "_")
	filename = strings.ReplaceAll(filename, "|", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "*", "_")

	// 限制文件名长度
	if len(filename) > 255 {
		ext := filepath.Ext(filename)
		name := filename[:255-len(ext)]
		filename = name + ext
	}

	return filename
}
