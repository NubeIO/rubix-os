package api

import (
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/user"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/bools"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
	"regexp"
)

func getBODYSchedule(ctx *gin.Context) (dto *model.Schedule, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYScheduleData(ctx *gin.Context) (dto *model.ScheduleData, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYFlowNetwork(ctx *gin.Context) (dto *model.FlowNetwork, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getMapBody(ctx *gin.Context) (dto *map[string]interface{}, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYNetwork(ctx *gin.Context) (dto *model.Network, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyHistory(ctx *gin.Context) (dto *model.ProducerHistory, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyBulkHistory(ctx *gin.Context) (dto []*model.ProducerHistory, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYDevice(ctx *gin.Context) (dto *model.Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYProducer(ctx *gin.Context) (dto *model.Producer, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodySyncProducer(ctx *gin.Context) (dto *interfaces.SyncProducer, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYConsumer(ctx *gin.Context) (dto *model.Consumer, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriter(ctx *gin.Context) (dto *model.Writer, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodySyncWriter(ctx *gin.Context) (dto *model.SyncWriter, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodySyncCOV(ctx *gin.Context) (dto *model.SyncCOV, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodySyncWriterAction(ctx *gin.Context) (dto *model.SyncWriterAction, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPlugin(ctx *gin.Context) (dto *model.PluginConf, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYMqttConnection(ctx *gin.Context) (dto *model.MqttConnection, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYIntegration(ctx *gin.Context) (dto *model.Integration, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriterBody(ctx *gin.Context) (dto *model.WriterBody, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriterBulk(ctx *gin.Context) (dto []*model.WriterBulkBody, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriterClone(ctx *gin.Context) (dto *model.WriterClone, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYStream(ctx *gin.Context) (dto *model.Stream, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodySyncStream(ctx *gin.Context) (dto *model.SyncStream, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYCommandGroup(ctx *gin.Context) (dto *model.CommandGroup, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYJobs(ctx *gin.Context) (dto *model.Job, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPointMapping(ctx *gin.Context) (dto *model.PointMapping, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYBulkPoints(ctx *gin.Context) (dto []*model.Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPoint(ctx *gin.Context) (dto *model.Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPointWriter(ctx *gin.Context) (dto *model.PointWriter, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyTag(ctx *gin.Context) (dto *model.Tag, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyUser(ctx *gin.Context) (dto *user.User, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyTokenCreate(ctx *gin.Context) (dto *interfaces.TokenCreate, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyTokenBlock(ctx *gin.Context) (dto *interfaces.TokenBlock, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyBulkNetworkMetaTags(ctx *gin.Context) (dto []*model.NetworkMetaTag, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyBulkDeviceMetaTag(ctx *gin.Context) (dto []*model.DeviceMetaTag, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyBulkPointMetaTag(ctx *gin.Context) (dto []*model.PointMetaTag, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getP2PBody(ctx *gin.Context) (dto *model.P2PBody, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyAutoMapping(ctx *gin.Context) (dto *interfaces.AutoMapping, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyLocation(ctx *gin.Context) (dto *model.Location, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyView(ctx *gin.Context) (dto *model.View, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyGroup(ctx *gin.Context) (dto *model.Group, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyHost(ctx *gin.Context) (dto *model.Host, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyHostComment(ctx *gin.Context) (dto *model.HostComment, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyHostTags(ctx *gin.Context) (dto []*model.HostTag, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodySnapshotLog(ctx *gin.Context) (dto *model.SnapshotLog, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyCreateSnapshot(ctx *gin.Context) (dto *interfaces.CreateSnapshot, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodySnapshotCreateLog(ctx *gin.Context) (dto *model.SnapshotCreateLog, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyRestoreSnapshot(ctx *gin.Context) (dto *interfaces.RestoreSnapshot, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodySnapshotRestoreLog(ctx *gin.Context) (dto *model.SnapshotRestoreLog, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyRebootJob(ctx *gin.Context) (dto *interfaces.RebootJob, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyRestartJob(ctx *gin.Context) (dto *interfaces.RestartJob, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyStreamLog(ctx *gin.Context) (dto *interfaces.StreamLog, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyAlert(ctx *gin.Context) (dto *model.Alert, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyMember(ctx *gin.Context) (dto *model.Member, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyMemberGroups(ctx *gin.Context) (dto []*string, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyChangePassword(ctx *gin.Context) (dto *interfaces.ChangePassword, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyMemberDevice(ctx *gin.Context) (dto *model.MemberDevice, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveGlobalUUID(ctx *gin.Context) string {
	return ctx.Param("global_uuid")
}

func resolveID(ctx *gin.Context) string {
	return ctx.Param("uuid")
}

func resolveWriterUUID(ctx *gin.Context) string {
	return ctx.Param("writer_uuid")
}

func resolveSourceUUID(ctx *gin.Context) string {
	return ctx.Param("source_uuid")
}

func resolveWriterThingClass(ctx *gin.Context) string {
	return ctx.Param("writer_thing_class")
}

func resolveProducerUUID(ctx *gin.Context) string {
	return ctx.Param("producer_uuid")
}

func getTagParam(ctx *gin.Context) string {
	return ctx.Param("tag")
}

func resolvePath(ctx *gin.Context) string {
	return ctx.Param("path")
}

func resolvePluginUUID(ctx *gin.Context) string {
	return ctx.Param("plugin_uuid")
}

func resolveName(ctx *gin.Context) string {
	return ctx.Param("name")
}

func resolvePluginName(ctx *gin.Context) string {
	return ctx.Param("plugin_name")
}

func resolveNetworkName(ctx *gin.Context) string {
	return ctx.Param("network_name")
}

func resolveDeviceId(ctx *gin.Context) string {
	return ctx.Param("device_id")
}

func resolveDeviceName(ctx *gin.Context) string {
	return ctx.Param("device_name")
}

func resolvePointName(ctx *gin.Context) string {
	return ctx.Param("point_name")
}

func resolveFlowNetworkCloneName(ctx *gin.Context) string {
	return ctx.Param("flow_network_clone_name")
}

func resolveStreamCloneName(ctx *gin.Context) string {
	return ctx.Param("stream_clone_name")
}

func resolveConsumerName(ctx *gin.Context) string {
	return ctx.Param("consumer_name")
}

func resolveWriterThingName(ctx *gin.Context) string {
	return ctx.Param("writer_thing_name")
}

func toBool(value string) (bool, error) {
	if value == "" {
		return false, nil
	} else {
		c, err := bools.Boolean(value)
		return c, err
	}
}

func validUsername(username string) bool {
	re, _ := regexp.Compile("^([A-Za-z0-9_-])+$")
	return re.FindString(username) != ""
}

func resolveHeaderHostID(ctx *gin.Context) string {
	return ctx.GetHeader("host-uuid")
}

func resolveHeaderHostName(ctx *gin.Context) string {
	return ctx.GetHeader("host-name")
}

func resolveUsername(ctx *gin.Context) string {
	return ctx.Param("username")
}

func resolveEmail(ctx *gin.Context) string {
	return ctx.Param("email")
}

func matchUUID(uuid string) bool {
	if len(uuid) == 16 {
		if uuid[0:4] == "hos_" {
			return true
		}
	}
	return false
}

func matchHostUUID(ctx *gin.Context) string {
	hostID := resolveHeaderHostID(ctx)
	if len(hostID) == 16 {
		if matchUUID(hostID) {
			return hostID
		}
	}
	return ""
}

func matchHostName(ctx *gin.Context) string {
	name := resolveHeaderHostName(ctx)
	return name
}

func matchHostUUIDName(ctx *gin.Context) (string, string) {
	return matchHostUUID(ctx), matchHostName(ctx)
}
