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

const tokenExpireDuration = time.Hour * 72
type MyClaims struct {
	Uid int64	`json:"uid"`
	jwt.StandardClaims
}

/**
 * @Description: 检查是够绑定mixin,github一定是登录了,从数据库中查询问是否绑定mixin,经过该中间件ctx中一定有mixinId
 * @return gin.HandlerFunc
 */
func MixinAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var val interface{}
		var ok bool
		if val, ok = ctx.Get(util.UID); !ok {
			util.HandleResponse(ctx, util.NewErr(errors.New(""), util.ErrDataBase, "ctx get uid error"), nil)
			return
		}
		uid := val.(int64)

		//从数据库查询mixin_id
		mixinId, err := service.GetMixinIdByUserId(uid)
		if err != nil {
			util.HandleResponse(ctx, err, nil)
			ctx.Abort()
		}

		if mixinId == "" {
			util.HandleResponse(ctx, util.NewErr(errors.New("error"), util.ErrUnauthorized, "mixin unauthorized"), nil)
			ctx.Abort()
			return
		} else {
			//set ctx,next
			ctx.Set(util.MIXINID,mixinId)
		}
		ctx.Next()
	}
}

/**
 * @Description: 生成Token,uid=github.ID
 * @param uid
 * @return string
 * @return error
 */
func GenToken(uid int64) (string, error) {

	c := MyClaims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(), // 过期时间
			Issuer:    "sky",                                           // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(util.MySecret)

}

/**
 * @Description: 解析jwt为Myclaim
 * @param tokenString
 * @return *MyClaims
 * @return error
 */
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return util.MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

/**
 * @Description: 判断请求的Token情况,经过该中间件验证,ctx中一定有uid,但是不一定授权了mixin,自动查询数据库，填充mixin是否登录
 * @return func(c *gin.Context)
*/
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")

		//无Token,需要授权github
		if authHeader == "" {
			log.Debug("No Token")
			util.HandleResponse(ctx, util.NewErr(errors.New(""), util.ErrUnauthorized, "request have no token"), nil)
			ctx.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			util.HandleResponse(ctx, util.NewErr(errors.New(""), util.ErrUnauthorized, "authorization format error"), nil)
			ctx.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		claim, err1 := ParseToken(parts[1])
		if err1 != nil {
			util.HandleResponse(ctx, util.NewErr(err1, util.ErrUnauthorized, "invalid token"), nil)
			ctx.Abort()
			return
		}

		//set Key
		ctx.Set(util.UID, claim.Uid)
		//uid = randomUid不是githubId
		ctx.Next()
	}
}
