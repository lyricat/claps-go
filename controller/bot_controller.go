package controller

import (
	"claps-test/service"
	"claps-test/util"
	"github.com/gin-gonic/gin"
)

/**
 * @Description: 通过botId和assetId获取要捐赠的虚拟货币的地址
 * @param ctx
 */
func Bot(ctx *gin.Context) {

	asset, err := service.GetAssetByBotIdAndAssetId(ctx.Param("botId"), ctx.Param("assetId"))
	util.HandleResponse(ctx, err, asset)

}
