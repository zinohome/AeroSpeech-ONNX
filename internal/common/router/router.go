package router

import (
	"github.com/gin-gonic/gin"
)

// Router 路由管理器
type Router struct {
	engine *gin.Engine
}

// NewRouter 创建路由管理器
func NewRouter() *Router {
	gin.SetMode(gin.ReleaseMode)
	return &Router{
		engine: gin.New(),
	}
}

// GetEngine 获取Gin引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// SetupMiddleware 设置中间件
func (r *Router) SetupMiddleware() {
	// 恢复中间件
	r.engine.Use(gin.Recovery())

	// 日志中间件
	r.engine.Use(gin.Logger())
}

// SetupStaticFiles 设置静态文件服务
func (r *Router) SetupStaticFiles(staticPath, templatesPath string) {
	r.engine.Static("/static", staticPath)
	if templatesPath != "" {
		r.engine.LoadHTMLGlob(templatesPath + "/*")
	}
}

// SetupRoutes 设置路由
func (r *Router) SetupRoutes(setupFunc func(*gin.Engine)) {
	if setupFunc != nil {
		setupFunc(r.engine)
	}
}

