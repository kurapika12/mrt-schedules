package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kurapika12/mrt-schedules/modules/station"
)


func main()  {
	InitiateRouter()
}

func InitiateRouter(){
	var  router = gin.Default()
	var api = router.Group("/v1/api")

	station.Initiate(api)

	router.Run(":8080")
}