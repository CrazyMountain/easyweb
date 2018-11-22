package main

import (
	"easyweb/routers"
	"easyweb/utils/setting"
	"fmt"
)

func main() {
	routers.InitRouter().Run(fmt.Sprintf(":%d", setting.Port))
}
