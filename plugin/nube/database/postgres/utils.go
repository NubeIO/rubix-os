package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func buildSelectQuery(hasFnc bool) string {
	query := "histories.value, histories.timestamp, points.uuid AS rubix_point_uuid, points.name AS rubix_point_name," +
		" devices.uuid AS rubix_device_uuid, devices.name AS rubix_device_name, networks.uuid AS rubix_network_uuid," +
		" networks.name AS rubix_network_name"
	if hasFnc {
		fncQuery := ", flow_network_clones.global_uuid, flow_network_clones.client_id, " +
			"flow_network_clones.client_name, flow_network_clones.site_id, flow_network_clones.site_name, " +
			"flow_network_clones.device_id, flow_network_clones.device_name"
		query += fmt.Sprintf("%s", fncQuery)
	}
	return query
}

func buildFilterQuery(filter *string) (filterQuery string, hasFnc bool, err error) {
	if filter == nil {
		return "", hasFnc, nil
	}
	isColumn := true
	re := regexp.MustCompile(filterRegex)
	filters := re.FindAllString(*filter, -1)
	for _, f := range filters {
		if contains(f, orderOperators) {
			filterQuery += f
			continue
		}
		if contains(f, comparisonOperators) {
			filterQuery = strings.Replace(filterQuery, operatorFormat, f, -1)
			continue
		}
		if contains(f, logicalOperators) {
			if f == "&&" {
				filterQuery += " AND "
			} else {
				filterQuery += " OR "
			}
			continue
		}
		if isColumn {
			column, err := getColumn(f)
			if err != nil {
				return "", hasFnc, err
			}
			hasFnc = hasFnc || contains(f, flowNetworkCloneColumns)
			filterQuery += column
		} else {
			value := fmt.Sprintf("'%s'", f)
			filterQuery = strings.Replace(filterQuery, valueFormat, value, -1)
		}
		isColumn = !isColumn
	}
	return filterQuery, hasFnc, err
}

func getColumn(key string) (string, error) {
	if filterQueryMap[key] == "" {
		return "", errors.New("invalid column")
	}
	return filterQueryMap[key], nil
}

func contains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}

func convertData(data interface{}, v interface{}) error {
	mData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(mData, &v); err != nil {
		return err
	}
	return nil
}
