package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Price struct {
	NOKPerKWh float64 `json:"NOK_per_kWh"`
	TimeStart string  `json:"time_start"`
	TimeEnd   string  `json:"time_end"`
}

func Setup(r *gin.Engine) {
	r.GET("/api/today/:area", func(c *gin.Context) {
		// Prisområde: NO1, NO2, NO3, NO4 eller NO5
		area := c.Param("area")

		now := time.Now()
		year := now.Year()
		month := now.Month()
		day := now.Day()

		url := fmt.Sprintf(
			"https://www.hvakosterstrommen.no/api/v1/prices/%d/%02d-%02d_%s.json",
			year,
			month,
			day,
			area,
		)

		resp, err := http.Get(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		var raw []Price
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		out := make([]Price, 0, len(raw))

		for _, p := range raw {
			start, err1 := time.Parse(time.RFC3339, p.TimeStart)
			end, err2 := time.Parse(time.RFC3339, p.TimeEnd)

			if err1 != nil || err2 != nil {
				continue
			}

			out = append(out, Price{
				NOKPerKWh: p.NOKPerKWh,
				TimeStart: start.Format("15:04"),
				TimeEnd:   end.Format("15:04"),
			})
		}

		c.JSON(http.StatusOK, out)
	})
}
