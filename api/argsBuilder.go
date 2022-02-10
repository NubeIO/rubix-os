package api

import (
	"github.com/gin-gonic/gin"
)

// TODO: REMOVE
func streamFieldsArgs(ctx *gin.Context) (flowUUID string, streamUUID string, producerUUID string, consumerUUID string, writerUUID string) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.FlowUUID = ctx.DefaultQuery(aType.FlowUUID, aDefault.FlowUUID)
	//args.StreamUUID = ctx.DefaultQuery(aType.StreamUUID, aDefault.StreamUUID)
	args.ProducerUUID = ctx.DefaultQuery(aType.ProducerUUID, aDefault.ProducerUUID)
	args.ConsumerUUID = ctx.DefaultQuery(aType.ConsumerUUID, aDefault.ConsumerUUID)
	args.WriterUUID = ctx.DefaultQuery(aType.WriterUUID, aDefault.WriterUUID)
	return args.FlowUUID, "todo remove", args.ProducerUUID, args.ConsumerUUID, args.WriterUUID
}

//withFieldsArgs
func withFieldsArgs(ctx *gin.Context) (field string, value string) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.Field = ctx.DefaultQuery(aType.Field, aDefault.Field)
	args.Value = ctx.DefaultQuery(aType.Value, aDefault.Value)
	return args.Field, args.Value
}

//withFieldsArgs
func networkDevicePointNames(ctx *gin.Context) (networkName, deviceName, pointName string) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.NetworkName = ctx.DefaultQuery(aType.NetworkName, aDefault.NetworkName)
	args.DeviceName = ctx.DefaultQuery(aType.DeviceName, aDefault.DeviceName)
	args.PointName = ctx.DefaultQuery(aType.PointName, aDefault.PointName)
	return args.NetworkName, args.DeviceName, args.PointName
}

//parentArgs
func parentArgs(ctx *gin.Context) (AddToParent string) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.AddToParent = ctx.DefaultQuery(aType.AddToParent, aDefault.AddToParent)
	return args.AddToParent
}

//withConsumerArgs
func withConsumerArgs(ctx *gin.Context) (askResponse bool, askRefresh bool, write bool, updateProducer bool) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.AskRefresh = ctx.DefaultQuery(aType.AskRefresh, aDefault.AskRefresh)
	args.AskResponse = ctx.DefaultQuery(aType.AskResponse, aDefault.AskResponse)
	args.Write = ctx.DefaultQuery(aType.Write, aDefault.Write)
	args.UpdateProducer = ctx.DefaultQuery(aType.UpdateProducer, aDefault.UpdateProducer)
	askRefresh, _ = toBool(args.AskRefresh)
	askResponse, _ = toBool(args.AskResponse)
	write, _ = toBool(args.Write)
	updateProducer, _ = toBool(args.UpdateProducer)
	return askRefresh, askResponse, write, updateProducer
}

func buildFlowNetworkArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithStreams, _ = toBool(ctx.DefaultQuery(aType.WithStreams, aDefault.WithStreams))
	args.WithProducers, _ = toBool(ctx.DefaultQuery(aType.WithProducers, aDefault.WithProducers))
	args.WithCommandGroups, _ = toBool(ctx.DefaultQuery(aType.WithCommandGroups, aDefault.WithCommandGroups))
	args.WithWriterClones, _ = toBool(ctx.DefaultQuery(aType.WithWriterClones, aDefault.WithWriterClones))
	if value, ok := ctx.GetQuery(aType.GlobalUUID); ok {
		args.GlobalUUID = &value
	}
	if value, ok := ctx.GetQuery(aType.ClientId); ok {
		args.ClientId = &value
	}
	if value, ok := ctx.GetQuery(aType.SiteId); ok {
		args.SiteId = &value
	}
	if value, ok := ctx.GetQuery(aType.DeviceId); ok {
		args.DeviceId = &value
	}
	return args
}

func buildFlowNetworkCloneArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithStreamClones, _ = toBool(ctx.DefaultQuery(aType.WithStreamClones, aDefault.WithStreamClones))
	args.WithConsumers, _ = toBool(ctx.DefaultQuery(aType.WithConsumers, aDefault.WithConsumers))
	args.WithWriters, _ = toBool(ctx.DefaultQuery(aType.WithWriters, aDefault.WithWriters))
	if value, ok := ctx.GetQuery(aType.GlobalUUID); ok {
		args.GlobalUUID = &value
	}
	if value, ok := ctx.GetQuery(aType.ClientId); ok {
		args.ClientId = &value
	}
	if value, ok := ctx.GetQuery(aType.SiteId); ok {
		args.SiteId = &value
	}
	if value, ok := ctx.GetQuery(aType.DeviceId); ok {
		args.DeviceId = &value
	}
	return args
}

func buildStreamArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithFlowNetworks, _ = toBool(ctx.DefaultQuery(aType.WithFlowNetworks, aDefault.WithFlowNetworks))
	args.WithProducers, _ = toBool(ctx.DefaultQuery(aType.WithProducers, aDefault.WithProducers))
	args.WithCommandGroups, _ = toBool(ctx.DefaultQuery(aType.WithCommandGroups, aDefault.WithCommandGroups))
	args.WithWriterClones, _ = toBool(ctx.DefaultQuery(aType.WithWriterClones, aDefault.WithWriterClones))
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	return args
}

func buildStreamCloneArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithConsumers, _ = toBool(ctx.DefaultQuery(aType.WithConsumers, aDefault.WithConsumers))
	args.WithWriters, _ = toBool(ctx.DefaultQuery(aType.WithWriters, aDefault.WithWriters))
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	if value, ok := ctx.GetQuery(aType.SourceUUID); ok {
		args.SourceUUID = &value
	}
	if value, ok := ctx.GetQuery(aType.FlowNetworkCloneUUID); ok {
		args.FlowNetworkCloneUUID = &value
	}
	return args
}

func buildConsumerArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithWriters, _ = toBool(ctx.DefaultQuery(aType.WithWriters, aDefault.WithWriters))
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	return args
}

func buildProducerArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithWriterClones, _ = toBool(ctx.DefaultQuery(aType.WithWriterClones, aDefault.WithWriterClones))
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	if value, ok := ctx.GetQuery(aType.StreamUUID); ok {
		args.StreamUUID = &value
	}
	return args
}

func buildNetworkArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithDevices, _ = toBool(ctx.DefaultQuery(aType.WithDevices, aDefault.WithDevices))
	args.WithPoints, _ = toBool(ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints))
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	return args
}

func buildDeviceArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithPriority, _ = toBool(ctx.DefaultQuery(aType.WithPriority, aDefault.WithPriority))
	args.WithPoints, _ = toBool(ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints))
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	return args
}

func buildPointArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithPriority, _ = toBool(ctx.DefaultQuery(aType.WithPriority, aDefault.WithPriority))
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	return args
}

func buildWriterArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	if value, ok := ctx.GetQuery(aType.WriterThingClass); ok {
		args.WriterThingClass = &value
	}
	return args
}

func buildWriterCloneArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	if value, ok := ctx.GetQuery(aType.WriterThingClass); ok {
		args.WriterThingClass = &value
	}
	return args
}

func buildPluginArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.PluginName, _ = toBool(ctx.DefaultQuery(aType.PluginName, aDefault.PluginName))
	return args
}

func buildTagArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.Networks, _ = toBool(ctx.DefaultQuery(aType.WithNetworks, aDefault.WithNetworks))
	args.WithDevices, _ = toBool(ctx.DefaultQuery(aType.WithDevices, aDefault.WithDevices))
	args.WithPoints, _ = toBool(ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints))
	args.WithStreams, _ = toBool(ctx.DefaultQuery(aType.WithStreams, aDefault.WithStreams))
	args.WithProducers, _ = toBool(ctx.DefaultQuery(aType.WithProducers, aDefault.WithProducers))
	args.WithConsumers, _ = toBool(ctx.DefaultQuery(aType.WithConsumers, aDefault.WithConsumers))
	return args
}

func buildProducerHistoryArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	if value, ok := ctx.GetQuery(aType.IdGt); ok {
		args.IdGt = &value
	}
	if value, ok := ctx.GetQuery(aType.TimestampGt); ok {
		args.TimestampGt = &value
	}
	if value, ok := ctx.GetQuery(aType.TimestampLt); ok {
		args.TimestampLt = &value
	}
	if order, ok := ctx.GetQuery(aType.Order); ok {
		args.Order = order
	}
	return args
}
