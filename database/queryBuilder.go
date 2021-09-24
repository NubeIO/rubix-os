package database

import (
	"github.com/NubeDev/flow-framework/api"
	"gorm.io/gorm"
	"strings"
)

func (d *GormDatabase) buildFlowNetworkQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Streams {
		query = query.Preload("Streams")
		if args.Producers {
			query = query.Preload("Streams.Producers")
			if args.Writers {
				query = query.Preload("Streams.Producers.WriterClones")
			}
		}
		if args.Consumers {
			query = query.Preload("Streams.Consumers")
			if args.Writers {
				query = query.Preload("Streams.Consumers.Writers")
			}
		}
		if args.CommandGroups {
			query = query.Preload("Streams.CommandGroups")
		}
	}
	if args.GlobalUUID != nil {
		query = query.Where("global_uuid = ?", *args.GlobalUUID)
	}
	if args.ClientId != nil {
		query = query.Where("client_id = ?", *args.ClientId)
	}
	if args.SiteId != nil {
		query = query.Where("site_id = ?", *args.SiteId)
	}
	if args.DeviceId != nil {
		query = query.Where("device_id = ?", *args.DeviceId)
	}
	return query
}

func (d *GormDatabase) buildStreamQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.FlowNetworks {
		query = query.Preload("FlowNetworks")
	}
	if args.Producers {
		query = query.Preload("Producers")
		if args.Writers {
			query = query.Preload("Producers.WriterClones")
		}
	}
	if args.Consumers {
		query = query.Preload("Consumers")
		if args.Writers {
			query = query.Preload("Consumers.Writers")
		}
	}
	if args.CommandGroups {
		query = query.Preload("CommandGroups")
	}
	if args.Tags {
		query = query.Preload("Tags")
	}
	return query
}

func (d *GormDatabase) buildConsumerQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Writers {
		query = query.Preload("Writers")
	}
	if args.Tags {
		query = query.Preload("Tags")
	}
	return query
}

func (d *GormDatabase) buildProducerQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Writers {
		query = query.Preload("WriterClones")
	}
	if args.Tags {
		query = query.Preload("Tags")
	}
	return query
}

func (d *GormDatabase) buildNetworkQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Devices {
		query = query.Preload("Devices")
	}
	if args.Points {
		query = query.Preload("Devices.Points")
	}
	if args.IpConnection {
		query = query.Preload("IpConnection")
	}
	if args.SerialConnection {
		query = query.Preload("SerialConnection")
	}
	if args.Tags {
		query = query.Preload("Tags")
	}
	return query
}

func (d *GormDatabase) buildDeviceQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Points {
		query = query.Preload("Points")
	}
	if args.Tags {
		query = query.Preload("Tags")
	}
	return query
}

func (d *GormDatabase) buildPointQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Priority {
		query = query.Preload("Priority")
	}
	if args.Tags {
		query = query.Preload("Tags")
	}
	return query
}

func (d *GormDatabase) buildTagQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Networks {
		query = query.Preload("Networks")
	}
	if args.Devices {
		query = query.Preload("Devices")
	}
	if args.Points {
		query = query.Preload("Points")
	}
	if args.Streams {
		query = query.Preload("Streams")
	}
	if args.Producers {
		query = query.Preload("Producers")
	}
	if args.Consumers {
		query = query.Preload("Consumers")
	}
	return query
}

func (d *GormDatabase) buildTagQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Devices {
		query = query.Preload("Devices")
	}
	if args.Points {
		query = query.Preload("Points")
	}
	if args.Streams {
		query = query.Preload("Streams")
	}
	if args.Producers {
		query = query.Preload("Producers")
	}
	if args.Consumers {
		query = query.Preload("Consumers")
	}
	return query
}

func (d *GormDatabase) buildProducerHistoryQuery(args api.Args) *gorm.DB {
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