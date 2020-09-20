package router

import (
	"claps-test/controller"
	"claps-test/middleware"
	"claps-test/util"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.LoggerToFile())
	r.Use(util.Cors())

	r.GET("/_hc", func(ctx *gin.Context) {
		ctx.JSON(200, "ok")
	})

	// /api
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/authInfo", middleware.JWTAuthMiddleware(),controller.AuthInfo)
		apiGroup.GET("/oauth", controller.Oauth)

		apiGroup.GET("/bots/:botId/assets/:assetId", controller.Bot)

		// /api/projects
		projectsGroup := apiGroup.Group("projects")
		{
			projectsGroup.GET("/", controller.Projects)
			projectsGroup.GET("/:name", controller.Project)
			projectsGroup.GET("/:name/members", controller.ProjectMembers)
			projectsGroup.GET("/:name/transactions", controller.ProjectTransactions)

		}

		// /api/mixin
		mixinGroup := apiGroup.Group("/mixin")
		{
			mixinGroup.GET("/assets", controller.MixinAssets)
			mixinGroup.GET("/oauth",  controller.MixinOauth)
		}

		// /api/user
		userGroup := apiGroup.Group("/user")
		userGroup.Use(middleware.JWTAuthMiddleware())
		{
			userGroup.GET("/profile", controller.UserProfile)
			//查询所有币种的total和balance
			userGroup.GET("/assets", controller.UserAssets)
			//查询所有完成和未完成的记录
			userGroup.GET("/transfers", middleware.MixinAuthMiddleware(),controller.UserTransfer)
			//请求获得某个用户的捐赠信息的汇总,包括总金额和捐赠人数
			userGroup.GET("/donation", controller.UserDonation)
			//提现
			userGroup.GET("/withdraw", middleware.MixinAuthMiddleware(), controller.UserWithdraw)
		}

	}

	return r
}
