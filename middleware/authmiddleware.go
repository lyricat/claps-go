package middleware

import (
	"claps-test/service"
	"claps-test/util"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type MyClaims struct {
	//MixinId string `json:"mixin_id"`
	//GithubId string `json:"github_id"`
	Uid string `json:"uid"`
	jwt.StandardClaims
}


var MySecret = []byte("claps-dev")

const (
	MIXINID = "mixin_id"
	GITHUBID = "gtihub_id"
	TOKEN = "token"
)

type userInfo struct {
	mixin_id string
	github_id string
}

/*
功能:判断github是否已经授权
说明:经过了JWT中间件,一定有cache key
 */
func GithubAuthMiddleware()gin.HandlerFunc  {
	return func(ctx *gin.Context) {
		var val interface{}
		var ok bool
		if val,ok = ctx.Get(util.UID);!ok{
			util.HandleResponse(ctx,util.NewErr(errors.New(""),util.ErrDataBase,"ctx get uid error"),nil)
			return
		}
		uid := val.(string)

		mcache := &util.MCache{}
		err1 := util.Rdb.Get(uid,mcache)
		if err1 != nil{
			util.HandleResponse(ctx,util.NewErr(err1,util.ErrDataBase,"cache get error"),nil)
			return
		}

		//github未登录
		if !mcache.GithubAuth{
			util.HandleResponse(ctx,util.NewErr(err1,util.ErrUnauthorized,"github unauthorized"),nil)
			return
		}
		ctx.Next()
	}
}

/*
功能:检查是够绑定mixin
说明:github一定是登录了,从数据库中查询问是否绑定mixin,绑定则更新缓存
 */
func MixinAuthMiddleware()gin.HandlerFunc  {
	return func(ctx *gin.Context) {
		var val interface{}
		var ok bool
		if val,ok = ctx.Get(util.UID);!ok{
			util.HandleResponse(ctx,util.NewErr(errors.New(""),util.ErrDataBase,"ctx get uid error"),nil)
			return
		}
		uid := val.(string)

		mcache := &util.MCache{}
		err1 := util.Rdb.Get(uid,mcache)
		if err1 != nil{
			util.HandleResponse(ctx,util.NewErr(err1,util.ErrDataBase,"cache get error"),nil)
			return
		}

		if mcache.MixinAuth{
			ctx.Next()
		}

		//从数据库查询mixin_id
		mixin_id,err := service.GetMixinIdByUserId(*mcache.Github.ID)
		if err != nil{
			util.HandleResponse(ctx,err,nil)
			ctx.Abort()
		}

		if mixin_id == ""{
			util.HandleResponse(ctx,util.NewErr(err1,util.ErrUnauthorized,"mixin unauthorized"),nil)
			ctx.Abort()
			return
		}else {
			//set cache ,next
			mcache.MixinId = mixin_id
			mcache.MixinAuth = true
			err1 = util.Rdb.Replace(uid,mcache,-1)
			if err1 != nil{
				err = util.NewErr(errors.New("cache error"), util.ErrDataBase, "")
				util.HandleResponse(ctx, err, nil)
				return
			}
		}
		ctx.Next()
	}
}

/*
功能:生成Tokenm
说明:uid=github.ID
 */
func GenToken(uid string) (string, error) {

	c := MyClaims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(util.TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "sky",                               // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)

}

/*
功能:解析jwt为Myclaim
参数:jwt字符号
 */
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}



/*
功能:判断请求的Token情况
说明:经过该中间件验证,ctx中一定有cache的key,但是不一定授权了github
 */
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		log.Debug("authHeader = ",authHeader)

		//无Token,需要授权github
		if authHeader == "" {
			log.Debug("No Token")
			util.HandleResponse(c,util.NewErr(errors.New(""),util.ErrUnauthorized,"request have no token"),nil)
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			util.HandleResponse(c,util.NewErr(errors.New(""),util.ErrUnauthorized,"authorization format error"),nil)
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		claim, err1 := ParseToken(parts[1])
		if err1 != nil {
			util.HandleResponse(c,util.NewErr(err1,util.ErrUnauthorized,"invalid token"),nil)
			c.Abort()
			return
		}

		//set Key
		c.Set(util.UID,claim.Uid)
		c.Next()
	}
}


