package util

import (
	"github.com/google/go-github/v32/github"
	"crypto/rand"
)

const (
	UID 	= 	"UID"
	MIXINID = 	"MIXINID"
)

type MCache struct {
	Github       github.User
	GithubEmails []github.UserEmail
	GithubAuth   bool
	MixinId      string
	MixinAuth    bool
}

var longLetters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_sky")

/**
 * @Description: 生成一个长度为n的随机字符串
 * @param n
 * @return string
 */
func RandUp(n int) string {
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
