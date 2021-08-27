package utils
import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uuid"
)


func MakeUUID() (string, error) {
	return uuid.MakeUUID()
}


func MakeTopicUUID(attribute string) string {
	u, _ := uuid.MakeUUID()
	divider := "_"
	fln := "fln" //flow network
	plg := "plg" //plugin
	net := "net" //network
	dev := "dev" //device
	pnt := "pnt" //point
	job := "job" //job
	str := "str" //stream gateway
	sub := "sub" //subscribers
	sus := "sus" //subscriptions
	sul := "sul" //subscriptions
	alt := "alt" //alerts
	cmd := "cmd" //command
	rub := "rbx" //rubix uuid
	rxg := "rxg" //rubix global uuid
	switch attribute {
	case model.CommonNaming.Plugin:
		return fmt.Sprintf("%s%s%s", plg, divider, u)
	case model.CommonNaming.FlowNetwork:
		return fmt.Sprintf("%s%s%s", fln, divider, u)
	case model.CommonNaming.Network:
		return fmt.Sprintf("%s%s%s", net, divider, u)
	case model.CommonNaming.Device:
		return fmt.Sprintf("%s%s%s", dev, divider, u)
	case model.CommonNaming.Point:
		return fmt.Sprintf("%s%s%s", pnt, divider, u)
	case model.CommonNaming.Stream:
		return fmt.Sprintf("%s%s%s", str, divider, u)
	case model.CommonNaming.Job:
		return fmt.Sprintf("%s%s%s", job, divider, u)
	case model.CommonNaming.Subscriber:
		return fmt.Sprintf("%s%s%s", sub, divider, u)
	case model.CommonNaming.Subscription:
		return fmt.Sprintf("%s%s%s", sus, divider, u)
	case model.CommonNaming.SubscriptionList:
		return fmt.Sprintf("%s%s%s", sul, divider, u)
	case model.CommonNaming.Alert:
		return fmt.Sprintf("%s%s%s", alt, divider, u)
	case model.CommonNaming.CommandGroup:
		return fmt.Sprintf("%s%s%s", cmd, divider, u)
	case model.CommonNaming.Rubix:
		return fmt.Sprintf("%s%s%s", rub, divider, u)
	case model.CommonNaming.RubixGlobal:
		return fmt.Sprintf("%s%s%s", rxg, divider, u)

	}
	return u
}

