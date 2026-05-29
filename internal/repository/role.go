package repository

import (
	"github.com/kar1hsu/frame/internal/app"
	"github.com/kar1hsu/frame/internal/model"
	"gorm.io/gorm"
)

type RoleRepo struct{}

func NewRoleRepo() *RoleRepo {
	return &RoleRepo{}
}

func (d *RoleRepo) db() *gorm.DB {
	return app.DB
}

func (d *RoleRepo) Create(role *model.SysRole) error {
	return d.db().Create(role).Error
}

func (d *RoleRepo) GetByID(id uint) (*model.SysRole, error) {
	var role model.SysRole
	err := d.db().Preload("Menus").First(&role, id).Error
	return &role, err
}

func (d *RoleRepo) GetByCode(code string) (*model.SysRole, error) {
	var role model.SysRole
	err := d.db().Where("code = ?", code).First(&role).Error
	return &role, err
}

func (d *RoleRepo) Update(role *model.SysRole) error {
	return d.db().Save(role).Error
}

func (d *RoleRepo) Delete(id uint) error {
	return d.db().Select("Menus").Delete(&model.SysRole{BaseModel: model.BaseModel{ID: id}}).Error
}

func (d *RoleRepo) List(page, pageSize int) ([]model.SysRole, int64, error) {
	var roles []model.SysRole
	var total int64

	db := d.db().Model(&model.SysRole{})
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("sort ASC, id ASC").Find(&roles).Error
	return roles, total, err
}

func (d *RoleRepo) ListAll() ([]model.SysRole, error) {
	var roles []model.SysRole
	err := d.db().Where("status = 1").Order("sort ASC").Find(&roles).Error
	return roles, err
}

func (d *RoleRepo) SetMenus(roleID uint, menuIDs []uint) error {
	role := &model.SysRole{BaseModel: model.BaseModel{ID: roleID}}
	var menus []model.SysMenu
	for _, id := range menuIDs {
		menus = append(menus, model.SysMenu{BaseModel: model.BaseModel{ID: id}})
	}
	return d.db().Model(role).Association("Menus").Replace(menus)
}

func (d *RoleRepo) GetMenusByRoleID(roleID uint) ([]model.SysMenu, error) {
	role := &model.SysRole{BaseModel: model.BaseModel{ID: roleID}}
	var menus []model.SysMenu
	err := d.db().Model(role).Association("Menus").Find(&menus)
	return menus, err
}
