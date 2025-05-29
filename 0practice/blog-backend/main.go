package main

import (
	"blog-backend/config"
	"blog-backend/routers"
	"blog-backend/utils"
)

func main() {
	config.ConnectDB()

	utils.InitLogger()
	r := routers.SetupRouter()
	r.Run(":8080")
}
