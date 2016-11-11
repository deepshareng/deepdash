package router

import (
	"net/http"

	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/handler/auth"
	"github.com/MISingularity/deepdash/handler/client"
	privateView "github.com/MISingularity/deepdash/handler/client/private"
	"github.com/MISingularity/deepdash/handler/resource"
	privateResource "github.com/MISingularity/deepdash/handler/resource/private"
	"github.com/gin-gonic/gin"

	"github.com/MISingularity/deepshare2/pkg/path"
)

func AddRouterByExistRule(router *gin.Engine) {
	curDir, _ := path.Getcurdir()
	router.StaticFS("/assets", http.Dir(curDir+"/../ds/assets"))
	router.LoadHTMLGlob(curDir + "/../ds/site/*")
	root := router.Group("/")
	{
		clientAPIRouter := root.Group("/")
		{
			//RouterAddRule(clientAPIRouter, client.RoutesClient)
			clientAPIRouter.GET("/auth/code", client.HandlerRedirectURIRetrieveToken)
			clientAPIRouter.POST("/authorization", client.HandlerAuthRequest)

			//RouterAddRule(clientAPIRouter, client.RoutesToken)
			clientAPIRouter.GET("/account/active", client.ActivateAccountHandler)
			clientAPIRouter.POST("/password/forgotpassword", client.ForgotPasswordHandler)
			clientAPIRouter.POST("/password/resetpassword", client.ResetPasswordHandler)

			//RouterAddRule(clientAPIRouter, client.RoutesSession)
			clientAPIRouter.GET("/session/loginselector", client.SessionLoginSelectorHandler)
			clientAPIRouter.GET("/session/logout", client.SessionLogoutHandler)
			clientAPIRouter.GET("/session/getuser", client.GetSessionUserInfoHandler)
			clientAPIRouter.POST("/session/checkusername", client.CheckUserNameHandler)
			clientAPIRouter.POST("/session/setappid/:appid", client.SetSessionAppIdHandler)

			//RouterAddRule(clientAPIRouter, client.RoutesGeneralView)
			clientAPIRouter.GET(handler.BaseURL.ForgetURL, client.ViewForgetHandler)
			clientAPIRouter.GET(handler.BaseURL.ResetURL, client.ViewResetHandler)
			//clientAPIRouter.GET(handler.BaseURL.RegisterURL, client.ViewRegisterHandler)

			//RouterAddRule(clientAPIRouter, resource.RoutesRegister)
			clientAPIRouter.PUT("/202cb962ac59075b964b07152d234b70/users", resource.PutUserHandler)
			clientAPIRouter.POST("/this-is-a-clandestine-resource/users", resource.PostUserHandler)
			clientAPIRouter.POST("/apply-try/users", resource.ApplyTryHandler)
			clientAPIRouter.POST("/this-is-a-clandestine-resource/resent-email", resource.ResentEmail)

		}
		clientGeneralViewRouter := root.Group("/")
		clientGeneralViewRouter.Use(auth.GeneralViewAuthentication)
		{
			//RouterAddRule(clientGeneralViewRouter, client.RoutesViewNeedRedirect)
			clientGeneralViewRouter.GET(handler.BaseURL.LoginURL, client.ViewLoginHandler)
			clientGeneralViewRouter.GET("unverified-email", client.ViewUnverifiedEmailHandler)

		}
		clientResourceViewRoute := router.Group("/")
		clientResourceViewRoute.Use(auth.ResourceViewAuthentication)
		{
			//RouterAddRule(clientResourceViewRoute, resource.RoutesDataView)
			clientResourceViewRoute.GET("/", resource.ViewIndexHandler)
		}

		resourceRouter := router.Group("/")
		resourceRouter.Use(auth.ResourceAuthentication)
		{

			//RouterAddRule(resourceRouter, resource.RoutesApp)
			resourceRouter.GET("/apps", resource.GetApplistHandler)
			resourceRouter.GET("/apps/:appid", resource.GetAppinfoHandler)
			resourceRouter.POST("/apps", resource.PostAppHandler)
			resourceRouter.POST("/apps/:appid/callback", resource.PostCallBackUrlHandler)
			resourceRouter.POST("/apps/:appid/uploadimage", resource.UploadImageHandler)
			resourceRouter.POST("/apps/:appid/uploadicon", resource.UploadIconHandler)
			resourceRouter.PUT("/apps/:appid", resource.PutAppHandler)
			resourceRouter.PUT("/apps/:appid/url", resource.PutAppUrlHandler)
			resourceRouter.DELETE("/apps/:appid", resource.DeleteAppHandler)

			//RouterAddRule(resourceRouter, resource.RoutesChannel)
			resourceRouter.GET("/apps/:appid/channels/statistics", resource.GetAppChannelsHandler)
			resourceRouter.POST("/apps/:appid/channels", resource.PostNewAppChannelHandler)
			resourceRouter.DELETE("/apps/:appid/channels/:channelname", resource.DeleteAppChannelHandler)
			resourceRouter.GET("/apps/:appid/channels", resource.GetAppChannelInfoHandler)
			resourceRouter.PUT("/apps/:appid/channelurl", resource.PutChannelUrlHandler)
			resourceRouter.GET("/apps/:appid/statistics", resource.GetChannelStatisticsHandler)

			resourceRouter.GET("/apps/:appid/smses", resource.GetSmsListHandler)
			resourceRouter.POST("/apps/:appid/smses", resource.PostSmsHandler)
			resourceRouter.DELETE("/apps/:appid/smses/:smsid", resource.DeleteSmsHandler)
			resourceRouter.PUT("/apps/:appid/smses/:smsid", resource.UpdateSmsHandler)

			//RouterAddRule(resourceRouter, resource.RoutesType)
			resourceRouter.GET("/apps/:appid/types", resource.TypesGetHandler)

			//RouterAddRule(resourceRouter, resource.RoutesEvent)
			resourceRouter.GET("/apps/:appid/events", resource.GetAppEventsHandler)

			//RouterAddRule(resourceRouter, resource.RoutesSelectedItems)
			resourceRouter.GET("/apps/:appid/selected_items", resource.GetSelectedItemsHandler)
			resourceRouter.PUT("/apps/:appid/selected_items", resource.PutSelectedItemsHandler)
		}
		privateResourceRouter := router.Group("/")
		privateResourceRouter.Use(auth.AdministratorResourceAuthentication)
		{
			//RouterAddRule(privateResourceRouter, privateResource.RoutesPrivateResource)
			privateResourceRouter.GET("/this-is-a-clandestine-resource/status", privateResource.GetStatusHandler)
			privateResourceRouter.GET("/this-is-a-clandestine-resource/total-status", privateResource.GetTotalStatusHandler)
			privateResourceRouter.GET("/this-is-a-clandestine-resource/apps", privateResource.GetSpecificAppStatus)
			privateResourceRouter.GET("/this-is-a-clandestine-resource/users", privateResource.GetUserList)
			privateResourceRouter.POST("/this-is-a-clandestine-resource/freeze", privateResource.FreezeUser)
			privateResourceRouter.POST("/this-is-a-clandestine-resource/unfreeze", privateResource.UnfreezeUser)
			privateResourceRouter.POST("/this-is-a-clandestine-resource/update/return-visit", privateResource.UpdateReturnVisitById)
			privateResourceRouter.POST("/this-is-a-clandestine-resource/user/add", privateResource.UserAdd)
			privateResourceRouter.POST("/this-is-a-clandestine-resource/permission/add", privateResource.PermissionAdd)
			privateResourceRouter.POST("/this-is-a-clandestine-resource/permission/remove", privateResource.PermissionRemove)
		}
		privateViewRouter := router.Group("/")
		privateViewRouter.Use(auth.AdministratorViewAuthentication)
		{
			//RouterAddRule(privateViewRouter, privateView.RoutesPrivateViews)
			privateViewRouter.GET("/private/status", privateView.PrivateStatusViewHandler)
			privateViewRouter.GET("/private/vital", privateView.PrivateVitalStatusViewHandler)
			privateViewRouter.GET("/private/add-user", privateView.PrivateUserAddViewHandler)
			privateViewRouter.GET("/private/permission", privateView.PermissionViewHandler)
		}

		// For website real-time pv
		root.GET("/statistics/pv-total", resource.TotalPvHandler)
	}
}
