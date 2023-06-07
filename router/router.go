package router

import (
	"github.com/NubeDev/location"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/auth"
	"github.com/NubeIO/rubix-os/config"
	"github.com/NubeIO/rubix-os/database"
	"github.com/NubeIO/rubix-os/eventbus"
	"github.com/NubeIO/rubix-os/global"
	"github.com/NubeIO/rubix-os/installer"
	"github.com/NubeIO/rubix-os/logger"
	"github.com/NubeIO/rubix-os/module"
	"github.com/NubeIO/rubix-os/nerrors"
	"github.com/NubeIO/rubix-os/plugin"
	"github.com/NubeIO/rubix-os/services/appstore"
	"github.com/NubeIO/rubix-os/services/system"
	"github.com/NubeIO/rubix-registry-go/rubixregistry"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func Create(db *database.GormDatabase, conf *config.Configuration, scheduler *gocron.Scheduler,
	systemCtl *systemctl.SystemCtl, system_ *system.System) *gin.Engine {
	engine := gin.New()
	engine.Use(logger.GinMiddlewareLogger(), gin.Recovery(), nerrors.Handler(), location.Default())
	engine.NoRoute(nerrors.NotFoundHandler())
	eventBus := eventbus.NewService(eventbus.GetBus())
	global.Installer = installer.New(&installer.Installer{})
	proxyHandler := api.Proxy{DB: db}
	healthHandler := api.HealthAPI{DB: db}

	authHandler := api.AuthAPI{}
	handleAuth := func(c *gin.Context) { c.Next() }
	if *conf.Auth {
		handleAuth = authHandler.HandleAuth()
	}
	apiRoutesTemp := engine.Group("/api", handleAuth) // TODO: remove this one and use the same one
	// http://localhost:1660/api/plugins/api/system/schema/json/device
	pluginManager, err := plugin.NewManager(db, conf.GetAbsPluginsDir(), apiRoutesTemp.Group("/plugins/api"))
	if err != nil {
		log.Error(err)
		panic(err)
	}

	modules, err := module.ReLoadModulesWithDir(config.Get().GetAbsModulesDir())
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
	historyHandler := api.HistoriesAPI{
		DB: db,
	}
	jobHandler := api.JobAPI{
		DB: db,
	}
	streamHandler := api.StreamAPI{
		DB: db,
	}
	remoteHandler := api.RemoteAPI{
		DB: db,
	}
	streamCloneHandler := api.StreamCloneAPI{
		DB: db,
	}
	producerHandler := api.ProducerAPI{
		DB: db,
	}
	consumerHandler := api.ConsumersAPI{
		DB: db,
	}
	writerCloneHandler := api.WriterCloneAPI{
		DB: db,
	}
	rubixCommandGroup := api.CommandGroupAPI{
		DB: db,
	}
	flowNetwork := api.FlowNetworksAPI{
		DB: db,
	}
	mapping := api.MappingAPI{
		DB: db,
	}
	flowNetworkCloneHandler := api.FlowNetworkClonesAPI{
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
		RubixRegistry: rubixregistry.New(),
	}
	writerHandler := api.WriterAPI{
		DB: db,
	}
	syncFlowNetworkHandler := api.SyncFlowNetworkAPI{
		DB: db,
	}
	syncStreamHandler := api.SyncStreamAPI{
		DB: db,
	}
	syncWriterHandler := api.SyncWriterAPI{
		DB: db,
	}
	autoMappingHandler := api.AutoMappingAPI{
		DB: db,
	}
	autoMappingScheduleHandler := api.AutoMappingScheduleAPI{
		DB: db,
	}
	syncProducerHandler := api.SyncProducerAPI{
		DB: db,
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
		RubixRegistry: rubixregistry.New(),
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

	engine.Use(cors.New(auth.CorsConfig(conf)))
	apiProxyWiresRoutes := engine.Group("/wires", handleAuth)
	apiProxyWiresRoutes.Any("/*proxyPath", wiresProxyHandler.WiresProxy) // EDGE-WIRES PROXY
	apiProxyChirpRoutes := engine.Group("/chirp", handleAuth)
	apiProxyChirpRoutes.Any("/*proxyPath", chirpProxyHandler.ChirpProxy) // CHIRP-STACK PROXY
	apiProxyHostRoutes := engine.Group("/proxy", handleAuth)
	apiProxyHostRoutes.Any("/*proxyPath", hostProxyHandler.HostProxy)

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
		fnProxy := apiRoutes.Group("/fn")
		{
			fnProxy.GET("/*any", proxyHandler.GetProxy(true))
			fnProxy.POST("/*any", proxyHandler.PostProxy(true))
			fnProxy.PUT("/*any", proxyHandler.PutProxy(true))
			fnProxy.PATCH("/*any", proxyHandler.PatchProxy(true))
			fnProxy.DELETE("/*any", proxyHandler.DeleteProxy(true))
		}

		fncProxy := apiRoutes.Group("/fnc")
		{
			fncProxy.GET("/*any", proxyHandler.GetProxy(false))
			fncProxy.POST("/*any", proxyHandler.PostProxy(false))
			fncProxy.PUT("/*any", proxyHandler.PutProxy(false))
			fncProxy.PATCH("/*any", proxyHandler.PatchProxy(false))
			fncProxy.DELETE("/*any", proxyHandler.DeleteProxy(false))
		}

		requireClientsGroupRoutes := apiRoutes.Group("")
		{
			plugins := requireClientsGroupRoutes.Group("/plugins")
			{
				plugins.GET("", pluginHandler.GetPlugins)
				plugins.GET("/:uuid", pluginHandler.GetPlugin)
				plugins.GET("/config/:uuid", pluginHandler.GetConfig)
				plugins.POST("/config/:uuid", pluginHandler.UpdateConfig)
				plugins.GET("/display/:uuid", pluginHandler.GetDisplay)
				plugins.POST("/enable/:uuid", pluginHandler.EnablePluginByUUID)
				plugins.POST("/restart/:uuid", pluginHandler.RestartPlugin)
				plugins.POST("/restart/name/:name", pluginHandler.RestartPluginByName)
				plugins.GET("/path/:path", pluginHandler.GetPluginByPath)
			}
		}

		historyProducerRoutes := apiRoutes.Group("/histories/producers")
		{
			historyProducerRoutes.GET("", historyHandler.GetProducerHistories)
			historyProducerRoutes.GET("/:producer_uuid", historyHandler.GetProducerHistoriesByProducerUUID)
			historyProducerRoutes.GET("/name/:name/one", historyHandler.GetLatestProducerHistoryByProducerName)
			historyProducerRoutes.GET("/name/:name", historyHandler.GetProducerHistoriesByProducerName)
			historyProducerRoutes.GET("/:producer_uuid/one", historyHandler.GetLatestProducerHistoryByProducerUUID)
			historyProducerRoutes.POST("/point_uuids", historyHandler.GetProducerHistoriesByPointUUIDs)
			historyProducerRoutes.GET("/points", historyHandler.GetProducerHistoriesPoints)
			historyProducerRoutes.GET("/points_for_sync", historyHandler.GetProducerHistoriesPointsForSync)
			historyProducerRoutes.DELETE("/:producer_uuid", historyHandler.DeleteProducerHistoriesByProducerUUID)
		}

		flowNetworkRoutes := apiRoutes.Group("/flow_networks")
		{
			flowNetworkRoutes.GET("", flowNetwork.GetFlowNetworks)
			flowNetworkRoutes.POST("", flowNetwork.CreateFlowNetwork)
			flowNetworkRoutes.GET("/:uuid", flowNetwork.GetFlowNetwork)
			flowNetworkRoutes.PATCH("/:uuid", flowNetwork.UpdateFlowNetwork)
			flowNetworkRoutes.DELETE("/:uuid", flowNetwork.DeleteFlowNetwork)
			flowNetworkRoutes.GET("/one/args", flowNetwork.GetOneFlowNetworkByArgs)
			flowNetworkRoutes.GET("/refresh_connections", flowNetwork.RefreshFlowNetworksConnections)
			flowNetworkRoutes.GET("/sync", flowNetwork.SyncFlowNetworks)
			flowNetworkRoutes.GET("/:uuid/sync/streams", flowNetwork.SyncFlowNetworkStreams)
		}

		flowNetworkCloneRoutes := apiRoutes.Group("/flow_network_clones")
		{
			flowNetworkCloneRoutes.GET("", flowNetworkCloneHandler.GetFlowNetworkClones)
			flowNetworkCloneRoutes.GET("/:uuid", flowNetworkCloneHandler.GetFlowNetworkClone)
			flowNetworkCloneRoutes.DELETE("/:uuid", flowNetworkCloneHandler.DeleteFlowNetworkClone)
			flowNetworkCloneRoutes.GET("/one/args", flowNetworkCloneHandler.GetOneFlowNetworkCloneByArgs)
			flowNetworkCloneRoutes.DELETE("/one/args", flowNetworkCloneHandler.DeleteOneFlowNetworkCloneByArgs)
			flowNetworkCloneRoutes.GET("/refresh_connections", flowNetworkCloneHandler.RefreshFlowNetworkClonesConnections)
			flowNetworkCloneRoutes.GET("/sync", flowNetworkCloneHandler.SyncFlowNetworkClones)
			flowNetworkCloneRoutes.GET("/:uuid/sync/stream_clones", flowNetworkCloneHandler.SyncFlowNetworkCloneStreamClones)
		}

		streamRoutes := apiRoutes.Group("/streams")
		{
			streamRoutes.GET("", streamHandler.GetStreams)
			streamRoutes.POST("", streamHandler.CreateStream)
			streamRoutes.GET("/:uuid", streamHandler.GetStream)
			streamRoutes.PATCH("/:uuid", streamHandler.UpdateStream)
			streamRoutes.DELETE("/:uuid", streamHandler.DeleteStream)
			streamRoutes.GET("/:uuid/sync/producers", streamHandler.SyncStreamProducers)
		}

		mappingRoutes := apiRoutes.Group("/mapping")
		{
			mappingRoutes.POST("/points", mapping.CreatePointMapping)
		}

		remoteRoutes := apiRoutes.Group("/remote")
		{
			remoteRoutes.GET("/flow_network_clones", remoteHandler.RemoteGetFlowNetworkClones)
			remoteRoutes.GET("/flow_network_clones/:uuid", remoteHandler.RemoteGetFlowNetworkClone)
			remoteRoutes.DELETE("/flow_network_clones", remoteHandler.RemoteDeleteFlowNetworkClone)

			remoteRoutes.GET("/networks", remoteHandler.RemoteGetNetworks)
			remoteRoutes.GET("/networks/:uuid", remoteHandler.RemoteGetNetwork)
			remoteRoutes.POST("/networks", remoteHandler.RemoteCreateNetwork)
			remoteRoutes.PATCH("/networks/:uuid", remoteHandler.RemoteEditNetwork)
			remoteRoutes.DELETE("/networks", remoteHandler.RemoteDeleteNetwork)

			remoteRoutes.GET("/devices", remoteHandler.RemoteGetDevices)
			remoteRoutes.GET("/devices/:uuid", remoteHandler.RemoteGetDevice)
			remoteRoutes.POST("/devices", remoteHandler.RemoteCreateDevice)
			remoteRoutes.PATCH("/devices/:uuid", remoteHandler.RemoteEditDevice)
			remoteRoutes.DELETE("/devices", remoteHandler.RemoteDeleteDevice)

			remoteRoutes.GET("/points", remoteHandler.RemoteGetPoints)
			remoteRoutes.GET("/points/:uuid", remoteHandler.RemoteGetPoint)
			remoteRoutes.POST("/points", remoteHandler.RemoteCreatePoint)
			remoteRoutes.PATCH("/points/:uuid", remoteHandler.RemoteEditPoint)
			remoteRoutes.DELETE("/points", remoteHandler.RemoteDeletePoint)

			remoteRoutes.GET("/streams", remoteHandler.RemoteGetStreams)
			remoteRoutes.GET("/streams/:uuid", remoteHandler.RemoteGetStream)
			remoteRoutes.POST("/streams", remoteHandler.RemoteCreateStream)
			remoteRoutes.PATCH("/streams/:uuid", remoteHandler.RemoteEditStream)
			remoteRoutes.DELETE("/streams", remoteHandler.RemoteDeleteStream)

			remoteRoutes.GET("/stream_clones", remoteHandler.RemoteGetStreamClones)
			remoteRoutes.DELETE("/stream_clones", remoteHandler.RemoteDeleteStreamClone)

			remoteRoutes.GET("/producers", remoteHandler.RemoteGetProducers)
			remoteRoutes.GET("/producers/:uuid", remoteHandler.RemoteGetProducer)
			remoteRoutes.POST("/producers", remoteHandler.RemoteCreateProducer)
			remoteRoutes.PATCH("/producers/:uuid", remoteHandler.RemoteEditProducer)
			remoteRoutes.DELETE("/producers", remoteHandler.RemoteDeleteProducer)

			remoteRoutes.GET("/consumers", remoteHandler.RemoteGetConsumers)
			remoteRoutes.GET("/consumers/:uuid", remoteHandler.RemoteGetConsumer)
			remoteRoutes.POST("/consumers", remoteHandler.RemoteCreateConsumer)
			remoteRoutes.PATCH("/consumers/:uuid", remoteHandler.RemoteEditConsumer)
			remoteRoutes.DELETE("/consumers", remoteHandler.RemoteDeleteConsumer)

			remoteRoutes.GET("/writers", remoteHandler.RemoteGetWriters)
			remoteRoutes.GET("/writers/:uuid", remoteHandler.RemoteGetWriter)
			remoteRoutes.POST("/writers", remoteHandler.RemoteCreateWriter)
			remoteRoutes.PATCH("/writers/:uuid", remoteHandler.RemoteEditWriter)
			remoteRoutes.DELETE("/writers", remoteHandler.RemoteDeleteWriter)
		}

		streamCloneRoutes := apiRoutes.Group("/stream_clones")
		{
			streamCloneRoutes.GET("", streamCloneHandler.GetStreamClones)
			streamCloneRoutes.GET("/:uuid", streamCloneHandler.GetStreamClone)
			streamCloneRoutes.DELETE("/:uuid", streamCloneHandler.DeleteStreamClone)
			streamCloneRoutes.DELETE("/one/args", streamCloneHandler.DeleteOneStreamCloneByArgs)
			streamCloneRoutes.GET("/:uuid/sync/consumers", streamCloneHandler.SyncStreamCloneConsumers)
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
			networkRoutes.GET("/sync", networkHandler.SyncNetworks)
			networkRoutes.GET("/:uuid/sync/devices", networkHandler.SyncNetworkDevices)
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
			deviceRoutes.GET("/:uuid/sync/points", deviceHandler.SyncDevicePoints)
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

		commandRoutes := apiRoutes.Group("/commands")
		{
			commandRoutes.GET("", rubixCommandGroup.GetCommandGroups)
			commandRoutes.POST("", rubixCommandGroup.CreateCommandGroup)
			commandRoutes.GET("/:uuid", rubixCommandGroup.GetCommandGroup)
			commandRoutes.PATCH("/:uuid", rubixCommandGroup.UpdateCommandGroup)
			commandRoutes.DELETE("/:uuid", rubixCommandGroup.DeleteCommandGroup)
		}

		producerRoutes := apiRoutes.Group("/producers")
		{
			producerRoutes.GET("", producerHandler.GetProducers)
			producerRoutes.POST("", producerHandler.CreateProducer)
			producerRoutes.GET("/:uuid", producerHandler.GetProducer)
			producerRoutes.PATCH("/:uuid", producerHandler.UpdateProducer)
			producerRoutes.DELETE("/:uuid", producerHandler.DeleteProducer)
			producerRoutes.GET("/one/args", producerHandler.GetOneProducerByArgs)
			producerRoutes.GET("/:uuid/sync/writer_clones", producerHandler.SyncProducerWriterClones)

			producerWriterCloneRoutes := producerRoutes.Group("/writer_clones")
			{
				producerWriterCloneRoutes.GET("", writerCloneHandler.GetWriterClones)
				producerWriterCloneRoutes.POST("", writerCloneHandler.CreateWriterClone)
				producerWriterCloneRoutes.GET("/:uuid", writerCloneHandler.GetWriterClone)
				producerWriterCloneRoutes.DELETE("/:uuid", writerCloneHandler.DeleteWriterClone)
				producerWriterCloneRoutes.DELETE("/one/args", writerCloneHandler.DeleteOneWriterCloneByArgs)
			}
		}

		consumerRoutes := apiRoutes.Group("/consumers")
		{
			consumerRoutes.GET("", consumerHandler.GetConsumers)
			consumerRoutes.POST("", consumerHandler.CreateConsumer)
			consumerRoutes.GET("/:uuid", consumerHandler.GetConsumer)
			consumerRoutes.PATCH("/:uuid", consumerHandler.UpdateConsumer)
			consumerRoutes.DELETE("/:uuid", consumerHandler.DeleteConsumer)
			consumerRoutes.DELETE("", consumerHandler.DeleteConsumers)
			consumerRoutes.GET("/:uuid/sync/writers", consumerHandler.SyncConsumerWriters)

			consumerWriterRoutes := consumerRoutes.Group("/writers")
			{
				consumerWriterRoutes.GET("", writerHandler.GetWriters)
				consumerWriterRoutes.POST("", writerHandler.CreateWriter)
				consumerWriterRoutes.GET("/:uuid", writerHandler.GetWriter)
				consumerWriterRoutes.GET("/name/:flow_network_clone_name/:stream_clone_name/:consumer_name/:writer_thing_name", writerHandler.GetWriterByName)
				consumerWriterRoutes.PATCH("/:uuid", writerHandler.UpdateWriter)
				consumerWriterRoutes.DELETE("/:uuid", writerHandler.DeleteWriter)
			}
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
			schRoutes.GET("/sync", schHandler.SyncSchedules)
			schRoutes.GET("/sync/:uuid", schHandler.SyncSchedule)
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
			deviceInfoRoutes.GET("/device_info", deviceInfoHandler.GetDeviceInfo)

			deviceInfoRoutes.GET("/device", deviceInfoHandler.GetDeviceInfo)
			deviceInfoRoutes.PATCH("/device", deviceInfoHandler.UpdateDeviceInfo)
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

		apiRoutes.POST("/writers/action/:uuid", writerHandler.WriterAction)
		apiRoutes.POST("/writers/action/bulk", writerHandler.WriterBulkAction)

		syncRoutes := apiRoutes.Group("/sync")
		{
			syncRoutes.POST("/flow_network", syncFlowNetworkHandler.SyncFlowNetwork)
			syncRoutes.POST("/stream", syncStreamHandler.SyncStream)
			syncRoutes.POST("/writer", syncWriterHandler.SyncWriter)
			syncRoutes.POST("/cov/:writer_uuid", syncWriterHandler.SyncCOV) // clone ---> source side
			syncRoutes.POST("/writer/write/:source_uuid", syncWriterHandler.SyncWriterWriteAction)
			syncRoutes.GET("/writer/read/:source_uuid", syncWriterHandler.SyncWriterReadAction)
			syncRoutes.POST("/producer", syncProducerHandler.SyncProducer)
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

		autoMappingRoutes := apiRoutes.Group("/auto_mappings")
		{
			autoMappingRoutes.POST("", autoMappingHandler.CreateAutoMapping)
		}

		scheduleAutoMappingRoutes := apiRoutes.Group("/auto_mapping_schedules") // RE
		{
			scheduleAutoMappingRoutes.POST("", autoMappingScheduleHandler.CreateAutoMappingSchedule)
		}

		memberRoutes := apiRoutes.Group("/members")
		{
			memberRoutes.GET("", memberHandler.GetMembers)
			memberRoutes.GET("/:uuid", memberHandler.GetMemberByUUID)
			memberRoutes.DELETE("/:uuid", memberHandler.DeleteMemberByUUID)
			memberRoutes.PATCH("/:uuid", memberHandler.UpdateMemberByUUID)
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
				networkingInterfaceRoutes.POST("/reset", networkingHandler.InterfaceUpDown) //
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

			viewSettingRoutes := serverApiRoutes.Group("/view-settings")
			{
				viewSettingRoutes.GET("", viewSettingHandler.GetViewSetting)
				viewSettingRoutes.POST("", viewSettingHandler.CreateViewSetting)
				viewSettingRoutes.DELETE("", viewSettingHandler.DeleteViewSetting)
			}

			viewRoutes := serverApiRoutes.Group("/views")
			{
				viewRoutes.GET("", viewHandler.GetViews)
				viewRoutes.GET("/:uuid", viewHandler.GetView)
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
				alertRoutes.GET("", alertHandler.GetAlerts)
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
		}
	}
	return engine
}
