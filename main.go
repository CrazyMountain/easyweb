package main

import (
	"easyweb/models"
	"easyweb/routers"
	"easyweb/utils/setting"
	"fmt"
)

func main() {
	fmt.Println(models.IsTableExists("users"))
	port := setting.Port
	routers.InitRouter().Run(fmt.Sprintf(":%d", port))
}
