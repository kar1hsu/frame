package errcode

const (
	Success         = 0
	ErrServer       = 10000
	ErrParam        = 10001
	ErrNotFound     = 10002
	ErrUnauthorized = 10003
	ErrForbidden    = 10004

	ErrUserExist     = 20001
	ErrUserNotFound  = 20002
	ErrPasswordWrong = 20003
	ErrUserDisabled  = 20004

	ErrTokenInvalid = 30001
	ErrTokenExpired = 30002

	ErrRoleExist    = 40001
	ErrRoleNotFound = 40002

	ErrMenuNotFound = 50001
)

var codeMessages = map[int]string{
	Success:         "操作成功",
	ErrServer:       "服务器内部错误",
	ErrParam:        "参数错误",
	ErrNotFound:     "资源不存在",
	ErrUnauthorized: "未登录或登录已过期",
	ErrForbidden:    "无访问权限",

	ErrUserExist:     "用户名已存在",
	ErrUserNotFound:  "用户不存在",
	ErrPasswordWrong: "密码错误",
	ErrUserDisabled:  "用户已被禁用",

	ErrTokenInvalid: "Token 无效",
	ErrTokenExpired: "Token 已过期",

	ErrRoleExist:    "角色编码已存在",
	ErrRoleNotFound: "角色不存在",

	ErrMenuNotFound: "菜单不存在",
}

func Message(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
