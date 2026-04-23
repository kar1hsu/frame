package service

import (
	"errors"

	"frame/internal/dao"
	"frame/internal/model"
	"frame/internal/pkg/utils"
	"gorm.io/gorm"
)

type UserService struct {
	userDAO *dao.UserDAO
}

func NewUserService() *UserService {
	return &UserService{userDAO: dao.NewUserDAO()}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int8   `json:"status"`
	RoleIDs  []uint `json:"role_ids"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   int8   `json:"status"`
	Password string `json:"password"`
	RoleIDs  []uint `json:"role_ids"`
}

func (s *UserService) Create(req *CreateUserRequest) error {
	_, err := s.userDAO.GetByUsername(req.Username)
	if err == nil {
		return errors.New("用户名已存在")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user := &model.SysUser{
		Username: req.Username,
		Password: hashed,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
	}

	if err := s.userDAO.Create(user); err != nil {
		return err
	}

	if len(req.RoleIDs) > 0 {
		return s.userDAO.SetRoles(user.ID, req.RoleIDs)
	}
	return nil
}

func (s *UserService) GetByID(id uint) (*model.SysUser, error) {
	return s.userDAO.GetByID(id)
}

func (s *UserService) Update(id uint, req *UpdateUserRequest) error {
	user, err := s.userDAO.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Status != 0 {
		user.Status = req.Status
	}
	if req.Password != "" {
		hashed, err := utils.HashPassword(req.Password)
		if err != nil {
			return errors.New("密码加密失败")
		}
		user.Password = hashed
	}

	if err := s.userDAO.Update(user); err != nil {
		return err
	}

	if req.RoleIDs != nil {
		return s.userDAO.SetRoles(id, req.RoleIDs)
	}
	return nil
}

func (s *UserService) Delete(id uint) error {
	return s.userDAO.Delete(id)
}

func (s *UserService) List(page, pageSize int) ([]model.SysUser, int64, error) {
	return s.userDAO.List(page, pageSize)
}

func (s *UserService) GetProfile(id uint) (*model.SysUser, error) {
	user, err := s.userDAO.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return user, nil
}
