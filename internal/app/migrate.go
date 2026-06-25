package app

import (
	"github.com/kar1hsu/frame/internal/model"
	"github.com/kar1hsu/frame/internal/pkg/utils"
	"gorm.io/gorm"
)

func AutoMigrate() error {
	return DB.AutoMigrate(
		&model.SysUser{},
		&model.SysRole{},
		&model.SysMenu{},
		&model.SysAPI{},
		&model.SysOperationLog{},
	)
}

func SeedData() error {
	var count int64
	DB.Model(&model.SysRole{}).Where("code = ?", "admin").Count(&count)
	if count > 0 {
		return nil
	}

	Log.Info("seeding initial data...")

	if err := DB.Transaction(func(tx *gorm.DB) error {
		// Seed APIs
		apis := []model.SysAPI{
			{ID: 1, Path: "/admin/users", Method: "GET", Group: "用户管理", Description: "用户列表"},
			{ID: 2, Path: "/admin/users", Method: "POST", Group: "用户管理", Description: "创建用户"},
			{ID: 3, Path: "/admin/users/:id", Method: "GET", Group: "用户管理", Description: "用户详情"},
			{ID: 4, Path: "/admin/users/:id", Method: "PUT", Group: "用户管理", Description: "更新用户"},
			{ID: 5, Path: "/admin/users/:id", Method: "DELETE", Group: "用户管理", Description: "删除用户"},
			{ID: 6, Path: "/admin/roles", Method: "GET", Group: "角色管理", Description: "角色列表"},
			{ID: 7, Path: "/admin/roles", Method: "POST", Group: "角色管理", Description: "创建角色"},
			{ID: 8, Path: "/admin/roles/:id", Method: "GET", Group: "角色管理", Description: "角色详情"},
			{ID: 9, Path: "/admin/roles/:id", Method: "PUT", Group: "角色管理", Description: "更新角色"},
			{ID: 10, Path: "/admin/roles/:id", Method: "DELETE", Group: "角色管理", Description: "删除角色"},
			{ID: 11, Path: "/admin/roles/:id/menus", Method: "PUT", Group: "角色管理", Description: "分配菜单"},
			{ID: 12, Path: "/admin/menus", Method: "POST", Group: "菜单管理", Description: "创建菜单"},
			{ID: 13, Path: "/admin/menus/:id", Method: "GET", Group: "菜单管理", Description: "菜单详情"},
			{ID: 14, Path: "/admin/menus/:id", Method: "PUT", Group: "菜单管理", Description: "更新菜单"},
			{ID: 15, Path: "/admin/menus/:id", Method: "DELETE", Group: "菜单管理", Description: "删除菜单"},
			{ID: 16, Path: "/admin/apis", Method: "GET", Group: "API管理", Description: "API列表"},
			{ID: 17, Path: "/admin/apis", Method: "POST", Group: "API管理", Description: "创建API"},
			{ID: 18, Path: "/admin/apis/:id", Method: "PUT", Group: "API管理", Description: "更新API"},
			{ID: 19, Path: "/admin/apis/:id", Method: "DELETE", Group: "API管理", Description: "删除API"},
		}
		for i := range apis {
			if err := tx.Create(&apis[i]).Error; err != nil {
				return err
			}
		}

		adminRole := &model.SysRole{
			Name:   "超级管理员",
			Code:   "admin",
			Sort:   0,
			Status: 1,
			Remark: "超级管理员，拥有所有权限",
		}
		if err := tx.Create(adminRole).Error; err != nil {
			return err
		}

		hashed, err := utils.HashPassword("admin123")
		if err != nil {
			return err
		}
		adminUser := &model.SysUser{
			Username: "admin",
			Password: hashed,
			Nickname: "管理员",
			Status:   1,
			Roles:    []model.SysRole{*adminRole},
		}
		if err := tx.Create(adminUser).Error; err != nil {
			return err
		}

		// Seed menus with API associations
		// API IDs: 1=用户列表 2=创建用户 3=用户详情 4=更新用户 5=删除用户
		//          6=角色列表 7=创建角色 8=角色详情 9=更新角色 10=删除角色 11=分配菜单
		//          12=创建菜单 13=菜单详情 14=更新菜单 15=删除菜单
		// a[0]=用户列表  a[1]=创建用户  a[2]=用户详情  a[3]=更新用户  a[4]=删除用户
		// a[5]=角色列表  a[6]=创建角色  a[7]=角色详情  a[8]=更新角色  a[9]=删除角色  a[10]=分配菜单
		// a[11]=创建菜单 a[12]=菜单详情 a[13]=更新菜单 a[14]=删除菜单
		a := apis
		menus := []model.SysMenu{
			// ── 系统管理（目录）──
			{ID: 1, ParentID: 0, Name: "系统管理", Path: "/system", Icon: "Setting", Sort: 1, Type: 0, Visible: 1, Status: 1},

			// ── 用户管理 ──
			{ID: 2, ParentID: 1, Name: "用户管理", Path: "/system/user", Component: "system/user/index", Icon: "User", Sort: 1, Type: 1, Permission: "system:user:list", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[0]}},
			{ID: 20, ParentID: 2, Name: "用户详情", Sort: 1, Type: 2, Permission: "system:user:query", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[2]}},
			{ID: 21, ParentID: 2, Name: "新增用户", Sort: 2, Type: 2, Permission: "system:user:add", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[1]}},
			{ID: 22, ParentID: 2, Name: "编辑用户", Sort: 3, Type: 2, Permission: "system:user:edit", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[2], a[3]}},
			{ID: 23, ParentID: 2, Name: "删除用户", Sort: 4, Type: 2, Permission: "system:user:delete", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[4]}},

			// ── 角色管理 ──
			{ID: 3, ParentID: 1, Name: "角色管理", Path: "/system/role", Component: "system/role/index", Icon: "UserFilled", Sort: 2, Type: 1, Permission: "system:role:list", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[5]}},
			{ID: 30, ParentID: 3, Name: "角色详情", Sort: 1, Type: 2, Permission: "system:role:query", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[7]}},
			{ID: 31, ParentID: 3, Name: "新增角色", Sort: 2, Type: 2, Permission: "system:role:add", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[6]}},
			{ID: 32, ParentID: 3, Name: "编辑角色", Sort: 3, Type: 2, Permission: "system:role:edit", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[7], a[8], a[10]}},
			{ID: 33, ParentID: 3, Name: "删除角色", Sort: 4, Type: 2, Permission: "system:role:delete", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[9]}},

			// ── 菜单管理 ──
			{ID: 4, ParentID: 1, Name: "菜单管理", Path: "/system/menu", Component: "system/menu/index", Icon: "Menu", Sort: 3, Type: 1, Permission: "system:menu:list", Visible: 1, Status: 1},
			{ID: 40, ParentID: 4, Name: "菜单详情", Sort: 1, Type: 2, Permission: "system:menu:query", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[12]}},
			{ID: 41, ParentID: 4, Name: "新增菜单", Sort: 2, Type: 2, Permission: "system:menu:add", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[11]}},
			{ID: 42, ParentID: 4, Name: "编辑菜单", Sort: 3, Type: 2, Permission: "system:menu:edit", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[12], a[13]}},
			{ID: 43, ParentID: 4, Name: "删除菜单", Sort: 4, Type: 2, Permission: "system:menu:delete", Visible: 1, Status: 1,
				APIs: []model.SysAPI{a[14]}},
		}

		for i := range menus {
			if err := tx.Create(&menus[i]).Error; err != nil {
				return err
			}
		}

		// Assign all menus to admin role
		menuModels := make([]model.SysMenu, len(menus))
		copy(menuModels, menus)
		if err := tx.Model(adminRole).Association("Menus").Replace(menuModels); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	Log.Warn("已创建默认管理员 admin/admin123，请登录后立即修改密码")
	Log.Info("seed data completed")
	return nil
}

// SeedMenus idempotently inserts menus/APIs that were added after the initial
// SeedData, so existing deployments pick them up on restart without a full
// re-seed. Safe to run on every startup.
func SeedMenus() error {
	return seedOperationLogMenu()
}

// seedOperationLogMenu adds the 系统监控 > 操作日志 menu (with its buttons and
// API records) and grants it to the admin role. Gated on the list permission so
// it runs at most once.
func seedOperationLogMenu() error {
	var cnt int64
	if err := DB.Model(&model.SysMenu{}).
		Where("permission = ?", "monitor:operlog:list").Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}

	Log.Info("seeding operation-log menu...")
	return DB.Transaction(func(tx *gorm.DB) error {
		// APIs (FirstOrCreate by path+method so reruns/manual rows don't dup)
		apis := []model.SysAPI{
			{Path: "/admin/operation-logs", Method: "GET", Group: "操作日志", Description: "日志列表"},
			{Path: "/admin/operation-logs/:id", Method: "GET", Group: "操作日志", Description: "日志详情"},
			{Path: "/admin/operation-logs/:id", Method: "DELETE", Group: "操作日志", Description: "删除日志"},
			{Path: "/admin/operation-logs", Method: "DELETE", Group: "操作日志", Description: "清空日志"},
		}
		for i := range apis {
			if err := tx.Where("path = ? AND method = ?", apis[i].Path, apis[i].Method).
				FirstOrCreate(&apis[i]).Error; err != nil {
				return err
			}
		}

		// 系统监控 directory (reuse if it already exists)
		dir := model.SysMenu{Name: "系统监控", Path: "/monitor", Icon: "Monitor", Sort: 2, Type: 0, Visible: 1, Status: 1}
		dirRes := tx.Where("name = ? AND parent_id = 0", "系统监控").
			Attrs(model.SysMenu{Path: "/monitor", Icon: "Monitor", Sort: 2, Type: 0, Visible: 1, Status: 1}).
			FirstOrCreate(&dir)
		if dirRes.Error != nil {
			return dirRes.Error
		}

		// 操作日志 menu + buttons
		listMenu := model.SysMenu{
			ParentID: dir.ID, Name: "操作日志", Path: "/monitor/operlog", Component: "monitor/operlog/index",
			Icon: "Document", Sort: 1, Type: 1, Permission: "monitor:operlog:list", Visible: 1, Status: 1,
			APIs: []model.SysAPI{apis[0], apis[1]},
		}
		deleteBtn := model.SysMenu{
			ParentID: 0, Name: "删除日志", Sort: 2, Type: 2, Permission: "monitor:operlog:delete", Visible: 1, Status: 1,
			APIs: []model.SysAPI{apis[2]},
		}
		clearBtn := model.SysMenu{
			ParentID: 0, Name: "清空日志", Sort: 3, Type: 2, Permission: "monitor:operlog:clear", Visible: 1, Status: 1,
			APIs: []model.SysAPI{apis[3]},
		}
		if err := tx.Create(&listMenu).Error; err != nil {
			return err
		}
		deleteBtn.ParentID = listMenu.ID
		clearBtn.ParentID = listMenu.ID
		if err := tx.Create(&deleteBtn).Error; err != nil {
			return err
		}
		if err := tx.Create(&clearBtn).Error; err != nil {
			return err
		}

		// Grant the new menus to the admin role for display in its menu tree.
		var admin model.SysRole
		if err := tx.Where("code = ?", "admin").First(&admin).Error; err != nil {
			return err
		}
		grant := []model.SysMenu{listMenu, deleteBtn, clearBtn}
		if dirRes.RowsAffected > 0 { // dir newly created → also grant it
			grant = append([]model.SysMenu{dir}, grant...)
		}
		if err := tx.Model(&admin).Association("Menus").Append(grant); err != nil {
			return err
		}
		return nil
	})
}
