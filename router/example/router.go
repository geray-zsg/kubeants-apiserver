package example

import (
	"github.com/gin-gonic/gin"
	"kubeants.io/api"
)

type ExampleRouter struct{}

func (ExampleRouter) InitExample(r *gin.Engine) {
	group := r.Group("example")
	apiGroup := api.ApiGroupApp.ExampleApiGroup
	group.GET("ping", apiGroup.ExampleTest)
}
