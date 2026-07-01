package model

// SysOperationLog is an audit record of a single admin mutation (or a
// login/logout event). Retention and "clear" purge rows with Unscoped (hard
// delete) so the table actually shrinks; a single-record delete uses the soft
// delete from BaseModel.
type SysOperationLog struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	TraceID    string `json:"trace_id" gorm:"size:32;comment:链路ID，关联应用日志"`
	UserID     uint   `json:"user_id" gorm:"index;comment:操作人ID"`
	Username   string `json:"username" gorm:"size:64;comment:操作人用户名(冗余存储)"`
	RoleCodes  string `json:"role_codes" gorm:"size:255;comment:操作人角色快照(逗号分隔)"`
	Module     string `json:"module" gorm:"size:64;index;comment:所属模块(取自sys_api.group)"`
	Action     string `json:"action" gorm:"size:128;comment:操作描述(取自sys_api.description)"`
	Method     string `json:"method" gorm:"size:16;comment:HTTP方法"`
	Route      string `json:"route" gorm:"size:255;comment:路由模板"`
	Path       string `json:"path" gorm:"size:255;comment:实际请求路径"`
	TargetID   string `json:"target_id" gorm:"size:64;comment:目标资源ID"`
	ReqParams  string `json:"req_params" gorm:"type:text;comment:请求参数(脱敏+截断)"`
	RespParams string `json:"resp_params" gorm:"type:text;comment:响应参数(脱敏+截断)"`
	Status     int    `json:"status" gorm:"comment:HTTP状态码"`
	BizCode    int    `json:"biz_code" gorm:"comment:业务返回码"`
	Success    bool   `json:"success" gorm:"index;comment:是否成功"`
	ErrorMsg   string `json:"error_msg" gorm:"type:text;comment:失败信息"`
	ClientIP   string `json:"client_ip" gorm:"size:64;index;comment:客户端IP"`
	UserAgent  string `json:"user_agent" gorm:"size:255;comment:User-Agent"`
	LatencyMs  int64  `json:"latency_ms" gorm:"comment:耗时(毫秒)"`
	BaseModel
}

func (SysOperationLog) TableName() string {
	return "sys_operation_log"
}
