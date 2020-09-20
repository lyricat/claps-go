package util

import (
	"encoding/gob"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/google/go-github/v32/github"
)

func RegisterType() {
	gob.Register(github.User{})
	gob.Register(mixin.User{})
	gob.Register(mixin.Client{})
}
