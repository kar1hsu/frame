package repository

import (
	"frame/internal/app"
	"frame/internal/model"
	"gorm.io/gorm"
)

type ApiRepo struct{}

func NewApiRepo() *ApiRepo {
	return &ApiRepo{}
}

func (d *ApiRepo) db() *gorm.DB {
	return app.DB
}

func (d *ApiRepo) Create(api *model.SysAPI) error {
	return d.db().Create(api).Error
}

func (d *ApiRepo) GetByID(id uint) (*model.SysAPI, error) {
	var api model.SysAPI
	err := d.db().First(&api, id).Error
	return &api, err
}

func (d *ApiRepo) Update(api *model.SysAPI) error {
	return d.db().Save(api).Error
}

func (d *ApiRepo) Delete(id uint) error {
	return d.db().Delete(&model.SysAPI{}, id).Error
}

func (d *ApiRepo) ListAll() ([]model.SysAPI, error) {
	var apis []model.SysAPI
	err := d.db().Order("`group` ASC, id ASC").Find(&apis).Error
	return apis, err
}

func (d *ApiRepo) List(page, pageSize int) ([]model.SysAPI, int64, error) {
	var apis []model.SysAPI
	var total int64
	db := d.db().Model(&model.SysAPI{})
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("`group` ASC, id ASC").Find(&apis).Error
	return apis, total, err
}
