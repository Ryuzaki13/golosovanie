package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// server()

	// создание роутера
	var router_setup = func() *gin.Engine {
		var router = gin.Default()

		router.Static("/static/", "./static")
		router.LoadHTMLGlob("templates/*")

		// данные аккаунта
		router.POST("/api/user/select-info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"Data": gin.H{"Phone": "14819"}, "Error": nil})
		})
		// главная
		router.GET("/", func(c *gin.Context) {
			// c.JSON(http.StatusOK, gin.H{"Data": gin.H{"Phone": "14819"}, "Error": nil})
			c.HTML(http.StatusOK, "index.html", gin.H{})
		})


		return router
	}



	router_setup().Run()

}
