package router

import (
	"claps-test/controller"
	"claps-test/middleware"
	"github.com/gin-gonic/gin"
)

/**
 * @Description: 注册所有路由
 * @param r
 * @return *gin.Engine
 */
func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.LoggerToFile())
	r.Use(middleware.Cors())
	r.Use(gin.Recovery())

	r.GET("/_hc", func(ctx *gin.Context) {
		ctx.JSON(200, "ok")
	})

	// /api
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/authInfo", middleware.JWTAuthMiddleware(), controller.AuthInfo)
		apiGroup.GET("/environments", controller.Environments)

		apiGroup.GET("/oauth", controller.Oauth)

		apiGroup.GET("/bots/:botId/assets/:assetId", controller.Bot)

		// /api/projects
		projectsGroup := apiGroup.Group("projects")
		{
			projectsGroup.GET("/", controller.Projects)
			projectsGroup.GET("/:id", controller.ProjectById)
			projectsGroup.GET("/:id/members", controller.ProjectMembers)
			projectsGroup.GET("/:id/transactions", controller.ProjectTransactions)
			projectsGroup.GET("/:id/svg", controller.ProjectSvg)

		}

		// /api/mixin
		mixinGroup := apiGroup.Group("/mixin")
		{
			mixinGroup.GET("/assets", controller.MixinAssets)
			mixinGroup.GET("/oauth", middleware.JWTAuthMiddleware(), controller.MixinOauth)
		}

		// /api/user
		userGroup := apiGroup.Group("/user")
		userGroup.Use(middleware.JWTAuthMiddleware())
		{
			userGroup.GET("/profile", controller.UserProfile)
			//查询所有币种的total和balance
			userGroup.GET("/assets", controller.UserAssets)
			//查询所有完成和未完成的记录
			userGroup.GET("/transfers", middleware.MixinAuthMiddleware(), controller.UserTransfer)
			//请求获得某个用户的捐赠信息的汇总,包括总金额和捐赠人数
			userGroup.GET("/donation", controller.UserDonation)
			//提现
			userGroup.POST("/withdraw", middleware.MixinAuthMiddleware(), controller.UserWithdraw)
			userGroup.POST("/withdrawalWay", middleware.MixinAuthMiddleware(), controller.UserWithdrawalWay)
		}

	}

	return r
}
