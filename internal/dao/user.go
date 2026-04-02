package dao

import (
	"github.com/karlhsu/frame/internal/app"
	"github.com/karlhsu/frame/internal/model"
	"gorm.io/gorm"
)

type UserDAO struct{}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

func (d *UserDAO) db() *gorm.DB {
	return app.DB
}

func (d *UserDAO) Create(user *model.SysUser) error {
	return d.db().Create(user).Error
}

func (d *UserDAO) GetByID(id uint) (*model.SysUser, error) {
	var user model.SysUser
	err := d.db().Preload("Roles").First(&user, id).Error
	return &user, err
}

func (d *UserDAO) GetByUsername(username string) (*model.SysUser, error) {
	var user model.SysUser
	err := d.db().Preload("Roles").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (d *UserDAO) Update(user *model.SysUser) error {
	return d.db().Save(user).Error
}

func (d *UserDAO) Delete(id uint) error {
	return d.db().Delete(&model.SysUser{}, id).Error
}

func (d *UserDAO) List(page, pageSize int) ([]model.SysUser, int64, error) {
	var users []model.SysUser
	var total int64

	db := d.db().Model(&model.SysUser{})
	db.Count(&total)
	err := db.Preload("Roles").Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id DESC").Find(&users).Error
	return users, total, err
}

func (d *UserDAO) SetRoles(userID uint, roleIDs []uint) error {
	user := &model.SysUser{BaseModel: model.BaseModel{ID: userID}}
	var roles []model.SysRole
	for _, id := range roleIDs {
		roles = append(roles, model.SysRole{BaseModel: model.BaseModel{ID: id}})
	}
	return d.db().Model(user).Association("Roles").Replace(roles)
}
