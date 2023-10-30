package api

import (
	"loa/user_content/database"
	"loa/user_content/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)
func HandleLoginPage(ctx *gin.Context){
	cookie, err := ctx.Cookie("token")
	if err != nil {
		ctx.HTML(http.StatusOK,"login.html", gin.H{})
		return
	}
	if cookie != "" {
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	jwt, err := utils.ValidateJWT(cookie)
	if err != nil {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
		return
	}
	if !jwt.Valid {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
		return
	}
	ctx.HTML(http.StatusOK,"login.html",gin.H{})
}
func HandleLogin(c *gin.Context){
	username := c.PostForm("username")
	password := c.PostForm("password")
	user, err := database.DB.GetUser(username)

	if err != nil{
		log.Default().Println(err)
		c.String(http.StatusUnauthorized,"Invalid credentials")
		return
	}
	if !utils.MatchPassword(password,user.Password) {
		c.String(http.StatusUnauthorized,"Invalid credentials")
		return
	}
	jwtToken, err := utils.CreateJWT(user.Username)
	if err != nil{
		log.Default().Println(err)
		c.String(http.StatusInternalServerError,"error creating JWT")
		return
	}
	c.SetCookie("username",user.Username, int(utils.Timeout),"/","localhost", false, true)
	c.SetCookie("token",jwtToken, int(utils.Timeout),"/","localhost", false, true)
    c.String(http.StatusOK,"User Logged in")
}
//check authentication
func HandleCheckAuth(ctx *gin.Context){
	cookie, err := ctx.Cookie("token")
	if err != nil {
		ctx.String(http.StatusUnauthorized,"")
		return
	}
	jwt, err := utils.ValidateJWT(cookie)
	if err != nil {
		ctx.String(http.StatusUnauthorized,"")
		return
	}
	if !jwt.Valid {
		ctx.String(http.StatusUnauthorized,"")
		return
	}
	ctx.String(http.StatusOK,"true")
}
func OnlyAuthed(ctx *gin.Context){
	cookie, err := ctx.Cookie("token")
	if err != nil {
		ctx.Redirect(http.StatusFound,"/about")
		return
	}
	jwt, err := utils.ValidateJWT(cookie)
	if err != nil {
		ctx.Redirect(http.StatusFound,"/about")
		return
	}
	if !jwt.Valid {
		ctx.Redirect(http.StatusFound,"/login")
		return
	}
	
}
func HandleAccount(ctx *gin.Context){
OnlyAuthed(ctx)
username, err := ctx.Cookie("username")
if err != nil {
	ctx.Redirect(http.StatusFound,"/login")
	return
}
docs, err := database.DB.GetUserDocs(username)
if err != nil {
	ctx.Redirect(http.StatusFound,"/login")
	return
}
ctx.HTML(http.StatusOK,"account.html",gin.H{
	"username":username,
	"docs":docs,
})

}
func HandleSignout(ctx *gin.Context){
	ctx.SetCookie("token","",-1,"/","localhost", false, true)
	ctx.Redirect(http.StatusFound,"/login")
}