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

	plg := "plg" // plugin
	net := "net" // network
	dev := "dev" // device
	pnt := "pnt" // point
	job := "job" // job
	sch := "sch" // schedule
	ing := "ing" // integration

	alt := "alt" // alerts
	rub := "rbx" // rubix uuid
	rxg := "rxg" // rubix global uuid
	tok := "tok" // token uuid
	loc := "loc" // location uuid
	grp := "grp" // group uuid
	hos := "hos" // host uuid
	hoc := "hoc" // host comment uuid
	scl := "scl" // snapshot create log uuid
	srl := "srl" // snapshot restore log uuid
	mem := "mem" // member uuid
	med := "med" // member device uuid
	viw := "viw" // view uuid
	vwi := "vwi" // view widget uuid
	vse := "vse" // view setting uuid
	vte := "vte" // view template uuid
	vtw := "vtw" // view template widget uuid
	vtp := "vtp" // view template widget pointer uuid
	tem := "tem" // team uuid

	switch attribute {
	case model.CommonNaming.Plugin:
		return fmt.Sprintf("%s%s%s", plg, divider, u)
	case model.ThingClass.Network:
		return fmt.Sprintf("%s%s%s", net, divider, u)
	case model.ThingClass.Device:
		return fmt.Sprintf("%s%s%s", dev, divider, u)
	case model.ThingClass.Point:
		return fmt.Sprintf("%s%s%s", pnt, divider, u)
	case model.CommonNaming.Job:
		return fmt.Sprintf("%s%s%s", job, divider, u)
	case model.CommonNaming.Schedule:
		return fmt.Sprintf("%s%s%s", sch, divider, u)
	case model.ThingClass.Integration:
		return fmt.Sprintf("%s%s%s", ing, divider, u)
	case model.ThingClass.Alert:
		return fmt.Sprintf("%s%s%s", alt, divider, u)
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
	case model.CommonNaming.Member:
		return fmt.Sprintf("%s%s%s", mem, divider, u)
	case model.CommonNaming.MemberDevice:
		return fmt.Sprintf("%s%s%s", med, divider, u)
	case model.CommonNaming.View:
		return fmt.Sprintf("%s%s%s", viw, divider, u)
	case model.CommonNaming.ViewWidget:
		return fmt.Sprintf("%s%s%s", vwi, divider, u)
	case model.CommonNaming.ViewSetting:
		return fmt.Sprintf("%s%s%s", vse, divider, u)
	case model.CommonNaming.ViewTemplate:
		return fmt.Sprintf("%s%s%s", vte, divider, u)
	case model.CommonNaming.ViewTemplateWidget:
		return fmt.Sprintf("%s%s%s", vtw, divider, u)
	case model.CommonNaming.ViewTemplateWidgetPointer:
		return fmt.Sprintf("%s%s%s", vtp, divider, u)
	case model.CommonNaming.Team:
		return fmt.Sprintf("%s%s%s", tem, divider, u)
	}
	return u
}

func ShortUUID(prefix ...string) string {
	u := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, u)
	if n != len(u) || err != nil {
		return "-error-uuid-"
	}
	uuid_ := fmt.Sprintf("%x%x", u[0:4], u[4:6])
	if len(prefix) > 0 {
		uuid_ = fmt.Sprintf("%s_%s", prefix[0], uuid_)
	}
	return uuid_
}
