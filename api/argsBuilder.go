package api

import "github.com/gin-gonic/gin"

//withChildrenArgs
func streamFieldsArgs(ctx *gin.Context) (flowUUID string, streamUUID string, producerUUID string, consumerUUID string, writerUUID string) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.FlowUUID = ctx.DefaultQuery(aType.FlowUUID, aDefault.FlowUUID)
	args.StreamUUID = ctx.DefaultQuery(aType.StreamUUID, aDefault.StreamUUID)
	args.ProducerUUID = ctx.DefaultQuery(aType.ProducerUUID, aDefault.ProducerUUID)
	args.ConsumerUUID = ctx.DefaultQuery(aType.ConsumerUUID, aDefault.ConsumerUUID)
	args.WriterUUID = ctx.DefaultQuery(aType.WriterUUID, aDefault.WriterUUID)
	return args.FlowUUID, args.StreamUUID, args.ProducerUUID, args.ConsumerUUID, args.WriterUUID
}

//withChildrenArgs
func withFieldsArgs(ctx *gin.Context) (field string, value string) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.Field = ctx.DefaultQuery(aType.Field, aDefault.Field)
	args.Value = ctx.DefaultQuery(aType.Value, aDefault.Value)
	return args.Field, args.Value
}

//withChildrenArgs
func parentArgs(ctx *gin.Context) (AddToParent string) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.AddToParent = ctx.DefaultQuery(aType.AddToParent, aDefault.AddToParent)
	return args.AddToParent
}

//withChildrenArgs
func queryFields(ctx *gin.Context) (order string, limit string) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault //ASC or DESC
	args.Order = ctx.DefaultQuery(aType.Order, aDefault.Order)
	args.Limit = ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints)
	return args.Order, args.Limit
}

//withChildrenArgs
func withChildrenArgs(ctx *gin.Context) (withChildren bool, withPoints bool, withParent bool) {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithChildren = ctx.DefaultQuery(aType.WithChildren, aDefault.WithChildren)
	args.WithPoints = ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints)
	args.WithParent = ctx.DefaultQuery(aType.WithParent, aDefault.WithParent)
	withChildren, _ = toBool(args.WithChildren) //?with_children=true&points=true
	withPoints, _ = toBool(args.WithPoints)
	withParent, _ = toBool(args.WithParent)
	return withChildren, withPoints, withParent
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
	args.Streams, _ = toBool(ctx.DefaultQuery(aType.Streams, aDefault.Streams))
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
	args.FlowNetworks, _ = toBool(ctx.DefaultQuery(aType.FlowNetworks, aDefault.FlowNetworks))
	args.Producers, _ = toBool(ctx.DefaultQuery(aType.Producers, aDefault.Producers))
	args.Consumers, _ = toBool(ctx.DefaultQuery(aType.Consumers, aDefault.Consumers))
	args.CommandGroups, _ = toBool(ctx.DefaultQuery(aType.CommandGroups, aDefault.CommandGroups))
	args.Writers, _ = toBool(ctx.DefaultQuery(aType.Writers, aDefault.Writers))
	args.Tags, _ = toBool(ctx.DefaultQuery(aType.Tags, aDefault.Tags))
	return args
}

func buildConsumerArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.Writers, _ = toBool(ctx.DefaultQuery(aType.Writers, aDefault.Writers))
	return args
}

func buildProducerArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.Writers, _ = toBool(ctx.DefaultQuery(aType.Writers, aDefault.Writers))
	return args
}

func buildNetworkArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.Devices, _ = toBool(ctx.DefaultQuery(aType.Devices, aDefault.Devices))
	args.Points, _ = toBool(ctx.DefaultQuery(aType.Points, aDefault.Points))
	args.IpConnection, _ = toBool(ctx.DefaultQuery(aType.IpConnection, aDefault.IpConnection))
	args.SerialConnection, _ = toBool(ctx.DefaultQuery(aType.SerialConnection, aDefault.SerialConnection))
	return args
}

func buildDeviceArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.Points, _ = toBool(ctx.DefaultQuery(aType.Points, aDefault.Points))
	args.IpConnection, _ = toBool(ctx.DefaultQuery(aType.IpConnection, aDefault.IpConnection))
	args.SerialConnection, _ = toBool(ctx.DefaultQuery(aType.SerialConnection, aDefault.SerialConnection))
	return args
}

func buildPluginArgs(ctx *gin.Context) Args {
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.PluginName, _ = toBool(ctx.DefaultQuery(aType.PluginName, aDefault.PluginName))
	return args
}
