package api

import (
	"loa/user_content/database"
	"loa/user_content/types"
	"loa/user_content/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)
func HandleHome(ctx *gin.Context){
	OnlyAuthed(ctx)
	var username, _  = ctx.Cookie("username")
	ctx.HTML(http.StatusOK,"index.html",gin.H{
		"username":username,
	})
	
}
func HandleCreateForm(ctx *gin.Context){
	title := ctx.DefaultPostForm("title","Untitled")
	if title == ""{
		title = "Untitled"
	}
	author:= ctx.DefaultPostForm("author","Anonymous")
	content := ctx.PostForm("content")
	if content == ""{
		ctx.String(http.StatusBadRequest,"Content cannot be empty")
		return
	}
	hash := utils.RanHash()
	req:= &types.CreateDocumentRequest{
		Hash: hash,
		Author: author,
		Title: title,
		Content: content,
	}
	err := database.DB.InsertDoc(req)
	if err != nil{
		log.Default().Println(err)
		ctx.JSON(http.StatusInternalServerError,err)
		return
	}
	ctx.String(http.StatusOK,"<a class='text-blue-600' href='/docs/"+hash+"'>View</a>")
}
func HandleGetByHash(ctx *gin.Context){
	hash := ctx.Param("hash")
	doc, err := database.DB.GetByHash(hash)
	if err != nil{
		log.Default().Println(err)
		ctx.JSON(http.StatusInternalServerError,err)
		return
	}
	ctx.HTML(http.StatusOK,"doc.html",gin.H{
		"title":doc.Title,
		"content":doc.Content,
		"author":doc.Author,
		"created_at":doc.CreatedAt,	
	})
}
func HandleAbout(ctx *gin.Context){
	ctx.HTML(http.StatusOK,"about.html",gin.H{})
}

func HandleSignupPage(ctx *gin.Context){
	ctx.HTML(http.StatusOK,"signup.html",gin.H{})
}
func HandleSignup(ctx *gin.Context){
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	userRequest := types.CreateUserRequest{
            Username: username,
			Password: password,
	}
	user, _ := database.DB.GetUser(username)
	// if err != nil{
	// 	log.Default().Println(err)
	// 	ctx.JSON(http.StatusInternalServerError,err)
	// 	return
	// }
	if user.Username  != ""{
		ctx.JSON(http.StatusOK,"user already exists.")
		return
	}
	err :=database.DB.CreateUser(userRequest)
	if err != nil{
		log.Default().Println(err)
		ctx.JSON(http.StatusInternalServerError,err)
		return
	}
	jwtToken, err := utils.CreateJWT(username)
	if err != nil{
		log.Default().Panicln(err)
		ctx.JSON(http.StatusInternalServerError,err)
		return
	}
	ctx.SetCookie("token",jwtToken, 3*60,"/","localhost", false,true)
	ctx.SetCookie("username",username, 3*60,"/","localhost", false,true)
	ctx.String(http.StatusOK,"User created")
}