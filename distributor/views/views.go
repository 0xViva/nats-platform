package views

import (
	"net/http"

	"encoding/json"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.GET("/today/:area", func(c *gin.Context) {
		area := c.Param("area")

		resp, err := http.Get("http://localhost:8080/api/today/" + area)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer resp.Body.Close()

		var prices []Price
		if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		render(c, PricesPage(prices, area))
	})
}
func render(c *gin.Context, cmp templ.Component) error {
	return cmp.Render(c.Request.Context(), c.Writer)
}
