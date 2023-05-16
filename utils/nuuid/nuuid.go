package nuuid

import (
	"crypto/rand"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"io"
)

func MakeUUID() (string, error) {
	return uuid.MakeUUID()
}

func MakeTopicUUID(attribute string) string {
	u, _ := uuid.MakeUUID()
	divider := "_"

	fln := "fln" // flow network
	str := "str" // stream
	stc := "stc" // stream clone
	pro := "pro" // producers
	prh := "prh" // producer history
	wrc := "wrc" // writerClone
	con := "con" // consumers
	wri := "wri" // writer

	fnc := "rfn" // flow network clone
	plg := "plg" // plugin
	net := "net" // network
	dev := "dev" // device
	pnt := "pnt" // point
	job := "job" // job
	sch := "sch" // schedule
	ing := "ing" // integration

	stl := "stl" // list of flow network gateway
	alt := "alt" // alerts
	cmd := "cmd" // command
	rub := "rbx" // rubix uuid
	rxg := "rxg" // rubix global uuid
	tok := "tok" // token uuid
	loc := "loc" // location uuid
	grp := "grp" // group uuid
	hos := "hos" // host uuid
	hoc := "hoc" // host comment uuid
	scl := "scl" // snapshot create log uuid
	srl := "srl" // snapshot restore log uuid
	tea := "tea" // team uuid

	switch attribute {
	case model.CommonNaming.Plugin:
		return fmt.Sprintf("%s%s%s", plg, divider, u)
	case model.CommonNaming.FlowNetwork:
		return fmt.Sprintf("%s%s%s", fln, divider, u)
	case model.CommonNaming.FlowNetworkClone:
		return fmt.Sprintf("%s%s%s", fnc, divider, u)
	case model.ThingClass.Network:
		return fmt.Sprintf("%s%s%s", net, divider, u)
	case model.ThingClass.Device:
		return fmt.Sprintf("%s%s%s", dev, divider, u)
	case model.ThingClass.Point:
		return fmt.Sprintf("%s%s%s", pnt, divider, u)
	case model.CommonNaming.Stream:
		return fmt.Sprintf("%s%s%s", str, divider, u)
	case model.CommonNaming.StreamClone:
		return fmt.Sprintf("%s%s%s", stc, divider, u)
	case model.CommonNaming.StreamList:
		return fmt.Sprintf("%s%s%s", stl, divider, u)
	case model.CommonNaming.Job:
		return fmt.Sprintf("%s%s%s", job, divider, u)
	case model.CommonNaming.Schedule:
		return fmt.Sprintf("%s%s%s", sch, divider, u)
	case model.CommonNaming.Producer:
		return fmt.Sprintf("%s%s%s", pro, divider, u)
	case model.ThingClass.Integration:
		return fmt.Sprintf("%s%s%s", ing, divider, u)
	case model.CommonNaming.WriterClone:
		return fmt.Sprintf("%s%s%s", wrc, divider, u)
	case model.CommonNaming.ProducerHistory:
		return fmt.Sprintf("%s%s%s", prh, divider, u)
	case model.CommonNaming.Consumer:
		return fmt.Sprintf("%s%s%s", con, divider, u)
	case model.CommonNaming.Writer:
		return fmt.Sprintf("%s%s%s", wri, divider, u)
	case model.ThingClass.Alert:
		return fmt.Sprintf("%s%s%s", alt, divider, u)
	case model.CommonNaming.CommandGroup:
		return fmt.Sprintf("%s%s%s", cmd, divider, u)
	case model.CommonNaming.Rubix:
		return fmt.Sprintf("%s%s%s", rub, divider, u)
	case model.CommonNaming.RubixGlobal:
		return fmt.Sprintf("%s%s%s", rxg, divider, u)
	case model.CommonNaming.Token:
		return fmt.Sprintf("%s%s%s", tok, divider, u)
	case model.CommonNaming.Location:
		return fmt.Sprintf("%s%s%s", loc, divider, u)
	case model.CommonNaming.Group:
		return fmt.Sprintf("%s%s%s", grp, divider, u)
	case model.CommonNaming.Host:
		return fmt.Sprintf("%s%s%s", hos, divider, u)
	case model.CommonNaming.HostComment:
		return fmt.Sprintf("%s%s%s", hoc, divider, u)
	case model.CommonNaming.SnapshotCreateLog:
		return fmt.Sprintf("%s%s%s", scl, divider, u)
	case model.CommonNaming.SnapshotRestoreLog:
		return fmt.Sprintf("%s%s%s", srl, divider, u)
	case model.CommonNaming.Team:
		return fmt.Sprintf("%s%s%s", tea, divider, u)
	}
	return u
}

func ShortUUID(prefix ...string) string {
	u := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, u)
	if n != len(u) || err != nil {
		return "-error-uuid-"
	}
	uuid := fmt.Sprintf("%x%x", u[0:4], u[4:6])
	if len(prefix) > 0 {
		uuid = fmt.Sprintf("%s_%s", prefix[0], uuid)
	}
	return uuid
}
