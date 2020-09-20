package controller

import (
	"claps-test/service"
	"claps-test/util"
	"github.com/gin-gonic/gin"
)

func Bot(ctx *gin.Context) {

	asset, err := service.GetAssetByBotIdAndAssetId(ctx.Param("botId"), ctx.Param("assetId"))
	util.HandleResponse(ctx, err, asset)

}
