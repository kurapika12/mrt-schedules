package station

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"github.com/kurapika12/mrt-schedules/common/client"
	"errors"
)

type Service interface {
	GetAllStations() (response []StationResponse, err error)
	CheckScheduleByStation(id string) (response []ScheduleResponse, err error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *service) GetAllStations() (response []StationResponse, err error) {
	url := "https://jakartamrt.co.id/id/val/stasiuns" 

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var stations []Station
	err = json.Unmarshal(byteResponse, &stations)

	for _, item := range stations {
		response = append(response, StationResponse(item))
	}

	return
}

func (s *service) CheckScheduleByStation(id string) (response []ScheduleResponse, err error) {
	url := "https://jakartamrt.co.id/id/val/stasiuns/"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var schedule []Schedule
	err = json.Unmarshal(byteResponse, &schedule)
	if err != nil {
		return
	}

	// schedule selection by id station
	var scheduleSelected Schedule
	for _, item := range schedule {
		if item.StationId == id {
			scheduleSelected = item
			break
		}
	}

	if scheduleSelected.StationId == "" {
		err = errors.New("station not found")
		return
	}

	response, err = ConvertDataToResponse(scheduleSelected)
	if err != nil {
		return
	}

	return
}


func ConvertDataToResponse(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		LebakBulusTripName = "Station Lebak Bulus Grab"
		BundaranHITripName = "Station Bundaran HI Bank DKI"
	)

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHI := schedule.ScheduleBundaranHI

	scheduleLebakBulusTimes, err := ConverScheduleToTimeFormat(scheduleLebakBulus)
	if err != nil {
		return
	}

	scheduleBundaranHITimes, err := ConverScheduleToTimeFormat(scheduleBundaranHI)
	if err != nil {
		return
	}	


	// convert to response
	for _, item := range scheduleLebakBulusTimes {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time:        item.Format("15:04"),
			})
	 	}
	}

	for _, item := range scheduleBundaranHITimes {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: BundaranHITripName,
				Time:        item.Format("15:04"),
			})
		}
	}
	return
}

func ConverScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	var (
		parsedTime time.Time
		schedules = strings.Split(schedule, ",")
	)

	for _, item := range schedules {
		trimedItem := strings.TrimSpace(item)
		if trimedItem == "" {
			continue
		}

		parsedTime, err = time.Parse("15:04", trimedItem)
		if err != nil {
			err = errors.New("invalid time format: " + trimedItem)
			return
		}	

		response = append(response, parsedTime)
	}
	return
}