package util

import (
	"encoding/json"
	"github.com/google/go-github/v32/github"
	"math/rand"
)

const (
	UID = "UID"
)

type MCache struct {
	Github github.User
	GithubEmails []github.UserEmail
	GithubAuth bool
	MixinId string
	MixinAuth bool
}

var longLetters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_sky")

func RandUp(n int) string{
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	arc := uint8(0)
	if _, err := rand.Read(b[:]); err != nil {
		return string(b)
	}
	for i, x := range b {
		arc = x & 63
		b[i] = longLetters[arc]
	}
	return string(b)
}

func UserToJson(user *github.User)(userJson string,err error)  {
	jsonBytes,err := json.Marshal(*user)
	userJson = string(jsonBytes)
	return
}

func JsonToUser(userJson string)(user *github.User,err error)  {
	user = &github.User{}
	err = json.Unmarshal([]byte(userJson), &user)
	return
}