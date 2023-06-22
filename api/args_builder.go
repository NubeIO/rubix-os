package api

import (
	"github.com/NubeIO/rubix-os/utils/boolean"
	"github.com/gin-gonic/gin"
)

func buildFlowNetworkArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithStreams, _ = toBool(ctx.DefaultQuery(aType.WithStreams, aDefault.WithStreams))
	args.WithProducers, _ = toBool(ctx.DefaultQuery(aType.WithProducers, aDefault.WithProducers))
	args.WithCommandGroups, _ = toBool(ctx.DefaultQuery(aType.WithCommandGroups, aDefault.WithCommandGroups))
	args.WithWriterClones, _ = toBool(ctx.DefaultQuery(aType.WithWriterClones, aDefault.WithWriterClones))
	args.IsMetadata, _ = toBool(ctx.DefaultQuery(aType.IsMetadata, aDefault.IsMetadata))
	if val, exists := ctx.Get(aType.IsRemote); exists {
		args.IsRemote = boolean.New(val.(bool))
	}
	if value, ok := ctx.GetQuery(aType.PluginName); ok {
		args.PluginName = value
	}
	if value, ok := ctx.GetQuery(aType.FlowNetworkUUID); ok {
		args.FlowNetworkUUID = value
	}
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
	if value, ok := ctx.GetQuery(aType.SourceUUID); ok {
		args.SourceUUID = &value
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
	args.IsMetadata, _ = toBool(ctx.DefaultQuery(aType.IsMetadata, aDefault.IsMetadata))
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
	if value, ok := ctx.GetQuery(aType.FlowNetworkUUID); ok {
		args.FlowNetworkUUID = value
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
	if value, ok := ctx.GetQuery(aType.ProducerUUID); ok {
		args.ProducerUUID = &value
	}
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
	if value, ok := ctx.GetQuery(aType.Name); ok {
		args.Name = &value
	}
	if value, ok := ctx.GetQuery(aType.ProducerThingUUID); ok {
		args.ProducerThingUUID = &value
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
	args.WithMetaTags, _ = toBool(ctx.DefaultQuery(aType.WithMetaTags, aDefault.WithMetaTags))
	if value, ok := ctx.GetQuery(aType.FlowNetworkUUID); ok {
		args.FlowNetworkUUID = value
	}
	if value, ok := ctx.GetQuery(aType.MetaTags); ok {
		args.MetaTags = &value
	}
	args.ShowCloneNetworks, _ = toBool(ctx.DefaultQuery(aType.ShowCloneNetworks, aDefault.ShowCloneNetworks))
	return args
}

func buildDeviceArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithPriority, _ = toBool(ctx.DefaultQuery(aType.WithPriority, aDefault.WithPriority))
	args.WithPoints, _ = toBool(ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints))
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	args.WithMetaTags, _ = toBool(ctx.DefaultQuery(aType.WithMetaTags, aDefault.WithMetaTags))
	if value, ok := ctx.GetQuery(aType.AddressUUID); ok {
		args.AddressUUID = &value
	}
	if value, ok := ctx.GetQuery(aType.MetaTags); ok {
		args.MetaTags = &value
	}
	return args
}

func buildPointArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithPriority, _ = toBool(ctx.DefaultQuery(aType.WithPriority, aDefault.WithPriority))
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	args.WithMetaTags, _ = toBool(ctx.DefaultQuery(aType.WithMetaTags, aDefault.WithMetaTags))
	if value, ok := ctx.GetQuery(aType.AddressUUID); ok {
		args.AddressUUID = &value
	}
	if value, ok := ctx.GetQuery(aType.IoNumber); ok {
		args.IoNumber = &value
	}
	if value, ok := ctx.GetQuery(aType.AddressID); ok {
		args.AddressID = &value
	}
	if value, ok := ctx.GetQuery(aType.ObjectType); ok {
		args.ObjectType = &value
	}
	if value, ok := ctx.GetQuery(aType.MetaTags); ok {
		args.MetaTags = &value
	}
	return args
}

func buildWriterArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	if value, ok := ctx.GetQuery(aType.ConsumerUUID); ok {
		args.ConsumerUUID = &value
	}
	if value, ok := ctx.GetQuery(aType.WriterThingClass); ok {
		args.WriterThingClass = &value
	}
	if value, ok := ctx.GetQuery(aType.WriterThingUUID); ok {
		args.WriterThingUUID = &value
	}
	return args
}

func buildWriterCloneArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	if value, ok := ctx.GetQuery(aType.ProducerUUID); ok {
		args.ProducerUUID = &value
	}
	if value, ok := ctx.GetQuery(aType.WriterThingClass); ok {
		args.WriterThingClass = &value
	}
	if value, ok := ctx.GetQuery(aType.SourceUUID); ok {
		args.SourceUUID = &value
	}
	return args
}

func buildPluginArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.ByPluginName, _ = toBool(ctx.DefaultQuery(aType.ByPluginName, aDefault.PluginName))
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

func buildPointHistoryArgs(ctx *gin.Context) Args {
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

func buildPointHistorySyncArgs(ctx *gin.Context) (string, string) {
	id := ""
	timeStamp := ""
	if value, ok := ctx.GetQuery("id"); ok {
		id = value
	}
	if value, ok := ctx.GetQuery("timestamp"); ok {
		timeStamp = value
	}
	return id, timeStamp
}

func buildScheduleArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	if value, ok := ctx.GetQuery(aType.Name); ok {
		args.Name = &value
	}
	return args
}
