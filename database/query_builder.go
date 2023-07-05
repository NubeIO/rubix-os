package database

import (
	argspkg "github.com/NubeIO/rubix-os/args"
	"gorm.io/gorm"
	"strings"
)

func (d *GormDatabase) buildNetworkQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithDevices {
		query = query.Preload("Devices")
		if args.WithTags {
			query = query.Preload("Devices.Tags")
		}
		if args.WithMetaTags {
			query = query.Preload("Devices.MetaTags")
		}
	}
	if args.WithPoints {
		query = query.Preload("Devices.Points").Preload("Devices.Points.Priority")
		if args.WithTags {
			query = query.Preload("Devices.Points.Tags")
		}
		if args.WithMetaTags {
			query = query.Preload("Devices.Points.MetaTags")
		}
	}
	if args.WithTags {
		query = query.Preload("Tags")
	}
	if args.WithMetaTags {
		query = query.Preload("MetaTags")
	}
	if args.MetaTags != nil {
		keyValues := metaTagsArgsToKeyValues(*args.MetaTags)
		subQuery := d.DB.Table("network_meta_tags").Select("network_uuid").
			Where("(key, value) IN ?", keyValues).
			Group("network_uuid").
			Having("COUNT(network_uuid) = ?", len(keyValues))
		query = query.Where("uuid IN (?)", subQuery)
	}
	if args.GlobalUUID != nil {
		query = query.Where("global_uuid = ?", *args.GlobalUUID)
	}
	if args.PointSourceUUID != nil || args.HostUUID != nil {
		subQuery := d.DB.Table("networks").Select("networks.uuid").
			Joins("JOIN devices ON devices.network_uuid = networks.uuid").
			Joins("JOIN points ON points.device_uuid = devices.uuid")
		if args.PointSourceUUID != nil {
			subQuery = subQuery.Where("points.source_uuid = ?", *args.PointSourceUUID)
		}
		if args.HostUUID != nil {
			subQuery = subQuery.Where("networks.host_uuid = ?", *args.HostUUID)
		}
		query = query.Where("uuid IN (?)", subQuery)
		args.ShowCloneNetworks = true
	}
	if !args.ShowCloneNetworks {
		query = query.Where("is_clone IS NOT TRUE") // to support older data where is_clone gets default NULL
	}
	return query
}

func (d *GormDatabase) buildDeviceQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithPoints {
		query = query.Preload("Points")
		if args.WithTags {
			query = query.Preload("Points.Tags")
		}
		if args.WithMetaTags {
			query = query.Preload("Points.MetaTags")
		}
	}
	if args.WithTags {
		query = query.Preload("Tags")
	}
	if args.WithPriority {
		query = query.Preload("Points.Priority")
	}
	if args.AddressUUID != nil {
		query = query.Where("address_uuid = ?", *args.AddressUUID)
	}
	if args.Name != nil {
		query = query.Where("name = ?", *args.Name)
	}
	if args.NetworkUUID != nil {
		query = query.Where("network_uuid = ?", args.NetworkUUID)
	}
	if args.WithMetaTags {
		query = query.Preload("MetaTags")
	}
	if args.MetaTags != nil {
		keyValues := metaTagsArgsToKeyValues(*args.MetaTags)
		subQuery := d.DB.Table("device_meta_tags").Select("device_uuid").
			Where("(key, value) IN ?", keyValues).
			Group("device_uuid").
			Having("COUNT(device_uuid) = ?", len(keyValues))
		query = query.Where("uuid IN (?)", subQuery)
	}
	if args.PointSourceUUID != nil || args.HostUUID != nil {
		subQuery := d.DB.Table("networks").Select("devices.uuid").
			Joins("JOIN devices ON devices.network_uuid = networks.uuid").
			Joins("JOIN points ON points.device_uuid = devices.uuid")
		if args.PointSourceUUID != nil {
			subQuery = subQuery.Where("points.source_uuid = ?", *args.PointSourceUUID)
		}
		if args.HostUUID != nil {
			subQuery = subQuery.Where("networks.host_uuid = ?", *args.HostUUID)
		}
		query = query.Where("uuid IN (?)", subQuery)
	}
	return query
}

func (d *GormDatabase) buildPointQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithPriority {
		query = query.Preload("Priority")
	}
	if args.WithTags {
		query = query.Preload("Tags")
	}
	if args.AddressUUID != nil {
		query = query.Where("address_uuid = ?", *args.AddressUUID)
	}
	if args.IoNumber != nil {
		query = query.Where("io_number = ?", *args.IoNumber)
	}
	if args.AddressID != nil {
		query = query.Where("address_id = ?", *args.AddressID)
	}
	if args.ObjectType != nil {
		query = query.Where("object_type = ?", *args.ObjectType)
	}
	if args.DeviceUUID != nil {
		query = query.Where("device_uuid = ?", *args.DeviceUUID)
	}
	if args.WithMetaTags {
		query = query.Preload("MetaTags")
	}
	if args.SourceUUID != nil {
		query = query.Where("source_uuid = ?", *args.SourceUUID)
	}
	if args.MetaTags != nil {
		keyValues := metaTagsArgsToKeyValues(*args.MetaTags)
		subQuery := d.DB.Table("point_meta_tags").Select("point_uuid").
			Where("(key, value) IN ?", keyValues).
			Group("point_uuid").
			Having("COUNT(point_uuid) = ?", len(keyValues))
		query = query.Where("uuid IN (?)", subQuery)
	}
	if args.PointSourceUUID != nil || args.HostUUID != nil {
		subQuery := d.DB.Table("networks").Select("points.uuid").
			Joins("JOIN devices ON devices.network_uuid = networks.uuid").
			Joins("JOIN points ON points.device_uuid = devices.uuid")
		if args.PointSourceUUID != nil {
			subQuery = subQuery.Where("points.source_uuid = ?", *args.PointSourceUUID)
		}
		if args.HostUUID != nil {
			subQuery = subQuery.Where("networks.host_uuid = ?", *args.HostUUID)
		}
		query = query.Where("uuid IN (?)", subQuery)
	}
	return query
}

func (d *GormDatabase) buildTagQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.Networks {
		query = query.Preload("Networks")
	}
	if args.WithDevices {
		query = query.Preload("Devices")
	}
	if args.WithPoints {
		query = query.Preload("Points")
	}
	return query
}

func (d *GormDatabase) buildPointHistoryQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.IdGt != nil {
		query = query.Where("Id > ?", args.IdGt)
	}
	if args.TimestampGt != nil {
		query = query.Where("timestamp > datetime(?)", args.TimestampGt)
	}
	if args.TimestampLt != nil {
		query = query.Where("timestamp < datetime(?)", args.TimestampLt)
	}
	if args.Order != "" {
		order := strings.ToUpper(strings.TrimSpace(args.Order))
		if order != "ASC" && order != "DESC" {
			args.Order = "DESC"
		}
		query = query.Order("timestamp " + args.Order)
	}
	return query
}

func (d *GormDatabase) buildHistoryQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.TimestampGt != nil {
		query = query.Where("timestamp > datetime(?)", args.TimestampGt)
	}
	if args.TimestampLt != nil {
		query = query.Where("timestamp < datetime(?)", args.TimestampLt)
	}
	if args.Order != "" {
		order := strings.ToUpper(strings.TrimSpace(args.Order))
		if order != "ASC" && order != "DESC" {
			args.Order = "DESC"
		}
		query = query.Order("timestamp " + args.Order)
	}
	return query
}

func (d *GormDatabase) buildScheduleQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.Name != nil {
		query = query.Where("name = ?", *args.Name)
	}
	if args.GlobalUUID != nil {
		query = query.Where("global_uuid = ?", *args.GlobalUUID)
	}
	return query
}

func (d *GormDatabase) buildLocationQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithViews {
		query = query.Preload("Views")
	}
	if args.WithGroups {
		query = query.Preload("Groups")
		if args.WithViews {
			query = query.Preload("Groups.Views")
		}
		if args.WithHosts {
			query = query.Preload("Groups.Hosts")
			if args.WithViews {
				query = query.Preload("Groups.Hosts.Views")
			}
		}
	}
	return query
}

func (d *GormDatabase) buildGroupQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithViews {
		query = query.Preload("Views")
	}
	if args.WithHosts {
		query = query.Preload("Hosts")
		if args.WithViews {
			query = query.Preload("Hosts.Views")
		}
	}
	return query
}

func (d *GormDatabase) buildHostQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithTags {
		query = query.Preload("Tags")
	}
	if args.WithComments {
		query = query.Preload("Comments")
	}
	if args.WithViews {
		query = query.Preload("Views")
	}
	if args.Name != nil {
		query = query.Where("name = ?", *args.Name)
	}
	return query
}

func (d *GormDatabase) buildMemberQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithMemberDevices {
		query = query.Preload("MemberDevices")
	}
	if args.WithTeams {
		query = query.Preload("Teams")
	}
	return query
}

func (d *GormDatabase) buildMemberDeviceQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.MemberUUID != nil {
		query = query.Where("member_uuid", *args.MemberUUID)
	}
	if args.DeviceId != nil {
		query = query.Where("device_id", *args.DeviceId)
	}
	return query
}

func (d *GormDatabase) buildTeamQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithMembers {
		query = query.Preload("Members")
	}
	if args.WithViews {
		query = query.Preload("Views")
	}
	return query
}

func (d *GormDatabase) buildViewQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithWidgets {
		query = query.Preload("Widgets")
	}
	return query
}

func (d *GormDatabase) buildViewTemplateQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithViewTemplateWidgets {
		query = query.Preload("ViewTemplateWidgets")
		if args.WithViewTemplateWidgetPointers {
			query = query.Preload("ViewTemplateWidgets.ViewTemplateWidgetPointers")
		}
	}
	return query
}
func (d *GormDatabase) buildViewTemplateWidgetQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithViewTemplateWidgetPointers {
		query = query.Preload("ViewTemplateWidgetPointers")
	}
	return query
}

func (d *GormDatabase) buildTicketQuery(args argspkg.Args) *gorm.DB {
	query := d.DB
	if args.WithComments {
		query = query.Preload("Comments")
	}
	if args.WithTeams {
		query = query.Preload("Teams")
	}
	return query
}
