package api

import (
	"github.com/gin-gonic/gin"
)

func buildNetworkArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithDevices, _ = toBool(ctx.DefaultQuery(aType.WithDevices, aDefault.WithDevices))
	args.WithPoints, _ = toBool(ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints))
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	args.WithMetaTags, _ = toBool(ctx.DefaultQuery(aType.WithMetaTags, aDefault.WithMetaTags))
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
