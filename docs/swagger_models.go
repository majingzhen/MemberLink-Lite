package docs

// 这个文件用于定义 Swagger 文档中使用的数据模型
// 确保所有的 API 响应和请求模型都有完整的文档

import (
	"MemberLink-Lite/common"
	"MemberLink-Lite/models"
	"MemberLink-Lite/services"
)

// AssetInfoResponse 资产信息响应
// @Description 资产信息查询响应
type AssetInfoResponse struct {
	Code    int                `json:"code" example:"200" description:"响应状态码"`
	Message string             `json:"message" example:"获取成功" description:"响应消息"`
	Data    services.AssetInfo `json:"data" description:"资产信息数据"`
	TraceID string             `json:"trace_id" example:"abc123" description:"请求追踪ID"`
}

// BalanceRecordsResponse 余额记录响应
// @Description 余额变动记录查询响应
type BalanceRecordsResponse struct {
	Code    int                   `json:"code" example:"200" description:"响应状态码"`
	Message string                `json:"message" example:"获取成功" description:"响应消息"`
	Data    common.PaginateResult `json:"data" description:"分页数据"`
	TraceID string                `json:"trace_id" example:"abc123" description:"请求追踪ID"`
}

// PointsRecordsResponse 积分记录响应
// @Description 积分变动记录查询响应
type PointsRecordsResponse struct {
	Code    int                   `json:"code" example:"200" description:"响应状态码"`
	Message string                `json:"message" example:"获取成功" description:"响应消息"`
	Data    common.PaginateResult `json:"data" description:"分页数据"`
	TraceID string                `json:"trace_id" example:"abc123" description:"请求追踪ID"`
}

// SuccessResponse 成功响应
// @Description 操作成功的通用响应
type SuccessResponse struct {
	Code    int         `json:"code" example:"200" description:"响应状态码"`
	Message string      `json:"message" example:"操作成功" description:"响应消息"`
	Data    interface{} `json:"data" description:"响应数据"`
	TraceID string      `json:"trace_id" example:"abc123" description:"请求追踪ID"`
}

// ErrorResponse 错误响应
// @Description 操作失败的通用响应
type ErrorResponse struct {
	Code    int         `json:"code" example:"400" description:"错误状态码"`
	Message string      `json:"message" example:"参数错误" description:"错误消息"`
	Data    interface{} `json:"data" description:"错误详情"`
	TraceID string      `json:"trace_id" example:"abc123" description:"请求追踪ID"`
}

// BalanceChangeRequest 余额变动请求（用于 Swagger 文档）
// @Description 余额变动请求参数
type BalanceChangeRequest struct {
	Amount  int64  `json:"amount" binding:"required" example:"1000" description:"变动金额(分为单位)，正数为增加，负数为减少"`
	Type    string `json:"type" binding:"required" example:"recharge" enums:"recharge,consume,refund,reward,deduct" description:"变动类型"`
	Remark  string `json:"remark" example:"用户充值" description:"变动备注说明"`
	OrderNo string `json:"order_no" example:"ORDER20240101001" description:"关联订单号（可选）"`
}

// PointsChangeRequest 积分变动请求（用于 Swagger 文档）
// @Description 积分变动请求参数
type PointsChangeRequest struct {
	Quantity   int64  `json:"quantity" binding:"required" example:"100" description:"变动数量，正数为增加，负数为减少"`
	Type       string `json:"type" binding:"required" example:"obtain" enums:"obtain,use,expire,reward,deduct" description:"变动类型"`
	Remark     string `json:"remark" example:"签到奖励" description:"变动备注说明"`
	OrderNo    string `json:"order_no" example:"ORDER20240101001" description:"关联订单号（可选）"`
	ExpireDays int    `json:"expire_days" example:"365" description:"过期天数，0表示永不过期"`
}

// PaginateResult 分页结果（用于 Swagger 文档）
// @Description 分页查询结果
type PaginateResult struct {
	List     interface{} `json:"list" description:"数据列表"`
	Total    int64       `json:"total" example:"100" description:"总记录数"`
	Page     int         `json:"page" example:"1" description:"当前页码"`
	PageSize int         `json:"page_size" example:"10" description:"每页数量"`
	Pages    int         `json:"pages" example:"10" description:"总页数"`
}

// BalanceRecordExample 余额记录示例（用于 Swagger 文档）
// @Description 余额变动记录示例
type BalanceRecordExample struct {
	ID           uint64 `json:"id" example:"1" description:"记录ID"`
	CreatedAt    string `json:"created_at" example:"2024-01-01T10:00:00Z" description:"创建时间"`
	UpdatedAt    string `json:"updated_at" example:"2024-01-01T10:00:00Z" description:"更新时间"`
	Status       int8   `json:"status" example:"1" description:"状态：1-正常，0-禁用，-1-删除"`
	TenantID     string `json:"tenant_id" example:"default" description:"租户ID"`
	UserID       uint64 `json:"user_id" example:"1" description:"用户ID"`
	Amount       int64  `json:"amount" example:"1000" description:"变动金额(分为单位)"`
	Type         string `json:"type" example:"recharge" description:"变动类型"`
	Remark       string `json:"remark" example:"用户充值" description:"变动备注"`
	BalanceAfter int64  `json:"balance_after" example:"10000" description:"变动后余额(分为单位)"`
	OrderNo      string `json:"order_no" example:"ORDER20240101001" description:"关联订单号"`
}

// PointsRecordExample 积分记录示例（用于 Swagger 文档）
// @Description 积分变动记录示例
type PointsRecordExample struct {
	ID          uint64  `json:"id" example:"1" description:"记录ID"`
	CreatedAt   string  `json:"created_at" example:"2024-01-01T10:00:00Z" description:"创建时间"`
	UpdatedAt   string  `json:"updated_at" example:"2024-01-01T10:00:00Z" description:"更新时间"`
	Status      int8    `json:"status" example:"1" description:"状态：1-正常，0-禁用，-1-删除"`
	TenantID    string  `json:"tenant_id" example:"default" description:"租户ID"`
	UserID      uint64  `json:"user_id" example:"1" description:"用户ID"`
	Quantity    int64   `json:"quantity" example:"100" description:"变动数量"`
	Type        string  `json:"type" example:"obtain" description:"变动类型"`
	Remark      string  `json:"remark" example:"签到奖励" description:"变动备注"`
	PointsAfter int64   `json:"points_after" example:"500" description:"变动后积分"`
	OrderNo     string  `json:"order_no" example:"ORDER20240101001" description:"关联订单号"`
	ExpireTime  *string `json:"expire_time" example:"2024-12-31T23:59:59Z" description:"过期时间"`
}

// AssetTypeEnums 资产变动类型枚举说明
// @Description 资产变动类型说明
type AssetTypeEnums struct {
	BalanceTypes []string `json:"balance_types" example:"[\"recharge\",\"consume\",\"refund\",\"reward\",\"deduct\"]" description:"余额变动类型：recharge-充值，consume-消费，refund-退款，reward-奖励，deduct-扣除"`
	PointsTypes  []string `json:"points_types" example:"[\"obtain\",\"use\",\"expire\",\"reward\",\"deduct\"]" description:"积分变动类型：obtain-获得，use-使用，expire-过期，reward-奖励，deduct-扣除"`
}

// UploadFileResponse 文件上传响应
// @Description 文件上传成功响应
type UploadFileResponse struct {
	Code    int                         `json:"code" example:"200" description:"响应状态码"`
	Message string                      `json:"message" example:"上传成功" description:"响应消息"`
	Data    services.UploadFileResponse `json:"data" description:"上传文件信息"`
	TraceID string                      `json:"trace_id" example:"abc123" description:"请求追踪ID"`
}

// FileInfoResponse 文件信息响应
// @Description 文件信息查询响应
type FileInfoResponse struct {
	Code    int         `json:"code" example:"200" description:"响应状态码"`
	Message string      `json:"message" example:"获取成功" description:"响应消息"`
	Data    models.File `json:"data" description:"文件信息"`
	TraceID string      `json:"trace_id" example:"abc123" description:"请求追踪ID"`
}

// FileListResponse 文件列表响应
// @Description 文件列表查询响应
type FileListResponse struct {
	Code    int                   `json:"code" example:"200" description:"响应状态码"`
	Message string                `json:"message" example:"获取成功" description:"响应消息"`
	Data    common.PaginateResult `json:"data" description:"分页文件列表"`
	TraceID string                `json:"trace_id" example:"abc123" description:"请求追踪ID"`
}

// SignedURLResponse 签名URL响应
// @Description 文件签名URL响应
type SignedURLResponse struct {
	Code    int    `json:"code" example:"200" description:"响应状态码"`
	Message string `json:"message" example:"获取成功" description:"响应消息"`
	Data    string `json:"data" example:"https://example.com/file.jpg?signature=xxx&expires=1234567890" description:"签名URL"`
	TraceID string `json:"trace_id" example:"abc123" description:"请求追踪ID"`
}

// FileExample 文件信息示例（用于 Swagger 文档）
// @Description 文件信息示例
type FileExample struct {
	ID        uint64 `json:"id" example:"1" description:"文件ID"`
	CreatedAt string `json:"created_at" example:"2024-01-01T10:00:00Z" description:"创建时间"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01T10:00:00Z" description:"更新时间"`
	Status    int8   `json:"status" example:"1" description:"状态：1-正常，0-禁用，-1-删除"`
	TenantID  string `json:"tenant_id" example:"default" description:"租户ID"`
	UserID    uint64 `json:"user_id" example:"1" description:"上传用户ID"`
	Filename  string `json:"filename" example:"avatar.jpg" description:"原始文件名"`
	Path      string `json:"path" example:"default/avatar/2024/01/01/1704067200_abc123.jpg" description:"存储路径"`
	URL       string `json:"url" example:"http://localhost:8080/uploads/default/avatar/2024/01/01/1704067200_abc123.jpg" description:"访问URL"`
	Size      int64  `json:"size" example:"1024000" description:"文件大小(字节)"`
	MimeType  string `json:"mime_type" example:"image/jpeg" description:"MIME类型"`
	Hash      string `json:"hash" example:"d41d8cd98f00b204e9800998ecf8427e" description:"文件哈希值"`
	Category  string `json:"category" example:"avatar" description:"文件分类"`
}

// UploadFileResponseExample 文件上传响应示例
// @Description 文件上传响应数据示例
type UploadFileResponseExample struct {
	ID       uint64 `json:"id" example:"1" description:"文件ID"`
	Filename string `json:"filename" example:"avatar.jpg" description:"文件名"`
	URL      string `json:"url" example:"http://localhost:8080/uploads/default/avatar/2024/01/01/1704067200_abc123.jpg" description:"访问URL"`
	Size     int64  `json:"size" example:"1024000" description:"文件大小(字节)"`
	Hash     string `json:"hash" example:"d41d8cd98f00b204e9800998ecf8427e" description:"文件哈希值"`
}

// FileCategoryEnums 文件分类枚举说明
// @Description 文件分类说明
type FileCategoryEnums struct {
	Categories []string `json:"categories" example:"[\"general\",\"avatar\",\"doc\",\"image\"]" description:"文件分类：general-通用文件，avatar-头像，doc-文档，image-图片"`
}

// FileTypeEnums 文件类型枚举说明
// @Description 支持的文件类型说明
type FileTypeEnums struct {
	ImageTypes []string `json:"image_types" example:"[\"image/jpeg\",\"image/png\",\"image/gif\",\"image/webp\"]" description:"支持的图片类型"`
	ImageExts  []string `json:"image_extensions" example:"[\".jpg\",\".jpeg\",\".png\",\".gif\",\".webp\"]" description:"支持的图片扩展名"`
	AvatarExts []string `json:"avatar_extensions" example:"[\".jpg\",\".jpeg\",\".png\"]" description:"头像支持的扩展名"`
}

// FileSizeLimits 文件大小限制说明
// @Description 文件大小限制说明
type FileSizeLimits struct {
	MaxAvatarSize  string `json:"max_avatar_size" example:"5MB" description:"头像文件最大大小"`
	MaxImageSize   string `json:"max_image_size" example:"10MB" description:"图片文件最大大小"`
	MaxGeneralSize string `json:"max_general_size" example:"50MB" description:"通用文件最大大小"`
}
