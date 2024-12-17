/// /gateway/main.go

package main

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"golang.org/x/time/rate"
)

// IP별 레이트 리미터를 저장할 맵과 이를 동기화하기 위한 뮤텍스
var (
	ips = make(map[string]*rate.Limiter)
	mu  sync.RWMutex
)

// 특정 IP 주소에 대한 레이트 리미터를 반환
func GetLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := ips[ip]
	if !exists {
		limiter = rate.NewLimiter(20, 20) // 레이트 리미팅 설정 조정
		ips[ip] = limiter
	}

	return limiter
}

func getClientIP(c *fiber.Ctx) string {
	// X-Real-IP 헤더를 확인
	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// X-Forwarded-For 헤더를 확인
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0] // 여러 IP가 쉼표로 구분되어 있을 수 있음
	}
	// 헤더가 없는 경우 기본 메서드 사용
	return c.IP()
}

// IP 주소별로 레이트 리미팅을 적용
func IPRateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Swagger UI에 대한 요청은 레이트 리미팅에서 제외
		if strings.HasPrefix(c.Path(), "/swagger/") {
			return c.Next()
		}

		ip := getClientIP(c)
		limiter := GetLimiter(ip)

		if !limiter.Allow() {
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error": "요청 수가 너무 많습니다",
			})
		}

		return c.Next()
	}
}
func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(IPRateLimitMiddleware())

	// 서비스로의 리버스 프록시 설정

	setupProxy(app, "/admin/*", "http://admin:44400")
	setupProxy(app, "/emotion/*", "http://emotion:44404")
	setupProxy(app, "/exercise/*", "http://exercise:44405")
	setupProxy(app, "/inquire/*", "http://inquire:44406")
	setupProxy(app, "/medicine/*", "http://medicine:44407")
	setupProxy(app, "/notification/*", "http://notification:44408")
	setupProxy(app, "/user/*", "http://user:44409")
	setupProxy(app, "/video/*", "http://video:44410")
	setupProxy(app, "/check/*", "http://check:44411")

	setupProxy(app, "/landing/*", "http://landing:44500")

	// Swagger UI 프록시 설정
	setupSwaggerUIProxy(app, "/admin-service/swagger/*", "http://admin:44400/swagger")
	setupSwaggerUIProxy(app, "/emotion-service/swagger/*", "http://emotion:44404/swagger")
	setupSwaggerUIProxy(app, "/exercise-service/swagger/*", "http://exercise:44405/swagger")
	setupSwaggerUIProxy(app, "/inquire-service/swagger/*", "http://inquire:44406/swagger")
	setupSwaggerUIProxy(app, "/medicine-service/swagger/*", "http://medicine:44407/swagger")
	setupSwaggerUIProxy(app, "/notification-service/swagger/*", "http://notification:44408/swagger")
	setupSwaggerUIProxy(app, "/user-service/swagger/*", "http://user:44409/swagger")
	setupSwaggerUIProxy(app, "/video-service/swagger/*", "http://video:44410/swagger")
	setupSwaggerUIProxy(app, "/check-service/swagger/*", "http://check:44411/swagger")

	setupSwaggerUIProxy(app, "/landing-service/swagger/*", "http://landing:44500/swagger")

	// Swagger JSON 파일 리다이렉트
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		referer := c.Get("Referer")
		if strings.Contains(referer, "/user-service/") {
			return c.Redirect("/user-service/swagger/doc.json")
		} else if strings.Contains(referer, "/inquire-service/") {
			return c.Redirect("/inquire-service/swagger/doc.json")

		} else if strings.Contains(referer, "/emotion-service/") {
			return c.Redirect("/emotion-service/swagger/doc.json")

		} else if strings.Contains(referer, "/medicine-service/") {
			return c.Redirect("/medicine-service/swagger/doc.json")

		} else if strings.Contains(referer, "/notification-service/") {
			return c.Redirect("/notification-service/swagger/doc.json")

		} else if strings.Contains(referer, "/exercise-service/") {
			return c.Redirect("/exercise-service/swagger/doc.json")

		} else if strings.Contains(referer, "/video-service/") {
			return c.Redirect("/video-service/swagger/doc.json")

		} else if strings.Contains(referer, "/admin-service/") {
			return c.Redirect("/admin-service/swagger/doc.json")

		} else if strings.Contains(referer, "/check-service/") {
			return c.Redirect("/check-service/swagger/doc.json")

		} else if strings.Contains(referer, "/landing-service/") {
			return c.Redirect("/landing-service/swagger/doc.json")
		}
		return c.SendStatus(fiber.StatusNotFound)
	})
	// API 게이트웨이 서버 시작
	app.Listen(":40000")
}

// 서비스로의 리버스 프록시 설정 함수
func setupProxy(app *fiber.App, path string, target string) {
	app.All(path, func(c *fiber.Ctx) error {
		originalPath := c.Params("*")
		originalQuery := c.Request().URI().QueryString()
		targetURL := target + "/" + originalPath
		if len(originalQuery) > 0 {
			targetURL += "?" + string(originalQuery)
		}
		return proxy.Do(c, targetURL)
	})
}

func setupSwaggerUIProxy(app *fiber.App, path string, target string) {
	app.All(path, func(c *fiber.Ctx) error {
		originalPath := c.Params("*")
		targetURL := target
		if originalPath != "" {
			targetURL += "/" + originalPath
		}
		return proxy.Do(c, targetURL)
	})
}
