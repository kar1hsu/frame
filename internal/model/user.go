package model

type SysUser struct {
	BaseModel
	Username string    `json:"username" gorm:"uniqueIndex;size:64;not null"`
	Password string    `json:"-" gorm:"size:128;not null"`
	Nickname string    `json:"nickname" gorm:"size:64"`
	Avatar   string    `json:"avatar" gorm:"size:255"`
	Email    string    `json:"email" gorm:"size:128"`
	Phone    string    `json:"phone" gorm:"size:20"`
	Status   int8      `json:"status" gorm:"default:1;comment:1-正常 0-禁用"`
	Roles    []SysRole `json:"roles" gorm:"many2many:sys_user_role;"`
}

func (SysUser) TableName() string {
	return "sys_user"
}
