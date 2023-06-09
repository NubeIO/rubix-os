package router

import (
	"github.com/NubeDev/location"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	authconstants "github.com/NubeIO/nubeio-rubix-lib-auth-go/constants"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/auth"
	"github.com/NubeIO/rubix-os/config"
	"github.com/NubeIO/rubix-os/constants"
	"github.com/NubeIO/rubix-os/database"
	"github.com/NubeIO/rubix-os/eventbus"
	"github.com/NubeIO/rubix-os/global"
	"github.com/NubeIO/rubix-os/installer"
	"github.com/NubeIO/rubix-os/logger"
	"github.com/NubeIO/rubix-os/module"
	"github.com/NubeIO/rubix-os/nerrors"
	"github.com/NubeIO/rubix-os/plugin"
	"github.com/NubeIO/rubix-os/rubixregistry"
	"github.com/NubeIO/rubix-os/services/appstore"
	"github.com/NubeIO/rubix-os/services/system"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func Create(db *database.GormDatabase, conf *config.Configuration, scheduler *gocron.Scheduler,
	systemCtl *systemctl.SystemCtl, system_ *system.System, registry *rubixregistry.RubixRegistry) *gin.Engine {
	engine := gin.New()
	engine.Use(logger.GinMiddlewareLogger(), gin.Recovery(), nerrors.Handler(), location.Default())
	engine.NoRoute(nerrors.NotFoundHandler())
	eventBus := eventbus.NewService(eventbus.GetBus())
	global.Installer = installer.New(&installer.Installer{})
	healthHandler := api.HealthAPI{DB: db}

	authHandler := api.AuthAPI{DB: db}
	handleAuth := func(c *gin.Context) { c.Next() }
	if *conf.Auth {
		handleAuth = authHandler.HandleAuth(false, authconstants.UserRole)
	}
	handleAuthWithMember := authHandler.HandleAuth(false, authconstants.UserRole, constants.MemberRole)
	handleAuthWithMemberHostLevel := authHandler.HandleAuth(true, authconstants.UserRole, constants.MemberRole)

	apiRoutesTemp := engine.Group("/api", handleAuth) // TODO: remove this one and use the same one
	// http://localhost:1660/api/plugins/api/system/schema/json/device
	pluginManager, err := plugin.NewManager(db, conf.GetAbsPluginsDir(), apiRoutesTemp.Group("/plugins/api"))
	if err != nil {
		log.Error(err)
		panic(err)
	}

	modules, err := module.ReLoadModulesWithDir(config.Get().GetAbsModulesDir(), apiRoutesTemp.Group("/modules"))
	if err != nil {
		log.Error(err)
		panic(err)
	}
	db.Modules = modules

	db.PluginManager = pluginManager
	pluginHandler := api.PluginAPI{
		Manager: pluginManager,
		Modules: modules,
		DB:      db,
	}
	networkHandler := api.NetworksAPI{
		DB:     db,
		Bus:    eventBus,
		Plugin: pluginManager,
	}
	deviceHandler := api.DeviceAPI{
		DB: db,
	}
	pointHandler := api.PointAPI{
		DB: db,
	}
	pointHistoryHandler := api.PointHistoryAPI{
		DB: db,
	}
	jobHandler := api.JobAPI{
		DB: db,
	}
	dbGroup := api.DatabaseAPI{
		DB: db,
	}
	integrationHandler := api.IntegrationAPI{
		DB: db,
	}
	mqttHandler := api.MqttConnectionAPI{
		DB: db,
	}
	schHandler := api.ScheduleAPI{
		DB: db,
	}
	thingHandler := api.ThingAPI{}

	tagHandler := api.TagAPI{
		DB: db,
	}

	deviceInfoHandler := api.DeviceInfoAPI{
		RubixRegistry: registry,
	}
	systemctlHandler := api.SystemctlAPI{
		SystemCtl: systemCtl,
	}
	syscallHandler := api.SyscallAPI{}
	dateHandler := api.DateAPI{
		System: system_,
	}
	networkingHandler := api.NetworkingAPI{
		System: system_,
	}
	fileHandler := api.FileAPI{
		FileMode: 0755,
	}
	dirHandler := api.DirApi{
		FileMode: 0755,
	}
	zipHandler := api.ZipApi{
		FileMode: 0755,
	}
	streamLogHandler := api.StreamLogApi{
		DB: db,
	}
	snapshotHandler := api.SnapshotAPI{
		SystemCtl:     systemCtl,
		FileMode:      0755,
		RubixRegistry: registry,
	}
	restartJobHandler := api.RestartJobApi{
		SystemCtl: systemCtl,
		Scheduler: scheduler,
		FileMode:  0755,
	}
	makeStore, _ := appstore.New(&appstore.Store{})
	appStoreHandler := api.AppStoreApi{
		Store: makeStore,
	}
	pluginStoreHandler := api.PluginStoreApi{
		Store: makeStore,
	}
	edgeBiosEdgeHandler := api.EdgeBiosEdgeApi{
		DB: db,
	}
	edgeAppHandler := api.EdgeAppApi{
		DB: db,
	}
	edgePluginHandler := api.EdgePluginApi{
		DB: db,
	}
	edgeConfigHandler := api.EdgeConfigApi{
		DB: db,
	}
	edgeSnapshotHandler := api.EdgeSnapshotApi{
		DB:       db,
		FileMode: 0755,
	}
	snapshotCreatLogHandler := api.SnapshotCreateLogAPI{
		DB: db,
	}
	snapshotRestoreLogHandler := api.SnapshotRestoreLogAPI{
		DB: db,
	}
	locationHandler := api.LocationAPI{
		DB: db,
	}
	groupHandler := api.GroupAPI{
		DB: db,
	}
	hostHandler := api.HostAPI{
		DB: db,
	}
	hostCommentHandler := api.HostCommentAPI{
		DB: db,
	}
	hostTagHandler := api.HostTagAPI{
		DB: db,
	}
	viewSettingHandler := api.ViewSettingAPI{
		DB: db,
	}
	viewTemplateHandler := api.ViewTemplateAPI{
		DB: db,
	}
	viewTemplateWidgetHandler := api.ViewTemplateWidgetAPI{
		DB: db,
	}
	viewHandler := api.ViewAPI{
		DB: db,
	}
	viewWidgetHandler := api.ViewWidgetAPI{
		DB: db,
	}
	systemHandler := api.SystemAPI{
		System:    system_,
		Scheduler: scheduler,
		FileMode:  0755,
	}
	alertHandler := api.AlertAPI{
		DB: db,
	}
	memberHandler := api.MemberAPI{
		DB: db,
	}
	memberDeviceHandler := api.MemberDeviceAPI{
		DB: db,
	}
	teamHandler := api.TeamAPI{
		DB: db,
	}
	cloudEdgeCloneHandler := api.CloneEdgeApi{
		DB: db,
	}
	ticketHandler := api.TicketAPI{
		DB: db,
	}
	ticketCommentHandler := api.TicketCommentAPI{
		DB: db,
	}
	fcmServerHandler := api.FcmServerAPI{
		DB: db,
	}
	userHandler := api.UserAPI{}
	tokenHandler := api.TokenAPI{}

	wiresProxyHandler := api.WiresProxyAPI{}
	chirpProxyHandler := api.ChirpProxyAPI{}
	hostProxyHandler := api.HostProxyAPI{
		DB: db,
	}
	dbGroup.SyncTopics()

	// for the custom plugin endpoints you need to use the plugin token
	engine.GET("/api/system/ping", healthHandler.Health)
	engine.POST("/api/users/login", userHandler.Login)
	engine.Static("/image", conf.GetAbsUploadedImagesDir())
	engine.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json") // if you comment it out it will detected as text on proxy-handlers
		for header, value := range conf.Server.ResponseHeaders {
			ctx.Header(header, value)
		}
	})

	apiProxyWiresRoutes := engine.Group("/wires", handleAuth)
	apiProxyWiresRoutes.Any("/*proxyPath", wiresProxyHandler.WiresProxy) // EDGE-WIRES PROXY
	apiProxyChirpRoutes := engine.Group("/chirp", handleAuth)
	apiProxyChirpRoutes.Any("/*proxyPath", chirpProxyHandler.ChirpProxy) // CHIRP-STACK PROXY

	apiProxyHostRoutes := engine.Group("/proxy")
	apiProxyHostRoutes.Use(auth.HostProxyOptions())

	apiProxyHostRoutesAuth := apiProxyHostRoutes.Group("/", handleAuth)
	apiProxyHostRoutesAuth.Any("/*proxyPath", hostProxyHandler.HostProxy)

	engine.Use(cors.New(auth.CorsConfig()))
	engine.Use(auth.HostProxyOptions())
	appApiRoutes := engine.Group("/api/apps")
	{
		memberRoutes := appApiRoutes.Group("/members")
		{
			memberRoutes.POST("", memberHandler.CreateMember)
			memberRoutes.POST("/login", memberHandler.Login)
			memberRoutes.GET("/check_username/:username", memberHandler.CheckUsername)
			memberRoutes.GET("/check_email/:email", memberHandler.CheckEmail)
		}

		ownMemberRoutes := appApiRoutes.Group("/o/members")
		{
			ownMemberRoutes.GET("", memberHandler.GetMember)
			ownMemberRoutes.PATCH("", memberHandler.UpdateMember)
			ownMemberRoutes.DELETE("", memberHandler.DeleteMember)
			ownMemberRoutes.POST("/change_password", memberHandler.ChangePassword)
			ownMemberRoutes.POST("/refresh_token", memberHandler.RefreshToken)
			ownMemberRoutes.GET("/sidebars", memberHandler.GetMemberSidebars)

			memberDevicesRoutes := ownMemberRoutes.Group("/devices")
			{
				memberDevicesRoutes.GET("", memberDeviceHandler.GetMemberDevices)
				memberDevicesRoutes.GET("/:device_id", memberDeviceHandler.GetMemberDevice)
				memberDevicesRoutes.POST("", memberDeviceHandler.CreateMemberDevice)
				memberDevicesRoutes.PATCH("/:device_id", memberDeviceHandler.UpdateMemberDevice)
				memberDevicesRoutes.DELETE("/:device_id", memberDeviceHandler.DeleteMemberDevice)
			}
		}
	}
	apiRoutes := engine.Group("/api", handleAuth)
	{
		requireClientsGroupRoutes := apiRoutes.Group("")
		{
			plugins := requireClientsGroupRoutes.Group("/plugins")
			{
				plugins.GET("", pluginHandler.GetPlugins)
				plugins.GET("/:uuid", pluginHandler.GetPlugin)
				plugins.GET("/config/:uuid", pluginHandler.GetConfig)
				plugins.POST("/config/:uuid", pluginHandler.UpdateConfig)
				plugins.GET("/display/:uuid", pluginHandler.GetDisplay) // todo: remove
				plugins.POST("/enable/:uuid", pluginHandler.EnablePluginByUUID)
				plugins.POST("/restart/:uuid", pluginHandler.RestartPlugin)
				plugins.GET("/path/:path", pluginHandler.GetPluginByPath)
			}
		}

		pointHistoryRoutes := apiRoutes.Group("/histories/points")
		{
			pointHistoryRoutes.GET("", pointHistoryHandler.GetPointHistories)
			pointHistoryRoutes.GET("/:point_uuid", pointHistoryHandler.GetPointHistoriesByPointUUID)
			pointHistoryRoutes.GET("/:point_uuid/one", pointHistoryHandler.GetLatestPointHistoryByPointUUID)
			pointHistoryRoutes.POST("/point_uuids", pointHistoryHandler.GetPointHistoriesByPointUUIDs)
			pointHistoryRoutes.GET("/sync", pointHistoryHandler.GetPointHistoriesForSync)
			pointHistoryRoutes.DELETE("/:point_uuid", pointHistoryHandler.DeletePointHistoriesByPointUUID)
		}

		networkRoutes := apiRoutes.Group("/networks")
		{
			networkRoutes.GET("", networkHandler.GetNetworks)
			networkRoutes.POST("", networkHandler.CreateNetwork)
			networkRoutes.GET("/:uuid", networkHandler.GetNetwork)
			networkRoutes.GET("/plugin/:name", networkHandler.GetNetworkByPluginName)
			networkRoutes.GET("/plugin/all/:name", networkHandler.GetNetworksByPluginName)
			networkRoutes.GET("/name/:name", networkHandler.GetNetworkByName)
			networkRoutes.PATCH("/:uuid", networkHandler.UpdateNetwork)
			networkRoutes.DELETE("/:uuid", networkHandler.DeleteNetwork)
			networkRoutes.DELETE("/one/args", networkHandler.DeleteOneNetworkByArgs)
			networkRoutes.DELETE("/name/:name", networkHandler.DeleteNetworkByName)
			networkRoutes.PUT("/meta_tags/uuid/:uuid", networkHandler.CreateNetworkMetaTags)
		}

		deviceRoutes := apiRoutes.Group("/devices")
		{
			deviceRoutes.GET("", deviceHandler.GetDevices)
			deviceRoutes.POST("", deviceHandler.CreateDevice)
			deviceRoutes.GET("/:uuid", deviceHandler.GetDevice)
			deviceRoutes.GET("/one/args", deviceHandler.GetOneDeviceByArgs)
			deviceRoutes.GET("/name/:network_name/:device_name", deviceHandler.GetDeviceByName)
			deviceRoutes.PATCH("/:uuid", deviceHandler.UpdateDevice)
			deviceRoutes.DELETE("/:uuid", deviceHandler.DeleteDevice)
			deviceRoutes.DELETE("/one/args", deviceHandler.DeleteOneDeviceByArgs)
			deviceRoutes.DELETE("/name/:network_name/:device_name", deviceHandler.DeleteDeviceByName)
			deviceRoutes.PUT("/meta_tags/uuid/:uuid", deviceHandler.CreateDeviceMetaTags)
		}

		pointRoutes := apiRoutes.Group("/points")
		{
			pointRoutes.GET("", pointHandler.GetPoints)
			pointRoutes.GET("/bulk/uuids", pointHandler.GetPointsBulkUUIs)
			pointRoutes.POST("/bulk", pointHandler.GetPointsBulk)
			pointRoutes.GET("/:uuid", pointHandler.GetPoint)
			pointRoutes.GET("/name", pointHandler.GetPointByNameArgs) // TODO remove
			pointRoutes.GET("/name/:network_name/:device_name/:point_name", pointHandler.GetPointByName)
			pointRoutes.GET("/one/args", pointHandler.GetOnePointByArgs)
			pointRoutes.POST("", pointHandler.CreatePoint)
			pointRoutes.PATCH("/:uuid", pointHandler.UpdatePoint)
			pointRoutes.PATCH("/write/:uuid", pointHandler.PointWrite)
			pointRoutes.DELETE("/:uuid", pointHandler.DeletePoint)
			pointRoutes.DELETE("/one/args", pointHandler.DeleteOnePointByArgs)
			pointRoutes.DELETE("/name/:network_name/:device_name/:point_name", pointHandler.DeletePointByName)
			pointRoutes.PATCH("/name", pointHandler.PointWriteByNameArgs) // TODO remove
			pointRoutes.PATCH("/name/:network_name/:device_name/:point_name", pointHandler.PointWriteByName)
			pointRoutes.PUT("/meta_tags/uuid/:uuid", pointHandler.CreatePointMetaTags)
			pointRoutes.GET("/with_parent/:uuid", pointHandler.GetPointWithParent)
		}

		jobRoutes := apiRoutes.Group("/jobs")
		{
			jobRoutes.GET("", jobHandler.GetJobs)
			jobRoutes.POST("", jobHandler.CreateJob)
			jobRoutes.GET("/:uuid", jobHandler.GetJob)
			jobRoutes.PATCH("/:uuid", jobHandler.UpdateJob)
			jobRoutes.DELETE("/:uuid", jobHandler.DeleteJob)
		}

		integrationRoutes := apiRoutes.Group("/integrations")
		{
			integrationRoutes.GET("", integrationHandler.GetIntegrations)
			integrationRoutes.POST("", integrationHandler.CreateIntegration)
			integrationRoutes.GET("/:uuid", integrationHandler.GetIntegration)
			integrationRoutes.PATCH("/:uuid", integrationHandler.UpdateIntegration)
			integrationRoutes.DELETE("/:uuid", integrationHandler.DeleteIntegration)
		}

		mqttClientRoutes := apiRoutes.Group("/localmqtt/clients")
		{
			mqttClientRoutes.GET("", mqttHandler.GetMqttConnectionsList)
			mqttClientRoutes.POST("", mqttHandler.CreateMqttConnection)
			mqttClientRoutes.GET("/:uuid", mqttHandler.GetMqttConnection)
			mqttClientRoutes.PATCH("/:uuid", mqttHandler.UpdateMqttConnection)
			mqttClientRoutes.DELETE("/:uuid", mqttHandler.DeleteMqttConnection)
		}

		schRoutes := apiRoutes.Group("/schedules")
		{
			schRoutes.GET("", schHandler.GetSchedules)
			schRoutes.GET("/:uuid", schHandler.GetSchedule)
			schRoutes.GET("/result", schHandler.GetSchedulesResult)
			schRoutes.GET("/result/:uuid", schHandler.GetScheduleResult)
			schRoutes.POST("", schHandler.CreateSchedule)
			schRoutes.GET("/one/args", schHandler.GetOneScheduleByArgs)
			schRoutes.PATCH("/:uuid", schHandler.UpdateSchedule)
			schRoutes.PATCH("/write/:uuid", schHandler.ScheduleWrite)
			schRoutes.DELETE("/:uuid", schHandler.DeleteSchedule)
		}

		thingRoutes := apiRoutes.Group("/things")
		{
			thingRoutes.GET("/class", thingHandler.ThingClass)
			thingRoutes.GET("/writers/actions", thingHandler.WriterActions)
			thingRoutes.GET("/units", thingHandler.ThingUnits)
		}

		tagRoutes := apiRoutes.Group("/tags")
		{
			tagRoutes.GET("", tagHandler.GetTags)
			tagRoutes.POST("", tagHandler.CreateTag)
			tagRoutes.GET("/:tag", tagHandler.GetTag)
			tagRoutes.DELETE(":tag", tagHandler.DeleteTag)
		}

		deviceInfoRoutes := apiRoutes.Group("/system")
		{
			deviceInfoRoutes.GET("/device", deviceInfoHandler.GetDeviceInfo)
			deviceInfoRoutes.POST("/scanner", systemHandler.RunScanner)
			deviceInfoRoutes.GET("/network_interfaces", systemHandler.GetNetworkInterfaces)
			deviceInfoRoutes.POST("/reboot", systemHandler.RebootHost)
			deviceInfoRoutes.GET("/reboot/job", systemHandler.GetRebootHostJob)
			deviceInfoRoutes.PUT("/reboot/job", systemHandler.UpdateRebootHostJob)
			deviceInfoRoutes.DELETE("/reboot/job", systemHandler.DeleteRebootHostJob)

			deviceInfoRoutes.GET("/time", systemHandler.HostTime)
			deviceInfoRoutes.GET("/info", systemHandler.GetSystem)
			deviceInfoRoutes.GET("/usage", systemHandler.GetMemoryUsage)
			deviceInfoRoutes.GET("/memory", systemHandler.GetMemory)
			deviceInfoRoutes.GET("/processes", systemHandler.GetTopProcesses) // /processes?sort=cpu&count=3
			deviceInfoRoutes.GET("/swap", systemHandler.GetSwap)
			deviceInfoRoutes.GET("/disc", systemHandler.DiscUsage)
			deviceInfoRoutes.GET("/disc/pretty", systemHandler.DiscUsagePretty)
		}

		userRoutes := apiRoutes.Group("/users")
		{
			userRoutes.PUT("", userHandler.UpdateUser)
			userRoutes.GET("", userHandler.GetUser)
		}

		tokenRoutes := apiRoutes.Group("/tokens")
		{
			tokenRoutes.GET("", tokenHandler.GetTokens)
			tokenRoutes.GET("/:uuid", tokenHandler.GetToken)
			tokenRoutes.POST("/generate", tokenHandler.GenerateToken)
			tokenRoutes.PUT("/:uuid/block", tokenHandler.BlockToken)
			tokenRoutes.PUT("/:uuid/regenerate", tokenHandler.RegenerateToken)
			tokenRoutes.DELETE("/:uuid", tokenHandler.DeleteToken)
		}

		memberRoutes := apiRoutes.Group("/members")
		{
			memberRoutes.GET("", memberHandler.GetMembers)
			memberRoutes.GET("/:uuid", memberHandler.GetMemberByUUID)
			memberRoutes.DELETE("/:uuid", memberHandler.DeleteMemberByUUID)
			memberRoutes.PATCH("/:uuid", memberHandler.UpdateMemberByUUID)
			memberRoutes.POST("/:uuid/change_password", memberHandler.ChangeMemberPassword)
			memberRoutes.GET("/username/:username", memberHandler.GetMemberByUsername)
			memberRoutes.POST("/verify/:username", memberHandler.VerifyMember)
		}

		systemctlRoutes := apiRoutes.Group("/systemctl")
		{
			systemctlRoutes.POST("/enable", systemctlHandler.SystemCtlEnable)
			systemctlRoutes.POST("/disable", systemctlHandler.SystemCtlDisable)
			systemctlRoutes.GET("/show", systemctlHandler.SystemCtlShow)
			systemctlRoutes.POST("/start", systemctlHandler.SystemCtlStart)
			systemctlRoutes.GET("/status", systemctlHandler.SystemCtlStatus)
			systemctlRoutes.POST("/stop", systemctlHandler.SystemCtlStop)
			systemctlRoutes.POST("/reset-failed", systemctlHandler.SystemCtlResetFailed)
			systemctlRoutes.POST("/daemon-reload", systemctlHandler.SystemCtlDaemonReload)
			systemctlRoutes.POST("/restart", systemctlHandler.SystemCtlRestart)
			systemctlRoutes.POST("/mask", systemctlHandler.SystemCtlMask)
			systemctlRoutes.POST("/unmask", systemctlHandler.SystemCtlUnmask)
			systemctlRoutes.GET("/state", systemctlHandler.SystemCtlState)
			systemctlRoutes.GET("/is-enabled", systemctlHandler.SystemCtlIsEnabled)
			systemctlRoutes.GET("/is-active", systemctlHandler.SystemCtlIsActive)
			systemctlRoutes.GET("/is-running", systemctlHandler.SystemCtlIsRunning)
			systemctlRoutes.GET("/is-failed", systemctlHandler.SystemCtlIsFailed)
			systemctlRoutes.GET("/is-installed", systemctlHandler.SystemCtlIsInstalled)
		}

		syscallRoutes := apiRoutes.Group("/syscall")
		{
			syscallRoutes.POST("/unlink", syscallHandler.SyscallUnlink)
			syscallRoutes.POST("/link", syscallHandler.SyscallLink)
		}

		timeRoutes := apiRoutes.Group("/time")
		{
			timeRoutes.GET("", dateHandler.SystemTime)
			timeRoutes.POST("", dateHandler.SetSystemTime)
			timeRoutes.POST("ntp/enable", dateHandler.NTPEnable)
			timeRoutes.POST("ntp/disable", dateHandler.NTPDisable)
		}

		timeZoneRoutes := apiRoutes.Group("/timezone")
		{
			timeZoneRoutes.GET("", dateHandler.GetHardwareTZ)
			timeZoneRoutes.POST("", dateHandler.UpdateTimezone)
			timeZoneRoutes.GET("/list", dateHandler.GetTimeZoneList)
			timeZoneRoutes.POST("/config", dateHandler.GenerateTimeSyncConfig)
		}

		networkingRoutes := apiRoutes.Group("/networking")
		{
			networkingRoutes.GET("", networkingHandler.Networking)
			networkingRoutes.GET("/interfaces", networkingHandler.GetInterfacesNames)
			networkingRoutes.GET("/internet", networkingHandler.InternetIP)

			networkingNetworkRoutes := networkingRoutes.Group("networks")
			{
				networkingNetworkRoutes.POST("/restart", networkingHandler.RestartNetworking)
			}

			networkingInterfaceRoutes := networkingRoutes.Group("interfaces")
			{
				networkingInterfaceRoutes.POST("/exists", networkingHandler.DHCPPortExists)
				networkingInterfaceRoutes.POST("/auto", networkingHandler.DHCPSetAsAuto)
				networkingInterfaceRoutes.POST("/static", networkingHandler.DHCPSetStaticIP)
				networkingInterfaceRoutes.POST("/reset", networkingHandler.InterfaceUpDown)
				networkingInterfaceRoutes.POST("/pp", networkingHandler.InterfaceUp)
				networkingInterfaceRoutes.POST("/down", networkingHandler.InterfaceDown)
			}

			networkingFirewallRoutes := networkingRoutes.Group("/firewall")
			{
				networkingFirewallRoutes.GET("", networkingHandler.UWFStatusList)
				networkingFirewallRoutes.POST("/status", networkingHandler.UWFStatus)
				networkingFirewallRoutes.POST("/active", networkingHandler.UWFActive)
				networkingFirewallRoutes.POST("/enable", networkingHandler.UWFEnable)
				networkingFirewallRoutes.POST("/disable", networkingHandler.UWFDisable)
				networkingFirewallRoutes.POST("/port/open", networkingHandler.UWFOpenPort)
				networkingFirewallRoutes.POST("/port/close", networkingHandler.UWFClosePort)
			}
		}

		fileRoutes := apiRoutes.Group("/files")
		{
			fileRoutes.GET("/exists", fileHandler.FileExists)            // needs to be a file
			fileRoutes.GET("/walk", fileHandler.WalkFile)                // similar as find in linux command
			fileRoutes.GET("/list", fileHandler.ListFiles)               // list all files and folders
			fileRoutes.POST("/create", fileHandler.CreateFile)           // create file only
			fileRoutes.POST("/copy", fileHandler.CopyFile)               // copy either file or folder
			fileRoutes.POST("/rename", fileHandler.RenameFile)           // rename either file or folder
			fileRoutes.POST("/move", fileHandler.MoveFile)               // move files or folders
			fileRoutes.POST("/upload", fileHandler.UploadFile)           // upload single file
			fileRoutes.POST("/download", fileHandler.DownloadFile)       // download single file
			fileRoutes.GET("/read", fileHandler.ReadFile)                // read single file
			fileRoutes.PUT("/write", fileHandler.WriteFile)              // write single file
			fileRoutes.DELETE("/delete", fileHandler.DeleteFile)         // delete single file
			fileRoutes.DELETE("/delete-all", fileHandler.DeleteAllFiles) // deletes file or folder
			fileRoutes.POST("/write/string", fileHandler.WriteStringFile)
			fileRoutes.POST("/write/json", fileHandler.WriteFileJson)
			fileRoutes.POST("/write/yml", fileHandler.WriteFileYml)
		}

		dirRoutes := apiRoutes.Group("/dirs")
		{
			dirRoutes.GET("/exists", dirHandler.DirExists)  // needs to be a folder
			dirRoutes.POST("/create", dirHandler.CreateDir) // create folder
		}

		zipRoutes := apiRoutes.Group("/zip")
		{
			zipRoutes.POST("/unzip", zipHandler.Unzip)
			zipRoutes.POST("/zip", zipHandler.ZipDir)
		}

		streamLogRoutes := apiRoutes.Group("/logs")
		{
			streamLogRoutes.GET("", streamLogHandler.GetStreamLogs)
			streamLogRoutes.GET("/:uuid", streamLogHandler.GetStreamLog)
			streamLogRoutes.POST("", streamLogHandler.CreateStreamLog)
			streamLogRoutes.POST("/create", streamLogHandler.CreateLogAndReturn)
			streamLogRoutes.DELETE("/:uuid", streamLogHandler.DeleteStreamLog)
			streamLogRoutes.DELETE("", streamLogHandler.DeleteStreamLogs)
		}

		snapshotRoutes := apiRoutes.Group("/snapshots")
		{
			snapshotRoutes.POST("create", snapshotHandler.CreateSnapshot)
			snapshotRoutes.POST("restore", snapshotHandler.RestoreSnapshot)
			snapshotRoutes.GET("status", snapshotHandler.SnapshotStatus)
		}

		restartJobRoutes := apiRoutes.Group("/restart-jobs")
		{
			restartJobRoutes.GET("", restartJobHandler.GetRestartJob)
			restartJobRoutes.PUT("", restartJobHandler.UpdateRestartJob)
			restartJobRoutes.DELETE("/:unit", restartJobHandler.DeleteRestartJob)
		}

		// These APIs are just needed on server level, later we can introduce a flag to restrict exposing
		serverApiRoutes := apiRoutes.Group("")
		{
			storeRoutes := serverApiRoutes.Group("/store")
			{
				appStoreRoutes := storeRoutes.Group("/apps")
				{
					appStoreRoutes.POST("", appStoreHandler.UploadAddOnAppStore)
					appStoreRoutes.GET("/exists", appStoreHandler.CheckAppExistence)
				}

				pluginStoreRoutes := storeRoutes.Group("/plugins")
				{
					pluginStoreRoutes.GET("", pluginStoreHandler.GetPluginsStorePlugins)
					pluginStoreRoutes.POST("", pluginStoreHandler.UploadPluginStorePlugin)
				}
			}

			edgeBiosAppRoutes := serverApiRoutes.Group("/eb/ros")
			{
				edgeBiosAppRoutes.POST("/upload", edgeBiosEdgeHandler.EdgeBiosRubixOsUpload)
				edgeBiosAppRoutes.POST("/install", edgeBiosEdgeHandler.EdgeBiosRubixOsInstall)
				edgeBiosAppRoutes.GET("/version", edgeBiosEdgeHandler.EdgeBiosGetRubixOsVersion)

				edgePluginRoutes := edgeBiosAppRoutes.Group("/plugins")
				{
					edgePluginRoutes.GET("", edgePluginHandler.EdgeListPlugins)
					edgePluginRoutes.POST("/upload", edgePluginHandler.EdgeUploadPlugin)
					edgePluginRoutes.POST("/move-from-download-to-install", edgePluginHandler.EdgeMoveFromDownloadToInstallPlugins)
					edgePluginRoutes.DELETE("/name/:plugin_name", edgePluginHandler.EdgeDeletePlugin)
					edgePluginRoutes.DELETE("/download-plugins", edgePluginHandler.EdgeDeleteDownloadPlugins)
				}
			}

			edgeRoutes := serverApiRoutes.Group("/edge")
			{
				edgeAppRoutes := edgeRoutes.Group("/apps")
				{
					edgeAppRoutes.POST("/upload", edgeAppHandler.EdgeAppUpload)
					edgeAppRoutes.POST("/install", edgeAppHandler.EdgeAppInstall)
					edgeAppRoutes.POST("/uninstall", edgeAppHandler.EdgeAppUninstall)
					edgeAppRoutes.GET("/status", edgeAppHandler.EdgeListAppsStatus)
					edgeAppRoutes.GET("/status/:app_name", edgeAppHandler.EdgeGetAppStatus)
				}

				edgeConfigRoutes := edgeRoutes.Group("/config")
				{
					edgeConfigRoutes.GET("", edgeConfigHandler.EdgeReadConfig)
					edgeConfigRoutes.POST("", edgeConfigHandler.EdgeWriteConfig)
				}

				edgeSnapshotRoutes := edgeRoutes.Group("/snapshots")
				{
					edgeSnapshotRoutes.GET("", edgeSnapshotHandler.GetSnapshots)
					edgeSnapshotRoutes.PATCH("/:file", edgeSnapshotHandler.UpdateSnapshot)
					edgeSnapshotRoutes.DELETE("", edgeSnapshotHandler.DeleteSnapshot)
					edgeSnapshotRoutes.POST("/create", edgeSnapshotHandler.CreateSnapshot)
					edgeSnapshotRoutes.POST("/restore", edgeSnapshotHandler.RestoreSnapshot)
					edgeSnapshotRoutes.POST("/download", edgeSnapshotHandler.DownloadSnapshot)
					edgeSnapshotRoutes.POST("/upload", edgeSnapshotHandler.UploadSnapshot)

					snapshotCreateLogRoutes := edgeSnapshotRoutes.Group("/create-logs")
					{
						snapshotCreateLogRoutes.GET("", snapshotCreatLogHandler.GetSnapshotCreateLogs)
						snapshotCreateLogRoutes.PATCH("/:uuid", snapshotCreatLogHandler.UpdateSnapshotCreateLog)
						snapshotCreateLogRoutes.DELETE("/:uuid", snapshotCreatLogHandler.DeleteSnapshotCreateLog)
					}

					snapshotRestoreLogRoutes := edgeSnapshotRoutes.Group("/restore-logs")
					{
						snapshotRestoreLogRoutes.GET("", snapshotRestoreLogHandler.GetSnapshotRestoreLogs)
						snapshotRestoreLogRoutes.PATCH("/:uuid", snapshotRestoreLogHandler.UpdateSnapshotRestoreLog)
						snapshotRestoreLogRoutes.DELETE("/:uuid", snapshotRestoreLogHandler.DeleteSnapshotRestoreLog)
					}
				}
			}

			locationRoutes := serverApiRoutes.Group("/locations")
			{
				locationRoutes.GET("/schema", locationHandler.GetLocationSchema)
				locationRoutes.GET("", locationHandler.GetLocations)
				locationRoutes.GET("/:uuid", locationHandler.GetLocation)
				locationRoutes.POST("", locationHandler.CreateLocation)
				locationRoutes.PATCH("/:uuid", locationHandler.UpdateLocation)
				locationRoutes.DELETE("/:uuid", locationHandler.DeleteLocation)
				locationRoutes.DELETE("/drop", locationHandler.DropLocations)
			}

			groupRoutes := serverApiRoutes.Group("/groups")
			{
				groupRoutes.GET("/schema", groupHandler.GetGroupSchema)
				groupRoutes.GET("", groupHandler.GetGroups)
				groupRoutes.GET("/:uuid", groupHandler.GetGroup)
				groupRoutes.POST("", groupHandler.CreateGroup)
				groupRoutes.PATCH("/:uuid", groupHandler.UpdateGroup)
				groupRoutes.DELETE("/:uuid", groupHandler.DeleteGroup)
				groupRoutes.DELETE("/drop", groupHandler.DropGroups)
				groupRoutes.GET("/:uuid/update-hosts-status", groupHandler.UpdateHostsStatus)
			}

			hostRoutes := serverApiRoutes.Group("/hosts")
			{
				hostRoutes.GET("/schema", hostHandler.GetHostSchema)
				hostRoutes.GET("", hostHandler.GetHosts)
				hostRoutes.POST("", hostHandler.CreateHost)
				hostRoutes.GET("/:uuid", hostHandler.GetHost)
				hostRoutes.PATCH("/:uuid", hostHandler.UpdateHost)
				hostRoutes.DELETE("/:uuid", hostHandler.DeleteHost)
				hostRoutes.DELETE("/drop", hostHandler.DropHosts)
				hostRoutes.GET("/:uuid/configure-openvpn", hostHandler.ConfigureOpenVPN)

				hostTagRoutes := hostRoutes.Group("/tags")
				{
					hostTagRoutes.PUT("/host_uuid/:host_uuid", hostTagHandler.UpdateHostTags)
				}

				hostCommentRoutes := hostRoutes.Group("/comments")
				{
					hostCommentRoutes.POST("", hostCommentHandler.CreateHostComment)
					hostCommentRoutes.PATCH("/:uuid", hostCommentHandler.UpdateHostComment)
					hostCommentRoutes.DELETE("/:uuid", hostCommentHandler.DeleteHostComment)
				}
			}

			viewRoutes := serverApiRoutes.Group("/views")
			{
				viewRoutes.POST("", viewHandler.CreateView)
				viewRoutes.PATCH("/:uuid", viewHandler.UpdateView)
				viewRoutes.DELETE("/:uuid", viewHandler.DeleteView)
				viewRoutes.POST("/generate-template", viewHandler.GenerateViewTemplate)
				viewRoutes.POST("/assign-template", viewHandler.AssignViewTemplate)

				viewWidgetRoutes := viewRoutes.Group("/widgets")
				{
					viewWidgetRoutes.POST("", viewWidgetHandler.CreateViewWidget)
					viewWidgetRoutes.PATCH("/:uuid", viewWidgetHandler.UpdateViewWidget)
					viewWidgetRoutes.DELETE("/:uuid", viewWidgetHandler.DeleteViewWidget)
				}
			}

			viewTemplateRoutes := serverApiRoutes.Group("/view-templates")
			{
				viewTemplateRoutes.GET("", viewTemplateHandler.GetViewTemplates)
				viewTemplateRoutes.GET("/:uuid", viewTemplateHandler.GetViewTemplate)
				viewTemplateRoutes.POST("", viewTemplateHandler.CreateViewTemplate)
				viewTemplateRoutes.PATCH("/:uuid", viewTemplateHandler.UpdateViewTemplate)
				viewTemplateRoutes.DELETE("/:uuid", viewTemplateHandler.DeleteViewTemplate)

				viewTemplateWidgetRoutes := viewTemplateRoutes.Group("/widgets")
				{
					viewTemplateWidgetRoutes.PATCH("/:uuid", viewTemplateWidgetHandler.UpdateViewTemplateWidget)
					viewTemplateWidgetRoutes.DELETE("/:uuid", viewTemplateWidgetHandler.DeleteViewTemplateWidget)
				}
			}

			alertRoutes := serverApiRoutes.Group("/alerts")
			{
				alertRoutes.GET("/schema", alertHandler.AlertsSchema)
				alertRoutes.POST("", alertHandler.CreateAlert)
				alertRoutes.GET("/:uuid", alertHandler.GetAlert)
				alertRoutes.GET("/host/:uuid", alertHandler.GetAlertsByHost)
				alertRoutes.PATCH("/:uuid/status", alertHandler.UpdateAlertStatus)
				alertRoutes.DELETE("/:uuid", alertHandler.DeleteAlert)
				alertRoutes.DELETE("/drop", alertHandler.DropAlerts)
			}

			teamRoutes := serverApiRoutes.Group("/teams")
			{
				teamRoutes.GET("", teamHandler.GetTeams)
				teamRoutes.GET("/:uuid", teamHandler.GetTeam)
				teamRoutes.POST("", teamHandler.CreateTeam)
				teamRoutes.PATCH("/:uuid", teamHandler.UpdateTeam)
				teamRoutes.DELETE("/:uuid", teamHandler.DeleteTeam)
				teamRoutes.DELETE("/drop", teamHandler.DropTeams)
				teamRoutes.PUT("/:uuid/members", teamHandler.UpdateTeamMembers)
				teamRoutes.PUT("/:uuid/views", teamHandler.UpdateTeamViews)
			}

			edgeCloneRoutes := serverApiRoutes.Group("/clone_edges")
			{
				edgeCloneRoutes.GET("", cloudEdgeCloneHandler.CloneEdge)
			}

			ticketRoutes := serverApiRoutes.Group("/tickets")
			{
				ticketRoutes.POST("", ticketHandler.CreateTicket)
				ticketRoutes.PATCH("/:uuid", ticketHandler.UpdateTicket)
				ticketRoutes.DELETE("/:uuid", ticketHandler.DeleteTicket)
				ticketRoutes.PUT("/:uuid/teams", ticketHandler.UpdateTicketTeams)
			}

			fcmServerRoutes := serverApiRoutes.Group("fcm-server")
			{
				fcmServerRoutes.GET("", fcmServerHandler.GetFcmServer)
				fcmServerRoutes.PUT("", fcmServerHandler.UpsertFcmServer)
			}
		}
	}

	authWithMember := engine.Group("/api", handleAuthWithMember)
	{
		viewSettingRoutes := authWithMember.Group("/view-settings")
		{
			viewSettingRoutes.GET("", viewSettingHandler.GetViewSetting)
			viewSettingRoutes.PUT("", viewSettingHandler.UpsertSetting)
			viewSettingRoutes.DELETE("", viewSettingHandler.DeleteViewSetting)
		}

		viewRoutes := authWithMember.Group("/views")
		{
			viewRoutes.GET("", viewHandler.GetViews)
			viewRoutes.GET("/:uuid", viewHandler.GetView)
		}

		alertRoutes := authWithMember.Group("/alerts")
		{
			alertRoutes.GET("", alertHandler.GetAlerts)
		}

		ticketRoutes := authWithMember.Group("/tickets")
		{
			ticketRoutes.GET("", ticketHandler.GetTickets)
			ticketRoutes.GET("/:uuid", ticketHandler.GetTicket)
			ticketRoutes.PUT("/:uuid/priority", ticketHandler.UpdateTicketPriority)
			ticketRoutes.PUT("/:uuid/status", ticketHandler.UpdateTicketStatus)

			ticketCommentRoutes := ticketRoutes.Group("/comments")
			{
				ticketCommentRoutes.POST("", ticketCommentHandler.CreateTicketComment)
				ticketCommentRoutes.GET("/:uuid", ticketCommentHandler.GetTicketComment)
				ticketCommentRoutes.PATCH("/:uuid", ticketCommentHandler.UpdateTicketComment)
				ticketCommentRoutes.DELETE("/:uuid", ticketCommentHandler.DeleteTicketComment)
			}
		}
	}

	authWithMemberHostLevel := engine.Group("/api/host_points", handleAuthWithMemberHostLevel)
	{
		authWithMemberHostLevel.GET("/:uuid", pointHandler.GetPointByHost)
		authWithMemberHostLevel.PATCH("/write/:uuid", pointHandler.WritePointByHost)
	}

	return engine
}
