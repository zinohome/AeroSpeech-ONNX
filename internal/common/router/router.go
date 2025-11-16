package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhangjun/AeroSpeech-ONNX/internal/common/middleware"
)

// Router 路由管理器
type Router struct {
	engine      *gin.Engine
	rateLimiter *middleware.RateLimiter
}

// NewRouter 创建路由管理器
func NewRouter() *Router {
	gin.SetMode(gin.ReleaseMode)
	return &Router{
		engine: gin.New(),
	}
}

// NewRouterWithRateLimit 创建带限流的路由管理器
func NewRouterWithRateLimit(rateLimiter *middleware.RateLimiter) *Router {
	gin.SetMode(gin.ReleaseMode)
	return &Router{
		engine:      gin.New(),
		rateLimiter: rateLimiter,
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

	// 限流中间件
	if r.rateLimiter != nil {
		r.engine.Use(func(c *gin.Context) {
			r.rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				c.Next()
			})).ServeHTTP(c.Writer, c.Request)
		})
	}
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

// GetRateLimiter 获取限流器
func (r *Router) GetRateLimiter() *middleware.RateLimiter {
	return r.rateLimiter
}

