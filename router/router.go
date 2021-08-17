package router

import (
	"time"

	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/api/stream"
	"github.com/NubeDev/flow-framework/auth"
	"github.com/NubeDev/flow-framework/config"
	"github.com/NubeDev/flow-framework/database"
	"github.com/NubeDev/flow-framework/docs"
	"github.com/NubeDev/flow-framework/error"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin"
	"github.com/NubeDev/location"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Create creates the gin engine with all routes.
func Create(db *database.GormDatabase, vInfo *model.VersionInfo, conf *config.Configuration) (*gin.Engine, func()) {
	g := gin.New()

	g.Use(gin.Logger(), gin.Recovery(), error.Handler(), location.Default())
	g.NoRoute(error.NotFound())

	streamHandler := stream.New(time.Duration(conf.Server.Stream.PingPeriodSeconds)*time.Second, 15*time.Second, conf.Server.Stream.AllowedOrigins)
	authentication := auth.Auth{DB: db}
	messageHandler := api.MessageAPI{Notifier: streamHandler, DB: db}
	healthHandler := api.HealthAPI{DB: db}
	clientHandler := api.ClientAPI{
		DB:            db,
		ImageDir:      conf.UploadedImagesDir,
		NotifyDeleted: streamHandler.NotifyDeletedClient,
	}
	applicationHandler := api.ApplicationAPI{
		DB:       db,
		ImageDir: conf.UploadedImagesDir,
	}
	userChangeNotifier := new(api.UserChangeNotifier)
	userHandler := api.UserAPI{DB: db, PasswordStrength: conf.PassStrength, UserChangeNotifier: userChangeNotifier}
	networkHandler := api.NetworksAPI{
		DB: db,
	}
	deviceHandler := api.DeviceAPI{
		DB: db,
	}
	pointHandler := api.PointAPI{
		DB: db,
	}
	pluginManager, err := plugin.NewManager(db, conf.PluginsDir, g.Group("/plugin/:id/custom/"), streamHandler)
	if err != nil {
		panic(err)
	}
	pluginHandler := api.PluginAPI{
		Manager:  pluginManager,
		Notifier: streamHandler,
		DB:       db,
	}

	userChangeNotifier.OnUserDeleted(streamHandler.NotifyDeletedUser)
	userChangeNotifier.OnUserDeleted(pluginManager.RemoveUser)
	userChangeNotifier.OnUserAdded(pluginManager.InitializeForUserID)

	g.GET("/ip", api.Hostname) //TODO remove
	g.GET("/health", healthHandler.Health)
	g.GET("/swagger", docs.Serve)
	g.Static("/image", conf.UploadedImagesDir)
	g.GET("/docs", docs.UI)

	g.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		for header, value := range conf.Server.ResponseHeaders {
			ctx.Header(header, value)
		}
	})
	g.Use(cors.New(auth.CorsConfig(conf)))

	{
		g.GET("/plugin", authentication.RequireClient(), pluginHandler.GetPlugins)
		pluginRoute := g.Group("/plugin/", authentication.RequireClient())
		{
			pluginRoute.GET("/:id/config", pluginHandler.GetConfig)
			pluginRoute.POST("/:id/config", pluginHandler.UpdateConfig)
			pluginRoute.GET("/:id/display", pluginHandler.GetDisplay)
			pluginRoute.POST("/:id/enable", pluginHandler.EnablePlugin)
			pluginRoute.POST("/:id/disable", pluginHandler.DisablePlugin)
			pluginRoute.POST("/:id/network", pluginHandler.EnablePlugin) //TODO
		}
	}

	g.OPTIONS("/*any")

	// swagger:operation GET /version version getVersion
	//
	// Get version information.
	//
	// ---
	// produces: [application/json]
	// responses:
	//   200:
	//     description: Ok
	//     schema:
	//         $ref: "#/definitions/VersionInfo"
	g.GET("version", func(ctx *gin.Context) {
		ctx.JSON(200, vInfo)
	})

	g.Group("/").Use(authentication.RequireApplicationToken()).POST("/message", messageHandler.CreateMessage)

	clientAuth := g.Group("")
	{
		clientAuth.Use(authentication.RequireClient())
		app := clientAuth.Group("/application")
		{
			app.GET("", applicationHandler.GetApplications)

			app.POST("", applicationHandler.CreateApplication)

			app.POST("/:id/image", applicationHandler.UploadApplicationImage)

			app.PUT("/:id", applicationHandler.UpdateApplication)

			app.DELETE("/:id", applicationHandler.DeleteApplication)

			tokenMessage := app.Group("/:id/message")
			{
				tokenMessage.GET("", messageHandler.GetMessagesWithApplication)

				tokenMessage.DELETE("", messageHandler.DeleteMessageWithApplication)
			}
		}

		client := clientAuth.Group("/client")
		{
			client.GET("", clientHandler.GetClients)

			client.POST("", clientHandler.CreateClient)

			client.DELETE("/:id", clientHandler.DeleteClient)

			client.PUT("/:id", clientHandler.UpdateClient)
		}

		message := clientAuth.Group("/message")
		{
			message.GET("", messageHandler.GetMessages)

			message.DELETE("", messageHandler.DeleteMessages)

			message.DELETE("/:id", messageHandler.DeleteMessage)
		}

		clientAuth.GET("/stream", streamHandler.Handle)

		clientAuth.GET("current/user", userHandler.GetCurrentUser)

		clientAuth.POST("current/user/password", userHandler.ChangePassword)
	}

	authAdmin := g.Group("/user")
	{
		authAdmin.Use(authentication.RequireAdmin())

		authAdmin.GET("", userHandler.GetUsers)

		authAdmin.POST("", userHandler.CreateUser)

		authAdmin.DELETE("/:id", userHandler.DeleteUserByID)

		authAdmin.GET("/:id", userHandler.GetUserByID)

		authAdmin.POST("/:id", userHandler.UpdateUserByID)
	}

	control := g.Group("api")
	{
		control.Use(authentication.RequireAdmin())

		control.GET("", api.Hostname)
		control.GET("/networks", networkHandler.GetNetworks)
		control.POST("/networks", networkHandler.CreateNetwork)
		control.GET("/networks/:uuid", networkHandler.GetNetwork)
		control.PATCH("/networks/:uuid", networkHandler.UpdateNetwork)
		control.DELETE("/networks/:uuid", networkHandler.DeleteNetwork)

		control.GET("/devices", deviceHandler.GetDevices)
		control.POST("/devices", deviceHandler.CreateDevice)
		control.GET("/devices/:uuid", deviceHandler.GetDevice)
		control.PATCH("/devices/:uuid", deviceHandler.UpdateDevice)
		control.DELETE("/devices/:uuid", deviceHandler.DeleteDevice)


		control.GET("/points", pointHandler.GetPoints)
		control.POST("/points", pointHandler.CreatePoint)
		control.GET("/points/:uuid", pointHandler.GetPoint)
		control.PATCH("/points/:uuid", pointHandler.UpdatePoint)
		control.DELETE("/points/:uuid", pointHandler.DeletePoint)

	}

	return g, streamHandler.Close
}
