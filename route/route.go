package route

import (
	"leave-order/config"
	"leave-order/controller"
	"leave-order/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Process(host string) {

	h := controller.Handler{}
	h.Db = config.Connect().DB

	e := echo.New()
	model.DBMigrate(h.Db)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))

	jwtGroup := e.Group("/jwt")
	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("secretpassword"),
	}))



	jwtGroup.GET("/availableleavesindividual", h.Availableleavesindividual)
	jwtGroup.GET("/availableleaves", h.Availableleaves)


	// e.GET("/users", h.AllUsers)
	jwtGroup.GET("/transactionindividual", h.IndividualTransactions)
	jwtGroup.GET("/user/:id", h.User)
	jwtGroup.GET("/users", h.AllUsers)
	jwtGroup.GET("/foodscheckbox/:checkedid", h.Foodscheckbox)
	// e.GET("/users", h.AllUsers)

	jwtGroup.GET("/transactions", h.AllTransactions)
	jwtGroup.POST("/addfunds", h.Addfunds)

	jwtGroup.POST("/approveleave/:id", h.Approveleave)
	jwtGroup.POST("/rejectleave/:id/:availableleaveid/:applydays", h.Rejectleave)
	jwtGroup.GET("/checkaccess", h.Checkaccess)

	jwtGroup.GET("/pendingleaves", h.Allpendingleaves)


	jwtGroup.GET("/notpendingleaves", h.Allnotpendingleaves)



	// jwtGroup.GET("/notpendingindividualleaves", h.Allnotpendingindividualleaves)

	jwtGroup.POST("/addleaves", h.Addleave)
	// Route
	e.POST("/createuser", h.CreateUser)
	e.POST("/login", h.Login)
	e.POST("/upload", h.Upload)
	e.GET("/myimage/:id", h.Uploadfilepath)





	e.POST("/userimageupload", h.Userimageupload)


	e.GET("/myuserimageupload/:id", h.Uploaduserfilepath)







	jwtGroup.GET("/foodscheckbox/:checkedid", h.Foodscheckbox)
	e.POST("/addfood", h.Addfood)
	
	e.POST("/updatefood", h.Updatefood)
	e.POST("/deletefood", h.Deletefood)
	e.POST("/addpoll", h.Addpoll)
	e.GET("/findpolldetail", h.Findpolldetail)
	e.GET("/lastpolldetail", h.Lastpolldetail)
	e.GET("/fooddisplay", h.Fooddisplay)
	e.GET("/findfood", h.Findfood)

	e.GET("/calculateselectedquantityprice", h.Calculateselectedquantityprice)
	
	e.GET("/calculatetotalselectedquantityprice", h.Calculatetotalselectedquantityprice)
	
	e.GET("/clearamountholder", h.Clearamountholder)



	e.GET("/ordersummarydisplay", h.Ordersummarydisplay)


	jwtGroup.GET("/transactionhistoryindividual",h.Transactionhistoryindividual)
	jwtGroup.GET("/orderhistory", h.Orderhistory)

	jwtGroup.GET("/leavetypedisplay",h.Leavetypedisplay)

	jwtGroup.GET("/leavehistoryindividual",h.Leavehistoryindividual)
	// jwtGroup.POST("/changepassword/:id", h.Changepassword)
	jwtGroup.POST("/changepassword", h.Changepassword)

	jwtGroup.POST("/createorder", h.Createorder)

	jwtGroup.POST("/closepoll", h.Closepoll)

	e.Logger.Fatal(e.Start(host))
}
