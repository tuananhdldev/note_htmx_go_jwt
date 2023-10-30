package routes

import (
	"loa/user_content/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine){
	r.GET("/", api.HandleHome)
	r.GET("/signup", api.HandleSignupPage)
	r.POST("/createUser", api.HandleSignup)
	r.GET("/login", api.HandleLoginPage)
	r.POST("/loginUser", api.HandleLogin)
	r.GET("/checkAuth", api.HandleCheckAuth)
	r.GET("/signout", api.HandleSignout)
	r.GET("/about", api.HandleAbout)
	r.POST("/create", api.HandleCreateForm)
	r.GET("/docs/:hash", api.HandleGetByHash)
	r.GET("/account", api.HandleAccount)
	r.GET("/checkExists", func (ctx *gin.Context){
      ctx.String(http.StatusOK,"check exists")
	})
}
