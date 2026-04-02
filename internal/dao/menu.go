package dao

import (
	"github.com/karlhsu/frame/internal/app"
	"github.com/karlhsu/frame/internal/model"
	"gorm.io/gorm"
)

type MenuDAO struct{}

func NewMenuDAO() *MenuDAO {
	return &MenuDAO{}
}

func (d *MenuDAO) db() *gorm.DB {
	return app.DB
}

func (d *MenuDAO) Create(menu *model.SysMenu) error {
	return d.db().Create(menu).Error
}

func (d *MenuDAO) GetByID(id uint) (*model.SysMenu, error) {
	var menu model.SysMenu
	err := d.db().Preload("APIs").First(&menu, id).Error
	return &menu, err
}

func (d *MenuDAO) Update(menu *model.SysMenu) error {
	return d.db().Save(menu).Error
}

func (d *MenuDAO) Delete(id uint) error {
	return d.db().Delete(&model.SysMenu{}, id).Error
}

func (d *MenuDAO) ListAll() ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := d.db().Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (d *MenuDAO) GetByIDs(ids []uint) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := d.db().Preload("APIs").Where("id IN ?", ids).Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (d *MenuDAO) SetAPIs(menuID uint, apiIDs []uint) error {
	menu := &model.SysMenu{BaseModel: model.BaseModel{ID: menuID}}
	var apis []model.SysAPI
	for _, id := range apiIDs {
		apis = append(apis, model.SysAPI{BaseModel: model.BaseModel{ID: id}})
	}
	return d.db().Model(menu).Association("APIs").Replace(apis)
}

func (d *MenuDAO) HasChildren(id uint) (bool, error) {
	var count int64
	err := d.db().Model(&model.SysMenu{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

func BuildMenuTree(menus []model.SysMenu, parentID uint) []*model.SysMenu {
	tree := make([]*model.SysMenu, 0)
	for i := range menus {
		if menus[i].ParentID == parentID {
			node := &menus[i]
			node.Children = BuildMenuTree(menus, node.ID)
			tree = append(tree, node)
		}
	}
	return tree
}
