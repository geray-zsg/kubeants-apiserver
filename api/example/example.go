package example

import (
	"github.com/gin-gonic/gin"
	"kubeants.com/response"
)

type ExampleTestApi struct{}

func (*ExampleTestApi) ExampleTest(c *gin.Context) {
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "pong",
	// })

	response.SuccessWithDetailed(c, "请求数据成功！", map[string]string{
		"message": "pong",
	})
}
