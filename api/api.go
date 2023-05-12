package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rng/engine"
)

type Api struct {
	engine *engine.Engine
}

func (api *Api) Run() {
	router := gin.Default()
	router.GET("/rand", func(c *gin.Context) {
		var result struct {
			Outcome []int `json:"outcome"`
		}
		result.Outcome = append(result.Outcome, int(api.engine.GetRand(100, 0)))
		fmt.Println(result)
		c.JSON(http.StatusOK, result)
	})
	router.Run(":8080")
}

func NewAPI(engine *engine.Engine) *Api {
	api := &Api{
		engine: engine,
	}
	return api
}
