package service

import (
	"context"
	"errors"

	"github.com/kar1hsu/frame/internal/model"
	"github.com/kar1hsu/frame/internal/pkg/setting"
	"github.com/kar1hsu/frame/internal/repository"
)

type ConfigService struct {
	repo *repository.ConfigRepo
}

func NewConfigService() *ConfigService {
	return &ConfigService{repo: repository.NewConfigRepo()}
}

func (s *ConfigService) List(ctx context.Context, group string) ([]model.SysConfig, error) {
	q := &repository.QueryOptions{Order: []string{"sort ASC", "id ASC"}}
	if group != "" {
		q.Where = map[string]interface{}{"group": group}
	}
	return s.repo.List(ctx, q)
}

type ConfigItem struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value"`
}

// BatchUpdate writes many values in one transaction, then re-syncs the whole
// cache from DB (so all instances converge on next read).
func (s *ConfigService) BatchUpdate(ctx context.Context, items []ConfigItem) error {
	if len(items) == 0 {
		return nil
	}
	if err := repository.Transaction(ctx, func(ctx context.Context) error {
		for _, it := range items {
			if err := s.repo.UpdateValue(ctx, it.Key, it.Value); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return setting.RefreshAll(context.Background())
}

type CreateConfigRequest struct {
	Group    string `json:"group"`
	Key      string `json:"key" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Value    string `json:"value"`
	Type     string `json:"type"`
	Options  string `json:"options"`
	IsPublic bool   `json:"is_public"`
	Remark   string `json:"remark"`
	Sort     int    `json:"sort"`
}

func (s *ConfigService) Create(ctx context.Context, req *CreateConfigRequest) error {
	exists, err := s.repo.Exists(ctx, &repository.QueryOptions{
		Where: map[string]interface{}{"key": req.Key},
	})
	if err != nil {
		return err
	}
	if exists {
		return errors.New("配置键已存在")
	}
	typ := req.Type
	if typ == "" {
		typ = "string"
	}
	c := &model.SysConfig{
		Group: req.Group, Key: req.Key, Name: req.Name, Value: req.Value,
		Type: typ, Options: req.Options, IsPublic: req.IsPublic,
		Remark: req.Remark, Sort: req.Sort, Editable: true, Builtin: false,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return err
	}
	_ = setting.RefreshKey(context.Background(), req.Key)
	return nil
}

func (s *ConfigService) Delete(ctx context.Context, id uint) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return notFoundOr(err, "配置不存在")
	}
	if c.Builtin {
		return errors.New("内置配置不可删除")
	}
	// Hard delete: configs are reference data, and a soft-deleted row would keep
	// the unique key, blocking re-creating the same key later.
	if err := s.repo.HardDelete(ctx, id); err != nil {
		return err
	}
	_ = setting.RefreshKey(context.Background(), c.Key) // drops the cache field
	return nil
}

// Refresh re-syncs the cache: a single key when given, otherwise everything.
func (s *ConfigService) Refresh(ctx context.Context, key string) error {
	if key != "" {
		return setting.RefreshKey(ctx, key)
	}
	return setting.RefreshAll(ctx)
}
