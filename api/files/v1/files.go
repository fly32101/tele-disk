package v1

import "github.com/gogf/gf/v2/frame/g"

// UploadReq 上传文件
// swagger:route POST /api/files/upload Files upload
//
// Responses:
//
//	default: UploadRes
type UploadReq struct {
	g.Meta `path:"/api/files/upload" tags:"Files" method:"post" summary:"上传文件"`
}

type UploadRes struct {
	Id        int64  `json:"id"`
	Md5       string `json:"md5,omitempty"`
	FileName  string `json:"file_name"`
	SizeBytes int64  `json:"size_bytes"`
	MimeType  string `json:"mime_type"`
	ProxyUrl  string `json:"proxy_url"`
}

// ListReq 获取文件列表
// swagger:route GET /api/files Files list
//
// Responses:
//
//	default: ListRes
type ListReq struct {
	g.Meta   `path:"/api/files" tags:"Files" method:"get" summary:"文件列表"`
	Page     int `json:"page" in:"query" d:"1"`
	PageSize int `json:"page_size" in:"query" d:"20"`
}

type FileItem struct {
	Id        int64  `json:"id"`
	Md5       string `json:"md5,omitempty"`
	FileName  string `json:"file_name"`
	SizeBytes int64  `json:"size_bytes"`
	MimeType  string `json:"mime_type"`
	CreatedAt string `json:"created_at"`
	ProxyUrl  string `json:"proxy_url"`
}

type ListRes struct {
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Files []FileItem `json:"files"`
}

// StatsReq 获取用户文件统计
// swagger:route GET /api/files/stats Files stats
//
// Responses:
//
//	default: StatsRes
type StatsReq struct {
	g.Meta `path:"/api/files/stats" tags:"Files" method:"get" summary:"用户文件统计"`
}

type StatsRes struct {
	UsedBytes int64 `json:"used_bytes"`
	FileCount int   `json:"file_count"`
	Quota     int64 `json:"quota_bytes"`
}

// ProxyReq 代理下载（通过 md5）
// swagger:route GET /api/files/proxy/{md5} Files proxy
//
// Responses:
//
//	default: string
type ProxyReq struct {
	g.Meta `path:"/api/files/proxy/{md5}" tags:"Files" method:"get" summary:"代理下载（按md5）"`
	Md5    string `json:"md5" in:"path" v:"required"` // 按 md5 查找文件
}

type ProxyRes struct {
	g.Meta `mime:"application/octet-stream"`
}
