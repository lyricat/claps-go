package util

import (
	"github.com/fox-one/mixin-sdk-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	//case DOGE,XRP,BCH,BTC,ETH:
	BTC  = "c6d0c728-2624-429b-8e0d-d9d19b6592fa"
	ETH  = "43d61dcd-e413-450d-80b8-101d5e903357"
	XRP  = "23dfb5a5-5d7b-48b6-905f-3970e3176e27"
	DOGE = "6770a1e5-6086-44d5-b60f-545f9d9e8ffd"
	BCH  = "fd11b6e3-0b87-41f1-a41f-f0e9b49e5bf0"
	//EOS的存币地址与其它的币有些不同，它由两部分组成： account_name and account tag, 如果你向Mixin Network存入EOS，你需要填两项数据： account name 是eoswithmixin,备注里输入你的account_tag,比如0aa2b00fad2c69059ca1b50de2b45569.

)

var MixinClient *mixin.Client

func InitMixin() *mixin.Client {
	s := &mixin.Keystore{
		ClientID:   viper.GetString("client_id"),
		SessionID:  viper.GetString("session_id"),
		PrivateKey: viper.GetString("private_key"),
		PinToken:   viper.GetString("pin_token"),
	}

	var err error
	MixinClient, err = mixin.NewFromKeystore(s)
	if err != nil {
		log.Error(err.Error())
	}
	return MixinClient
}

func GetMixin() *mixin.Client {
	return MixinClient
}

/**
 * @Description: 默认支持如下5中币种
 * @param assetId
 * @return bool
 */
func CheckAsset(assetId *string) bool {
	switch *assetId {
	case DOGE, XRP, BCH, BTC, ETH:
		return true
	}
	return false
}
