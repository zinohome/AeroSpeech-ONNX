package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Fatal("NewRouter() returned nil")
	}
	if router.engine == nil {
		t.Fatal("Router engine is nil")
	}
}

func TestGetEngine(t *testing.T) {
	router := NewRouter()
	engine := router.GetEngine()
	if engine == nil {
		t.Fatal("GetEngine() returned nil")
	}
}

func TestSetupMiddleware(t *testing.T) {
	router := NewRouter()
	router.SetupMiddleware()

	// 测试中间件是否设置
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.engine.ServeHTTP(w, req)
}

func TestSetupStaticFiles(t *testing.T) {
	router := NewRouter()
	// 不设置模板路径，避免模板文件不存在的问题
	router.SetupStaticFiles("/tmp/static", "")

	// 测试静态文件路由是否设置
	req := httptest.NewRequest("GET", "/static/test.js", nil)
	w := httptest.NewRecorder()
	router.engine.ServeHTTP(w, req)
	// 404是预期的，因为文件不存在
}

func TestSetupRoutes(t *testing.T) {
	router := NewRouter()
	
	setupFunc := func(engine *gin.Engine) {
		engine.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "test"})
		})
	}
	
	router.SetupRoutes(setupFunc)

	// 测试路由是否设置
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSetupRoutesNil(t *testing.T) {
	router := NewRouter()
	router.SetupRoutes(nil)
	// 不应该panic
}

