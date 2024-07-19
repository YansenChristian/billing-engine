package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handle(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}
