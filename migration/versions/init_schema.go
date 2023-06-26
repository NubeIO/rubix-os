package versions

import "github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"

func GetInitInterfaces() []interface{} {
	return []interface{}{
		&model.Alert{},
		&model.PluginConf{},
		&model.Network{},
		&model.Device{},
		&model.Point{},
		&model.Priority{},
		&model.Job{},
		&model.Integration{},
		&model.MqttConnection{},
		&model.Schedule{},
		&model.Tag{},
		&model.History{},
		&model.HistoryLog{},
		&model.HistoryPostgresLog{},
		&model.NetworkMetaTag{},
		&model.DeviceMetaTag{},
		&model.PointMetaTag{},
		&model.Location{},
		&model.Group{},
		&model.Host{},
		&model.HostTag{},
		&model.HostComment{},
		&model.SnapshotLog{},
		&model.SnapshotCreateLog{},
		&model.SnapshotRestoreLog{},
		&model.Member{},
		&model.MemberDevice{},
		&model.ViewSetting{},
		&model.View{},
		&model.ViewWidget{},
		&model.ViewTemplate{},
		&model.ViewTemplateWidget{},
		&model.ViewTemplateWidgetPointer{},
		&model.Team{},
		&model.TeamView{},
		&model.PointHistory{},
	}
}
