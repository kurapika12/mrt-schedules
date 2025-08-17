package station

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/kurapika12/mrt-schedules/common/response"
)

func Initiate(router *gin.RouterGroup) {

	stationService := NewService()
	station := router.Group("/stations")
	station.GET("", func(c *gin.Context) {
		GetAllStation(c, stationService)

	})

	station.GET("/:id", func(c *gin.Context) {
		CheckSchedulesByStation(c, stationService)
	})
}

func GetAllStation(c *gin.Context, service Service) {
	datas, err := service.GetAllStations()
	if err != nil {
		c.JSON(http.StatusBadRequest, 
			response.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		response.APIResponse{
			Success: true,
			Message: "successfully get all stations",
			Data:    datas,
		},
	)

}

func CheckSchedulesByStation(c *gin.Context, service Service) {
	id := c.Param("id")

	datas, err := service.CheckScheduleByStation(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, 
			response.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		response.APIResponse{
			Success: true,
			Message: "successfully get schedules by station",
			Data:    datas,
		},
	)
}