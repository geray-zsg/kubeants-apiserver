package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	success = iota
	fail
)

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": success,
		"msg":  "成功！",
	})
}
func SuccessWithMessage(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": success,
		"msg":  msg,
	})
}
func SuccessWithDetailed(c *gin.Context, msg string, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": success,
		"msg":  msg,
		"data": data,
	})
}

func Fail(c *gin.Context) {
	c.JSON(http.StatusFailedDependency, gin.H{
		"code": fail,
		"msg":  "失败！",
	})
}
func FailWithMessage(c *gin.Context, msg string) {
	c.JSON(http.StatusFailedDependency, gin.H{
		"code": fail,
		"msg":  msg,
	})
}
func FailWithDetailed(c *gin.Context, msg string, data any) {
	c.JSON(http.StatusFailedDependency, gin.H{
		"code": fail,
		"msg":  msg,
		"data": data,
	})
}
