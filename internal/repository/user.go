package repository

import (
	"frame/internal/app"
	"frame/internal/model"
	"gorm.io/gorm"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (d *UserRepo) db() *gorm.DB {
	return app.DB
}

func (d *UserRepo) Create(user *model.SysUser) error {
	return d.db().Create(user).Error
}

func (d *UserRepo) GetByID(id uint) (*model.SysUser, error) {
	var user model.SysUser
	err := d.db().Preload("Roles").First(&user, id).Error
	return &user, err
}

func (d *UserRepo) GetByUsername(username string) (*model.SysUser, error) {
	var user model.SysUser
	err := d.db().Preload("Roles").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (d *UserRepo) Update(user *model.SysUser) error {
	return d.db().Save(user).Error
}

func (d *UserRepo) Delete(id uint) error {
	return d.db().Delete(&model.SysUser{}, id).Error
}

func (d *UserRepo) List(page, pageSize int) ([]model.SysUser, int64, error) {
	var users []model.SysUser
	var total int64

	db := d.db().Model(&model.SysUser{})
	db.Count(&total)
	err := db.Preload("Roles").Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id DESC").Find(&users).Error
	return users, total, err
}

func (d *UserRepo) SetRoles(userID uint, roleIDs []uint) error {
	user := &model.SysUser{BaseModel: model.BaseModel{ID: userID}}
	var roles []model.SysRole
	for _, id := range roleIDs {
		roles = append(roles, model.SysRole{BaseModel: model.BaseModel{ID: id}})
	}
	return d.db().Model(user).Association("Roles").Replace(roles)
}
