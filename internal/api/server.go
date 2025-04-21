package api

import (
	"html/template"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/time/rate"
)

func newServer(ctx *AppContext) *echo.Echo {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := newRequestContext(c, ctx)
			return next(cc)
		}
	})
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, path=${path}, method=${method} uri=${uri}, status=${status}, error=${error}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(100))))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.getSecretKey()))))

	templateFiles, err := template.ParseGlob("./internal/templates/*.html")
	if err != nil {
		templateFiles, err = template.ParseGlob("../templates/*.html")
	}

	renderer := &TemplateRenderer{
		templates: template.Must(templateFiles, err),
	}
	e.Renderer = renderer

	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetHeader("${time_rfc3339} ${level} ${short_file}:L${line} ${message}")
	return e
}

func CreateServer() *echo.Echo {
	ctx := NewAppContext(make(map[string]interface{}))
	err := ctx.Init()
	if err != nil {
		log.Fatal(err)
	}
	e := newServer(ctx)

	// Routes
	e.GET("/", IndexHandler)
	e.GET(EndpointWellKnown, AuthorizationServerWellKnownHandler)
	e.GET(EndpointJWK, JWKHandler)
	e.POST(EndpointToken, NewTokenHandler)
	e.GET(EndpointAuthorization, AuthorizationHandler)
	e.GET(EndpointAuthorizationLogin, LoginPageHandler)
	e.POST(EndpointAuthorizationLogin, UserLoginHandler)
	return e
}

func newRequestContext(c echo.Context, ctx *AppContext) RequestContext {
	tokenRepo := ctx.TokenRepo.NewTokenRepositoryWithCtx(c.Request().Context())
	return RequestContext{c, ctx.JWKManager, ctx.ClientRepo, ctx.UserRepo, tokenRepo}
}
