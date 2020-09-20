package service

import (
	"claps-test/dao"
	"claps-test/util"
	"context"
	"github.com/fox-one/mixin-sdk-go"
)

func GetAssetByBotIdAndAssetId(botId string, assetId string) (asset *mixin.Asset, err *util.Err) {

	bot, err1 := dao.GetBotById(botId)
	if err1 != nil {
		err = util.NewErr(err, util.ErrDataBase, "通过botid获取bot信息失败")
	}

	client, err := CreateMixinClient(bot)
	if err != nil {
		return
	}
	asset, err1 = client.ReadAsset(context.Background(), assetId)
	if err1 != nil {
		err = util.NewErr(err, util.ErrThirdParty, "通过botid读取asset信息失败")
	}
	return
}
