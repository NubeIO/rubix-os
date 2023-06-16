package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/api"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/src/cli/cligetter"
	"github.com/NubeIO/rubix-os/utils/nuuid"
)

func (d *GormDatabase) GetHost(uuid string) (*model.Host, error) {
	hostModel := model.Host{}
	query := d.buildHostQuery(api.Args{})
	if err := query.Where("uuid = ? ", uuid).First(&hostModel).Error; err != nil {
		return nil, errors.New(fmt.Sprintf("no host was found with uuid: %s", uuid))
	}
	return &hostModel, nil
}

func (d *GormDatabase) GetHostByName(name string) (*model.Host, error) {
	hostModel := model.Host{}
	query := d.buildHostQuery(api.Args{Name: &name})
	if err := query.First(&hostModel).Error; err != nil {
		return nil, errors.New(fmt.Sprintf("no host was found with name: %s", name))
	}
	return &hostModel, nil
}

func (d *GormDatabase) GetHosts(withOpenVPN bool) ([]*model.Host, error) {
	var hostsModel []*model.Host
	query := d.buildHostQuery(api.Args{})
	if err := query.Find(&hostsModel).Error; err != nil {
		return nil, err
	}
	if withOpenVPN {
		attachOpenVPN(hostsModel)
	}
	return hostsModel, nil
}

func (d *GormDatabase) GetFirstHost() (*model.Host, error) {
	hostsModel, err := d.GetHosts(false)
	if err != nil {
		return nil, err
	}
	if len(hostsModel) > 0 {
		return hostsModel[0], err
	}
	return nil, err
}

func (d *GormDatabase) GetHostsByUUIDs(uuids []*string) ([]*model.Host, error) {
	var groupsModel []*model.Host
	query := d.buildHostQuery(api.Args{})
	if err := query.Where("uuid IN ?", uuids).Find(&groupsModel).Error; err != nil {
		return nil, err
	}
	return groupsModel, nil
}

func (d *GormDatabase) CreateHost(body *model.Host) (*model.Host, error) {
	body.UUID = nuuid.MakeTopicUUID(model.CommonNaming.Host)
	if body.HTTPS == nil {
		body.HTTPS = nils.NewFalse()
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateHostByName(name string, body *model.Host) (*model.Host, error) {
	m := new(model.Host)
	query := d.DB.Where("name = ?", name).Find(&m).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return m, nil
}

func (d *GormDatabase) UpdateHost(uuid string, body *model.Host) (*model.Host, error) {
	m := new(model.Host)
	query := d.DB.Where("uuid = ?", uuid).Find(&m).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return m, nil
}

func (d *GormDatabase) DeleteHost(uuid string) (*interfaces.Message, error) {
	var m *model.Host
	query := d.DB.Where("uuid = ? ", uuid).Delete(&m)
	return d.deleteResponse(query)
}

func (d *GormDatabase) DropHosts() (*interfaces.Message, error) {
	var m *model.Host
	query := d.DB.Where("1 = 1").Delete(&m)
	return d.deleteResponse(query)
}

func (d *GormDatabase) ConfigureOpenVPN(uuid string) (*interfaces.Message, error) {
	host := model.Host{}
	if err := d.DB.Where("uuid = ? ", uuid).First(&host).Error; err != nil {
		return nil, errors.New(fmt.Sprintf("no host was found with uuid: %s", uuid))
	}
	cli := cligetter.GetEdgeClient(&host)
	globalUUID, _, pingable, isValidToken := cli.Ping()
	if pingable == false {
		return nil, errors.New("make it accessible at first")
	}
	if isValidToken == false || globalUUID == nil {
		return nil, errors.New("configure valid token at first")
	}
	host.GlobalUUID = *globalUUID
	oCli, err := cligetter.GetOpenVPNClient()
	if err != nil {
		return nil, err
	}
	openVPNConfig, err := oCli.GetOpenVPNConfig(host.GlobalUUID)
	if err != nil {
		return nil, err
	}
	_, err = cli.ConfigureOpenVPN(openVPNConfig)
	if err != nil {
		return nil, err
	}
	host.IsOnline = &pingable
	host.IsValidToken = &isValidToken
	if err := d.DB.Where("uuid = ?", host.UUID).Updates(&host).Error; err != nil {
		return nil, err
	}
	return &interfaces.Message{Message: "OpenVPN is configured!"}, nil
}

func (d *GormDatabase) ResolveHost(uuid string, name string) (*model.Host, error) {
	if uuid == "" && name == "" {
		return nil, errors.New("host-uuid and host-name both can not be empty")
	}
	if uuid != "" {
		host, _ := d.GetHost(uuid)
		if host != nil {
			return host, nil
		}
	}
	if name != "" {
		host, _ := d.GetHostByName(name)
		if host != nil {
			return host, nil
		}
	}
	var hostNames []string
	var hostUUIDs []string
	var count int
	hosts, err := d.GetHosts(false)
	if err != nil {
		return nil, err
	}
	for _, h := range hosts {
		hostNames = append(hostNames, h.Name)
		hostUUIDs = append(hostUUIDs, h.UUID)
		count++
	}
	return nil, errors.New(fmt.Sprintf("no valid host was found: host count: %d, host names found: %v uuids: %v", count, hostNames, hostUUIDs))
}

func attachOpenVPN(hosts []*model.Host) {
	resetHostClient := func(host *model.Host) {
		host.VirtualIP = ""
		host.ReceivedBytes = 0
		host.SentBytes = 0
		host.ConnectedSince = ""
	}

	resetHostsClient := func(hosts []*model.Host) {
		for _, host := range hosts {
			resetHostClient(host)
		}
	}

	oCli, _ := cligetter.GetOpenVPNClient()
	if oCli != nil {
		clients, _ := oCli.GetClients()
		if clients != nil {
			for _, host := range hosts {
				if client, found := (*clients)[host.GlobalUUID]; found {
					host.VirtualIP = client.VirtualIP
					host.ReceivedBytes = client.ReceivedBytes
					host.SentBytes = client.SentBytes
					host.ConnectedSince = client.ConnectedSince
				} else {
					resetHostClient(host)
				}
			}
		} else {
			resetHostsClient(hosts)
		}
	} else {
		resetHostsClient(hosts)
	}
}
