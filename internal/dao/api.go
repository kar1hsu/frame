package dao

import (
	"github.com/karlhsu/frame/internal/app"
	"github.com/karlhsu/frame/internal/model"
	"gorm.io/gorm"
)

type APIDAO struct{}

func NewAPIDAO() *APIDAO {
	return &APIDAO{}
}

func (d *APIDAO) db() *gorm.DB {
	return app.DB
}

func (d *APIDAO) Create(api *model.SysAPI) error {
	return d.db().Create(api).Error
}

func (d *APIDAO) GetByID(id uint) (*model.SysAPI, error) {
	var api model.SysAPI
	err := d.db().First(&api, id).Error
	return &api, err
}

func (d *APIDAO) Update(api *model.SysAPI) error {
	return d.db().Save(api).Error
}

func (d *APIDAO) Delete(id uint) error {
	return d.db().Delete(&model.SysAPI{}, id).Error
}

func (d *APIDAO) ListAll() ([]model.SysAPI, error) {
	var apis []model.SysAPI
	err := d.db().Order("`group` ASC, id ASC").Find(&apis).Error
	return apis, err
}

func (d *APIDAO) List(page, pageSize int) ([]model.SysAPI, int64, error) {
	var apis []model.SysAPI
	var total int64
	db := d.db().Model(&model.SysAPI{})
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("`group` ASC, id ASC").Find(&apis).Error
	return apis, total, err
}
