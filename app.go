package main

import (
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/services/localmqtt"
	"github.com/NubeIO/flow-framework/services/system"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/internaltoken"
	"github.com/go-co-op/gocron"
	"github.com/robfig/cron/v3"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/database"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/history"
	"github.com/NubeIO/flow-framework/logger"
	"github.com/NubeIO/flow-framework/router"
	"github.com/NubeIO/flow-framework/runner"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/flow-framework/src/jobs"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

var flushMqttPublishBufferInterval = 1 * time.Second

func intHandler(db *database.GormDatabase) {
	dh := new(dbhandler.Handler)
	dh.DB = db
	dbhandler.Init(dh)

	s := new(cachestore.Handler)
	s.Store = cache.New(5*time.Minute, 10*time.Minute)
	cachestore.Init(s)

	j := new(jobs.Jobs)
	j.InitCron()
}

func initHistorySchedulers(db *database.GormDatabase, conf *config.Configuration) {
	h := new(history.History)
	h.DB = db
	if *conf.ProducerHistory.Enable && *conf.ProducerHistory.Cleaner.Enable {
		h.InitProducerHistoryCleaner(
			conf.ProducerHistory.Cleaner.Frequency,
			conf.ProducerHistory.Cleaner.DataPersistingHours)
	}
	if *conf.ProducerHistory.Enable && *conf.ProducerHistory.IntervalHistoryCreator.Enable {
		h.InitIntervalHistoryCreator(conf.ProducerHistory.IntervalHistoryCreator.Frequency)
	}
}

func initFlushBuffers() {
	if boolean.IsTrue(config.Get().MQTT.Enable) {
		go func() {
			for {
				time.Sleep(flushMqttPublishBufferInterval)
				localmqtt.GetLocalMqtt().Client.FlushMqttPublishBuffers()
			}
		}()
	}
}

func setupCron() (*gocron.Scheduler, *systemctl.SystemCtl, *system.System) {
	scheduler := gocron.NewScheduler(time.Local)
	systemCtl := systemctl.New(false, 30)
	system_ := system.New(&system.System{})
	restartJobs := utils.GetRestartJobs()
	for _, restartJob := range restartJobs {
		_, err := cron.ParseStandard(restartJob.Expression)
		if err != nil {
			log.Errorln(err)
		} else {
			_, err = scheduler.Cron(restartJob.Expression).Tag(restartJob.Unit).Do(func() {
				_ = systemCtl.Restart(restartJob.Unit)
			})
			if err != nil {
				log.Errorln(err)
			}
		}
	}

	rebootJob := utils.GetRebootJob()
	if rebootJob != nil {
		_, err := cron.ParseStandard(rebootJob.Expression)
		if err != nil {
			log.Errorln(err)
		}
		_, err = scheduler.Cron(rebootJob.Expression).Tag(rebootJob.Tag).Do(func() {
			_, _ = system_.RebootHost()
		})
	}
	scheduler.StartAsync()
	return scheduler, systemCtl, system_
}

var db *database.GormDatabase

func main() {
	defer db.Close()
	conf := config.CreateApp()

	logger.SetLogger(conf.LogLevel)
	logger.SetGinMode(conf.LogLevel)

	if err := os.MkdirAll(conf.GetAbsPluginDir(), 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(conf.GetAbsUploadedImagesDir(), 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(conf.GetAbsSnapShotDir(), 0755); err != nil {
		panic(err)
	}
	internaltoken.CreateInternalTokenIfDoesNotExist()

	mqttBroker := "tcp://" + conf.MQTT.Address + ":" + strconv.Itoa(conf.MQTT.Port)
	_, err := mqttclient.InternalMQTT(mqttBroker)
	if err != nil {
		log.Errorln(err)
	}

	eventbus.Init()

	connection := path.Join(conf.GetAbsDataDir(), conf.Database.Connection)
	db, err = database.New(conf.Database.Dialect, connection, conf.Database.LogLevel)
	if err != nil {
		panic(err)
	}

	if *conf.MQTT.Enable {
		err = localmqtt.Init(mqttBroker, conf, onConnected)
		if err != nil {
			log.Errorln(err)
		}
	}

	intHandler(db)
	scheduler, systemCtl, system_ := setupCron()
	engine := router.Create(db, conf, scheduler, systemCtl, system_)
	eventbus.RegisterMQTTBus(false)
	initHistorySchedulers(db, conf)
	initFlushBuffers()
	runner.Run(engine, conf)
}

func onConnected() {
	go db.PublishPointsList("")
	go db.RePublishPointsCov()
	go db.PublishPointsListListener()
	go db.PublishSchedulesListener()
	go db.PublishDeviceInfo()
	go db.PublishFetchPointListener()
	go db.PublishPointWriteListener()
	go db.RePublishSelectedPointsCovListener()
}
