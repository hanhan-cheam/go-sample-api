package route

import (
	"sample-go/config"
	"sample-go/controller"
	"sample-go/controller/base"
	"sample-go/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Process(host string) {
	h := controller.Handler{}
	b := base.Handler{}
	h.Db = config.Connect().DB
	b.Db = config.Connect().DB
	e := echo.New()
	model.DBMigrate(h.Db)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))

	jwtGroup := e.Group("/jwt")
	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("secretpassword"),
	}))

	e.GET("/user/:id", b.User)
	e.GET("/users", b.AllUsers)
	e.POST("/createuser", b.CreateUser)
	e.POST("/deleteuser", b.DeleteUser)
	e.POST("/updateuser", b.UpdateUser)

	e.Logger.Fatal(e.Start(host))
}
