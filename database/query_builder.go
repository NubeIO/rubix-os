package database

import (
	"github.com/NubeIO/rubix-os/api"
	"gorm.io/gorm"
	"strings"
)

func (d *GormDatabase) buildNetworkQuery(args api.Args) *gorm.DB {
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
	if !args.ShowCloneNetworks {
		query = query.Where("is_clone IS NOT TRUE") // to support older data where is_clone gets default NULL
	}
	return query
}

func (d *GormDatabase) buildDeviceQuery(args api.Args) *gorm.DB {
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
	return query
}

func (d *GormDatabase) buildPointQuery(args api.Args) *gorm.DB {
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
	return query
}

func (d *GormDatabase) buildTagQuery(args api.Args) *gorm.DB {
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

func (d *GormDatabase) buildPointHistoryQuery(args api.Args) *gorm.DB {
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

func (d *GormDatabase) buildHistoryQuery(args api.Args) *gorm.DB {
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

func (d *GormDatabase) buildScheduleQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Name != nil {
		query = query.Where("name = ?", *args.Name)
	}
	if args.GlobalUUID != nil {
		query = query.Where("global_uuid = ?", *args.GlobalUUID)
	}
	return query
}

func (d *GormDatabase) buildLocationQuery() *gorm.DB {
	return d.DB.Preload("Views").Preload("Groups.Views").Preload("Groups.Hosts.Views")
}

func (d *GormDatabase) buildGroupQuery() *gorm.DB {
	return d.DB.Preload("Views").Preload("Hosts.Views")
}

func (d *GormDatabase) buildHostQuery(args api.Args) *gorm.DB {
	query := d.DB.Preload("Comments").Preload("Tags").Preload("Views")
	if args.Name != nil {
		query = query.Where("name = ?", *args.Name)
	}
	return query
}

func (d *GormDatabase) buildMemberQuery(args api.Args) *gorm.DB {
	query := d.DB.Preload("MemberDevices")
	if args.Name != nil {
		query = query.Where("name = ?", *args.Name)
	}
	return query
}

func (d *GormDatabase) buildMemberDeviceQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.MemberUUID != nil {
		query = query.Where("member_uuid", *args.MemberUUID)
	}
	if args.DeviceId != nil {
		query = query.Where("device_id", *args.DeviceId)
	}
	return query
}

func (d *GormDatabase) buildTeamQuery() *gorm.DB {
	return d.DB.Preload("Members").Preload("Views")
}

func (d *GormDatabase) buildViewQuery() *gorm.DB {
	return d.DB.Preload("Widgets")
}

func (d *GormDatabase) buildViewTemplateQuery() *gorm.DB {
	return d.DB.Preload("ViewTemplateWidgets.ViewTemplateWidgetPointers")
}

func (d *GormDatabase) buildViewTemplateWidgetQuery() *gorm.DB {
	return d.DB.Preload("ViewTemplateWidgetPointers")
}
