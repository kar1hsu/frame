package repository

import (
	"github.com/kar1hsu/frame/internal/app"
	"github.com/kar1hsu/frame/internal/model"
	"gorm.io/gorm"
)

type MenuRepo struct{}

func NewMenuRepo() *MenuRepo {
	return &MenuRepo{}
}

func (d *MenuRepo) db() *gorm.DB {
	return app.DB
}

func (d *MenuRepo) Create(menu *model.SysMenu) error {
	return d.db().Create(menu).Error
}

func (d *MenuRepo) GetByID(id uint) (*model.SysMenu, error) {
	var menu model.SysMenu
	if err := d.db().Preload("APIs").First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (d *MenuRepo) Update(menu *model.SysMenu) error {
	// Only update base columns; the APIs association is managed by SetAPIs.
	return d.db().Model(&model.SysMenu{ID: menu.ID}).
		Select("ParentID", "Name", "Path", "Component", "Icon", "Sort", "Type", "Permission", "Visible", "Status").
		Updates(menu).Error
}

func (d *MenuRepo) Delete(id uint) error {
	return d.db().Delete(&model.SysMenu{}, id).Error
}

func (d *MenuRepo) ListAll() ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := d.db().Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (d *MenuRepo) GetByIDs(ids []uint) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := d.db().Preload("APIs").Where("id IN ?", ids).Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (d *MenuRepo) SetAPIs(menuID uint, apiIDs []uint) error {
	menu := &model.SysMenu{ID: menuID}
	var apis []model.SysAPI
	for _, id := range apiIDs {
		apis = append(apis, model.SysAPI{ID: id})
	}
	return d.db().Model(menu).Association("APIs").Replace(apis)
}

func (d *MenuRepo) HasChildren(id uint) (bool, error) {
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
