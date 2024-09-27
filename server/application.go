package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mheers/vcluster-operator/auth"
	"github.com/mheers/vcluster-operator/config"
	"github.com/mheers/vcluster-operator/vclustermanagement"

	jwt "github.com/appleboy/gin-jwt/v2"
)

func NewApplicaton(cfg *config.ServerConfig) *Application {
	router := gin.Default()
	return &Application{
		Router: router,
		Config: cfg,
	}
}

type Application struct {
	Router *gin.Engine
	Config *config.ServerConfig
}

func (a *Application) Run() error {
	err := a.startWebServer()
	if err != nil {
		return err
	}
	return nil
}

func (a *Application) startWebServer() error {
	r := gin.Default()

	// the jwt middleware
	authMiddleware, err := auth.GetAuthMiddleware(a.Config.SecretKey, a.Config.AdminUser, a.Config.AdminPassword)
	if err != nil {
		return errors.New("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		return errors.New("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/api")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())

	vclusters := auth.Group("/vclusters")
	vclusters.GET("/", vclustermanagement.ListHandler)
	vclusters.GET("/:name", vclustermanagement.GetHandler)
	vclusters.GET("/:name/kubeconfig", vclustermanagement.KubeconfigHandler)
	vclusters.GET("/:name/token", authMiddleware.TokenHandler)
	vclusters.POST("/:name", vclustermanagement.CreateHandler)
	vclusters.DELETE("/:name", vclustermanagement.DeleteHandler)

	listenAddress := fmt.Sprintf(":%d", a.Config.Port)
	if err := http.ListenAndServe(listenAddress, r); err != nil {
		return err
	}
	return nil
}
