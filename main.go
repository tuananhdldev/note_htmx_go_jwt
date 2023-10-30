package main

import (
	"loa/user_content/database"
	"loa/user_content/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
  //inital Database
   err := database.NewStoreInstance()
   if err != nil {
	   panic(err)
   }
   if err = database.DB.CreateTable(); err != nil {
	panic(err)
   }

	r:= gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static","./public")
	r.Use(gin.Logger())
	r.SetTrustedProxies(nil)
	routes.SetRoutes(r);
	log.Fatal(r.Run(":3000"))
}