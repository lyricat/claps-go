package util

import (
	"github.com/fox-one/mixin-sdk-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	EOS  = "6cfe566e-4aad-470b-8c9a-2fd35b49c68d"
	CNB  = "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	BTC  = "c6d0c728-2624-429b-8e0d-d9d19b6592fa"
	ETC  = "2204c1ee-0ea2-4add-bb9a-b3719cfff93a"
	XRP  = "23dfb5a5-5d7b-48b6-905f-3970e3176e27"
	XEM  = "27921032-f73e-434e-955f-43d55672ee31"
	ETH  = "43d61dcd-e413-450d-80b8-101d5e903357"
	DASH = "6472e7e3-75fd-48b6-b1dc-28d294ee1476"
	DOGE = "6770a1e5-6086-44d5-b60f-545f9d9e8ffd"
	LTC  = "76c802a2-7c88-447f-a93e-c29c9e5dd9c8"
	SC   = "990c4c29-57e9-48f6-9819-7d986ea44985"
	ZEN  = "a2c5d22b-62a2-4c13-b3f0-013290dbac60"
	ZEC  = "c996abc9-d94e-4494-b1cf-2a3fd3ac5714"
	BCH  = "fd11b6e3-0b87-41f1-a41f-f0e9b49e5bf0"
	USDT = "815b0b1a-2764-3736-8faa-42d694fa620a"
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

func CheckAsset(assetId *string) bool{
	switch *assetId{
	case DOGE,USDT,XEM,XRP,BCH,BTC,EOS,ETC:
		return true
	}
	return false
}