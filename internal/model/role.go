package model

type SysRole struct {
	BaseModel
	Name   string    `json:"name" gorm:"size:64;not null"`
	Code   string    `json:"code" gorm:"uniqueIndex;size:64;not null"`
	Sort   int       `json:"sort" gorm:"default:0"`
	Status int8      `json:"status" gorm:"default:1;comment:1-正常 0-禁用"`
	Remark string    `json:"remark" gorm:"size:255"`
	Menus  []SysMenu `json:"menus,omitempty" gorm:"many2many:sys_role_menu;"`
}

func (SysRole) TableName() string {
	return "sys_role"
}
