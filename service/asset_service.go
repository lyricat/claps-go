package service

import (
	"claps-test/model"
	"claps-test/util"
)

/**
 * @Description: 通过assetId 将数据库中对应的asset信息读出来 暂时弃用
 * @param assetID
 * @return asset
 * @return err
 */
func GetAssetById(assetID string) (asset *model.Asset, err *util.Err) {
	asset, err1 := model.ASSET.GetAssetById(assetID)
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "通过id获取asset信息失败")
	}
	return
}

/**
 * @Description: 把数据库中asset表中的数据全读出来
 * @return assets
 * @return err
 */
func ListAssetsAllByDB() (assets *[]model.Asset, err *util.Err) {
	assets, err1 := model.ASSET.ListAssetsAllByDB()
	if err1 != nil {
		err = util.NewErr(err1, util.ErrDataBase, "获取全部assets信息失败")
	}
	return
}
