package api

import (
	argspkg "github.com/NubeIO/rubix-os/args"
	"github.com/gin-gonic/gin"
)

func buildNetworkArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
	args.WithDevices, _ = toBool(ctx.DefaultQuery(aType.WithDevices, aDefault.WithDevices))
	args.WithPoints, _ = toBool(ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints))
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	args.WithMetaTags, _ = toBool(ctx.DefaultQuery(aType.WithMetaTags, aDefault.WithMetaTags))
	if value, ok := ctx.GetQuery(aType.MetaTags); ok {
		args.MetaTags = &value
	}
	args.ShowCloneNetworks, _ = toBool(ctx.DefaultQuery(aType.ShowCloneNetworks, aDefault.ShowCloneNetworks))
	if value, ok := ctx.GetQuery(aType.PointSourceUUID); ok {
		args.PointSourceUUID = &value
	}
	if value, ok := ctx.GetQuery(aType.HostUUID); ok {
		args.HostUUID = &value
	}
	return args
}

func buildDeviceArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
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
	if value, ok := ctx.GetQuery(aType.PointSourceUUID); ok {
		args.PointSourceUUID = &value
	}
	if value, ok := ctx.GetQuery(aType.HostUUID); ok {
		args.HostUUID = &value
	}
	return args
}

func buildPointArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
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
	if value, ok := ctx.GetQuery(aType.PointSourceUUID); ok {
		args.PointSourceUUID = &value
	}
	if value, ok := ctx.GetQuery(aType.HostUUID); ok {
		args.HostUUID = &value
	}
	return args
}

func buildPluginArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
	args.ByPluginName, _ = toBool(ctx.DefaultQuery(aType.ByPluginName, aDefault.PluginName))
	return args
}

func buildTagArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
	args.Networks, _ = toBool(ctx.DefaultQuery(aType.WithNetworks, aDefault.WithNetworks))
	args.WithDevices, _ = toBool(ctx.DefaultQuery(aType.WithDevices, aDefault.WithDevices))
	args.WithPoints, _ = toBool(ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints))
	return args
}

func buildPointHistoryArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
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

func buildScheduleArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	if value, ok := ctx.GetQuery(aType.Name); ok {
		args.Name = &value
	}
	return args
}

func buildTeamArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
	args.WithMembers, _ = toBool(ctx.DefaultQuery(aType.WithMembers, aDefault.WithMembers))
	args.WithViews, _ = toBool(ctx.DefaultQuery(aType.WithViews, aDefault.WithViews))
	return args
}

func buildLocationArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
	args.WithViews, _ = toBool(ctx.DefaultQuery(aType.WithViews, aDefault.WithViews))
	args.WithGroups, _ = toBool(ctx.DefaultQuery(aType.WithGroups, aDefault.WithGroups))
	args.WithHosts, _ = toBool(ctx.DefaultQuery(aType.WithHosts, aDefault.WithHosts))
	return args
}

func buildGroupArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
	args.WithViews, _ = toBool(ctx.DefaultQuery(aType.WithViews, aDefault.WithViews))
	args.WithHosts, _ = toBool(ctx.DefaultQuery(aType.WithHosts, aDefault.WithHosts))
	return args
}

func buildHostArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
	args.WithTags, _ = toBool(ctx.DefaultQuery(aType.WithTags, aDefault.WithTags))
	args.WithComments, _ = toBool(ctx.DefaultQuery(aType.WithComments, aDefault.WithComments))
	args.WithViews, _ = toBool(ctx.DefaultQuery(aType.WithViews, aDefault.WithViews))
	if value, ok := ctx.GetQuery(aType.Name); ok {
		args.Name = &value
	}
	return args
}

func buildMemberArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
	args.WithMemberDevices, _ = toBool(ctx.DefaultQuery(aType.WithMemberDevices, aDefault.WithMemberDevices))
	args.WithTeams, _ = toBool(ctx.DefaultQuery(aType.WithTeams, aDefault.WithTeams))
	return args
}

func buildViewArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
	args.WithWidgets, _ = toBool(ctx.DefaultQuery(aType.WithWidgets, aDefault.WithWidgets))
	return args
}

func buildViewTemplateArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
	args.WithViewTemplateWidgets, _ = toBool(ctx.DefaultQuery(aType.WithViewTemplateWidgets,
		aDefault.WithViewTemplateWidgets))
	args.WithViewTemplateWidgetPointers, _ = toBool(ctx.DefaultQuery(aType.WithViewTemplateWidgetPointers,
		aDefault.WithViewTemplateWidgetPointers))
	return args
}

func buildTicketArgs(ctx *gin.Context) argspkg.Args {
	var args argspkg.Args
	var aType = argspkg.ArgsType
	var aDefault = argspkg.ArgsDefault
	args.WithTeams, _ = toBool(ctx.DefaultQuery(aType.WithTeams, aDefault.WithTeams))
	args.WithComments, _ = toBool(ctx.DefaultQuery(aType.WithComments, aDefault.WithComments))
	return args
}

func buildAlertArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithTickets, _ = toBool(ctx.DefaultQuery(aType.WithTickets, aDefault.WithTickets))
	if value, ok := ctx.GetQuery(aType.Target); ok {
		args.Target = &value
	}
	return args
}
