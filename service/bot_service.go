package service

import (
	"claps-test/model"
	"claps-test/util"
	"context"
	"github.com/fox-one/mixin-sdk-go"
)

/**
 * @Description: 根据对应捐赠方式所对应的botId和要捐赠的assetId获取asset信息,主要需要Destination
 * @param botId
 * @param assetId
 * @return asset
 * @return err
 */
func GetAssetByBotIdAndAssetId(botId string, assetId string) (asset *mixin.Asset, err *util.Err) {

	bot, err1 := model.BOT.GetBotById(botId)
	if err1 != nil {
		err = util.NewErr(err, util.ErrDataBase, "通过botId获取bot信息失败")
	}

	client, err := CreateMixinClient(bot)
	if err != nil {
		return
	}
	asset, err1 = client.ReadAsset(context.Background(), assetId)
	if err1 != nil {
		err = util.NewErr(err, util.ErrThirdParty, "通过botId读取asset信息失败")
	}
	return
}
