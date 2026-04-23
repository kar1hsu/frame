package server

import (
	"embed"
	"io/fs"
	"net/http"

	"frame/internal/app"
	"frame/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Module interface {
	Name() string
	RegisterRoutes(rg *gin.RouterGroup)
}

func NewRouter(adminDist embed.FS, modules ...Module) *gin.Engine {
	gin.SetMode(app.Cfg.Server.Mode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())

	for _, m := range modules {
		group := r.Group("/" + m.Name())
		m.RegisterRoutes(group)
		app.Log.Infof("module [%s] registered at /%s", m.Name(), m.Name())
	}

	setupStaticFiles(r, adminDist)

	return r
}

func setupStaticFiles(r *gin.Engine, adminDist embed.FS) {
	subFS, err := fs.Sub(adminDist, "web/admin/dist")
	if err != nil {
		app.Log.Warnf("embed admin dist not found: %v", err)
		return
	}

	staticHandler := http.FileServer(http.FS(subFS))

	r.NoRoute(func(c *gin.Context) {
		// Serve static files for the admin panel
		path := c.Request.URL.Path
		// Try to serve the file directly
		f, err := subFS.Open(path[1:]) // strip leading /
		if err == nil {
			f.Close()
			staticHandler.ServeHTTP(c.Writer, c.Request)
			return
		}
		// Fallback to index.html for SPA routing
		c.Request.URL.Path = "/"
		staticHandler.ServeHTTP(c.Writer, c.Request)
	})

	app.Log.Info("admin panel static files registered")
}
