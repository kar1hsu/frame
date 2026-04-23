package service

import (
	"errors"

	"frame/internal/dao"
	"frame/internal/model"
)

type MenuService struct {
	menuDAO *dao.MenuDAO
	roleDAO *dao.RoleDAO
}

func NewMenuService() *MenuService {
	return &MenuService{
		menuDAO: dao.NewMenuDAO(),
		roleDAO: dao.NewRoleDAO(),
	}
}

type CreateMenuRequest struct {
	ParentID   uint   `json:"parent_id"`
	Name       string `json:"name" binding:"required"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	Type       int8   `json:"type"`
	Permission string `json:"permission"`
	Visible    int8   `json:"visible"`
	Status     int8   `json:"status"`
	APIIDs     []uint `json:"api_ids"`
}

type UpdateMenuRequest struct {
	ParentID   *uint   `json:"parent_id"`
	Name       string  `json:"name"`
	Path       string  `json:"path"`
	Component  string  `json:"component"`
	Icon       string  `json:"icon"`
	Sort       int     `json:"sort"`
	Type       int8    `json:"type"`
	Permission string  `json:"permission"`
	Visible    int8    `json:"visible"`
	Status     int8    `json:"status"`
	APIIDs     *[]uint `json:"api_ids"`
}

func (s *MenuService) Create(req *CreateMenuRequest) error {
	menu := &model.SysMenu{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Path:       req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		Sort:       req.Sort,
		Type:       req.Type,
		Permission: req.Permission,
		Visible:    req.Visible,
		Status:     req.Status,
	}
	if err := s.menuDAO.Create(menu); err != nil {
		return err
	}
	if len(req.APIIDs) > 0 {
		return s.menuDAO.SetAPIs(menu.ID, req.APIIDs)
	}
	return nil
}

func (s *MenuService) GetByID(id uint) (*model.SysMenu, error) {
	return s.menuDAO.GetByID(id)
}

func (s *MenuService) Update(id uint, req *UpdateMenuRequest) error {
	menu, err := s.menuDAO.GetByID(id)
	if err != nil {
		return errors.New("菜单不存在")
	}

	if req.ParentID != nil {
		menu.ParentID = *req.ParentID
	}
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Component != "" {
		menu.Component = req.Component
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	menu.Sort = req.Sort
	menu.Type = req.Type
	if req.Permission != "" {
		menu.Permission = req.Permission
	}
	menu.Visible = req.Visible
	menu.Status = req.Status

	if err := s.menuDAO.Update(menu); err != nil {
		return err
	}
	if req.APIIDs != nil {
		return s.menuDAO.SetAPIs(id, *req.APIIDs)
	}
	return nil
}

func (s *MenuService) Delete(id uint) error {
	has, err := s.menuDAO.HasChildren(id)
	if err != nil {
		return err
	}
	if has {
		return errors.New("存在子菜单，无法删除")
	}
	return s.menuDAO.Delete(id)
}

func (s *MenuService) GetTree() ([]*model.SysMenu, error) {
	menus, err := s.menuDAO.ListAll()
	if err != nil {
		return nil, err
	}
	return dao.BuildMenuTree(menus, 0), nil
}

func (s *MenuService) GetUserMenuTree(userID uint) ([]*model.SysMenu, error) {
	userDAO := dao.NewUserDAO()
	user, err := userDAO.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	menuIDSet := make(map[uint]bool)
	for _, role := range user.Roles {
		menus, err := s.roleDAO.GetMenusByRoleID(role.ID)
		if err != nil {
			continue
		}
		for _, m := range menus {
			menuIDSet[m.ID] = true
		}
	}

	if len(menuIDSet) == 0 {
		return []*model.SysMenu{}, nil
	}

	ids := make([]uint, 0, len(menuIDSet))
	for id := range menuIDSet {
		ids = append(ids, id)
	}

	menus, err := s.menuDAO.GetByIDs(ids)
	if err != nil {
		return nil, err
	}

	// Filter: only directories (0) and menus (1), visible and enabled
	visible := make([]model.SysMenu, 0, len(menus))
	for _, m := range menus {
		if m.Type <= 1 && m.Visible == 1 && m.Status == 1 {
			visible = append(visible, m)
		}
	}

	return dao.BuildMenuTree(visible, 0), nil
}

// GetUserPermissions returns all permission identifiers (including buttons) for a user
func (s *MenuService) GetUserPermissions(userID uint) ([]string, error) {
	userDAO := dao.NewUserDAO()
	user, err := userDAO.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// Super admin has all permissions
	for _, role := range user.Roles {
		if role.Code == "admin" {
			return []string{"*"}, nil
		}
	}

	menuIDSet := make(map[uint]bool)
	for _, role := range user.Roles {
		menus, err := s.roleDAO.GetMenusByRoleID(role.ID)
		if err != nil {
			continue
		}
		for _, m := range menus {
			menuIDSet[m.ID] = true
		}
	}

	if len(menuIDSet) == 0 {
		return []string{}, nil
	}

	ids := make([]uint, 0, len(menuIDSet))
	for id := range menuIDSet {
		ids = append(ids, id)
	}

	menus, err := s.menuDAO.GetByIDs(ids)
	if err != nil {
		return nil, err
	}

	perms := make([]string, 0)
	for _, m := range menus {
		if m.Permission != "" && m.Status == 1 {
			perms = append(perms, m.Permission)
		}
	}
	return perms, nil
}
