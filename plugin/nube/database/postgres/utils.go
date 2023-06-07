package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/rubix-os/plugin/nube/database/postgres/pgmodel"
	"regexp"
	"strings"
)

func buildSelectQuery() string {
	query := "histories.value, histories.timestamp, points.uuid AS rubix_point_uuid, points.name AS rubix_point_name," +
		" points.device_uuid AS rubix_device_uuid, points.device_name AS rubix_device_name, points.network_uuid AS " +
		"rubix_network_uuid, points.network_name AS rubix_network_name, points.global_uuid, points.location_uuid, " +
		"points.location_name, points.group_uuid, points.group_name, points.host_uuid, points.host_name"
	return query
}

func buildFilterQuery(filter *string) (filterQuery string, err error) {
	if filter == nil {
		return "", nil
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
				return "", err
			}
			filterQuery += column
		} else {
			value := fmt.Sprintf("'%s'", f)
			filterQuery = strings.Replace(filterQuery, valueFormat, value, -1)
		}
		isColumn = !isColumn
	}
	return filterQuery, err
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

func convertHistoryDataToResponse(historyData []*pgmodel.HistoryData) []*pgmodel.HistoryDataResponse {
	historyResponses := make([]*pgmodel.HistoryDataResponse, 0)
	indexMap := make(map[string]int)

	for _, history := range historyData {
		key := fmt.Sprintf("%s-%s-%s-%s-%s-%s", history.RubixNetworkUUID, history.RubixNetworkName,
			history.RubixDeviceUUID, history.RubixDeviceName, history.RubixPointUUID, history.RubixPointName)

		if index, ok := indexMap[key]; ok {
			historyResponses[index].Histories = append(historyResponses[index].Histories, &pgmodel.HistoryResponse{
				Value:     history.Value,
				Timestamp: history.Timestamp,
			})
		} else {
			historyResponse := &pgmodel.HistoryDataResponse{
				RubixNetworkUUID: history.RubixNetworkUUID,
				RubixNetworkName: history.RubixNetworkName,
				RubixDeviceUUID:  history.RubixDeviceUUID,
				RubixDeviceName:  history.RubixDeviceName,
				RubixPointUUID:   history.RubixPointUUID,
				RubixPointName:   history.RubixPointName,
				Host: &pgmodel.HostData{
					GlobalUUID:   history.GlobalUUID,
					HostUUID:     history.HostUUID,
					HostName:     history.HostName,
					GroupUUID:    history.GroupUUID,
					GroupName:    history.GroupName,
					LocationUUID: history.LocationUUID,
					LocationName: history.LocationName,
				},
				Histories: []*pgmodel.HistoryResponse{{
					Value:     history.Value,
					Timestamp: history.Timestamp,
				}},
			}

			historyResponses = append(historyResponses, historyResponse)
			indexMap[key] = len(historyResponses) - 1
		}
	}

	return historyResponses
}
