package router

import (
	"demo/handlers"

	"github.com/gin-gonic/gin"
)

func Userroutes(router *gin.Engine) {

	router.POST("/resetting", handlers.Resetting)
	router.GET("/statement", handlers.Statement)
	router.POST("/Withdrawal", handlers.Withdrawal)

	router.POST("/transfer_money", handlers.Transfer_money)
	router.POST("/deposite_money", handlers.Deposite_money)
	router.POST("/create_acc", handlers.Create_acc)
}
