package urls

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"strings"
)

const FlowNetworkUrl string = "/api/flow_networks"
const FlowNetworkCloneUrl string = "/api/flow_network_clones"
const StreamUrl string = "/api/streams"
const StreamCloneUrl string = "/api/stream_clones"
const ProducerUrl string = "/api/producers"
const ConsumerUrl string = "/api/consumers"
const WriterCloneUrl string = "/api/producers/writer_clones"
const WriterUrl string = "/api/consumers/writers"
const NetworkUrl string = "/api/networks"
const DeviceUrl string = "/api/devices"
const PointUrl string = "/api/points"

const FlowNetworkStreamsSyncUrl string = "/api/flow_networks/:uuid/sync/streams"
const StreamProducersSyncUrl string = "/api/streams/:uuid/sync/producers"
const ProducerWriterClonesSyncUrl string = "/api/producers/:uuid/sync/writer_clones"

const FlowNetworkCloneStreamClonesSyncUrl string = "/api/flow_network_clones/:uuid/sync/stream_clones"
const StreamCloneConsumersSyncUrl string = "/api/stream_clones/:uuid/sync/consumers"
const ConsumerWritersSyncUrl string = "/api/consumers/:uuid/sync/writers"

const NetworkDevicesSyncUrl string = "/api/networks/:uuid/sync/devices"
const DevicePointsSyncUrl string = "api/devices/:uuid/sync/points"

func SingularUrl(url, uuid string) string {
	return fmt.Sprintf("%s/%s", url, uuid)
}

func SingularUrlByArg(url, name, value string) string {
	return fmt.Sprintf("%s/one/args?%s=%s", url, name, value)
}

func PluralUrlByArg(url, name, value string) string {
	return fmt.Sprintf("%s?%s=%s", url, name, value)
}

func GetUrl(url, uuid string) string {
	return strings.Replace(url, ":uuid", uuid, -1)
}

func GenerateFNUrlParams(args api.Args) string {
	params := "?"
	var aType = api.ArgsType
	params += fmt.Sprintf("%s=%v", aType.WithProducers, args.WithProducers)
	params += fmt.Sprintf("&%s=%v", aType.WithWriterClones, args.WithWriterClones)
	return params
}

func GenerateFNCUrlParams(args api.Args) string {
	params := "?"
	var aType = api.ArgsType
	params += fmt.Sprintf("%s=%v", aType.WithConsumers, args.WithConsumers)
	params += fmt.Sprintf("&%s=%v", aType.WithWriters, args.WithWriters)
	return params
}

func GenerateNetworkUrlParams(args api.Args) string {
	params := "?"
	var aType = api.ArgsType
	params += fmt.Sprintf("%s=%v", aType.WithDevices, args.WithDevices)
	params += fmt.Sprintf("&%s=%v", aType.WithPoints, args.WithPoints)
	params += fmt.Sprintf("&%s=%v", aType.WithTags, args.WithTags)
	return params
}
