package api

import "github.com/gin-gonic/gin"

func networkArgs(ctx *gin.Context) (withChildren bool, withPoints bool){
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithChildren = ctx.DefaultQuery(aType.WithChildren, aDefault.WithChildren)
	args.WithPoints = ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints)
	withChildren, _ = toBool(args.WithChildren) //?with_children=true&points=true
	withPoints, _ = toBool(args.WithPoints)
	return withChildren, withPoints

}
//withChildrenArgs
func withChildrenArgs(ctx *gin.Context) (withChildren bool, withPoints bool){
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.WithChildren = ctx.DefaultQuery(aType.WithChildren, aDefault.WithChildren)
	args.WithPoints = ctx.DefaultQuery(aType.WithPoints, aDefault.WithPoints)
	withChildren, _ = toBool(args.WithChildren) //?with_children=true&points=true
	withPoints, _ = toBool(args.WithPoints)
	return withChildren, withPoints
}

//withConsumerArgs
func withConsumerArgs(ctx *gin.Context) (askResponse bool, askRefresh bool, write bool, thingType string, flowNetworkUUID string){
	var args Args
	var aType = ArgsType
	var aDefault = ArgsDefault
	args.AskRefresh = ctx.DefaultQuery(aType.AskRefresh, aDefault.AskRefresh)
	args.AskResponse = ctx.DefaultQuery(aType.AskResponse, aDefault.AskResponse)
	args.Write = ctx.DefaultQuery(aType.Write, aDefault.Write)
	args.ThingType = ctx.DefaultQuery(aType.ThingType, aDefault.ThingType)
	args.FlowNetworkUUID = ctx.DefaultQuery(aType.FlowNetworkUUID, aDefault.FlowNetworkUUID)
	askRefresh, _ = toBool(args.AskRefresh)
	askResponse, _ = toBool(args.AskResponse)
	write, _ = toBool(args.Write)
	return askRefresh, askResponse, write, args.ThingType, args.FlowNetworkUUID
}

