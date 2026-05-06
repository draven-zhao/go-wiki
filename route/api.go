package route

import (
	"go-wiki/controller"

	"github.com/gin-gonic/gin"
)

func Api(r *gin.Engine) {
	r.POST("/query", controller.QueryHandler)
}
