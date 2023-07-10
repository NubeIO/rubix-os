package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/utils/nuuid"
	"github.com/NubeIO/rubix-os/utils/structs"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"unicode"
)

const nameExcludeChars = "#+/!@%&*()\\}{[]}:;'\",.?"

func truncateString(str string, num int) string {
	ret := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		ret = str[0:num] + ""
	}
	return ret
}

func typeIsNil(t string, use string) string {
	if t == "" {
		return use
	}
	return t
}

func pluginIsNil(name string) string {
	if name == "" {
		return "system"
	}
	return name
}

func nameIsNil(name string) string {
	if name == "" {
		uuid := nuuid.MakeTopicUUID("")
		return fmt.Sprintf("n_%s", truncateString(uuid, 8))
	}
	return name
}

func checkTransport(t string) (string, error) {
	if t == "" {
		return model.TransType.IP, nil
	}
	i := structs.ArrayValues(model.TransType)
	if !structs.ArrayContains(i, t) {
		return "", errors.New("please provide a valid transport type ie: ip or serial")
	}
	return t, nil
}

func checkObjectType(t string) (model.ObjectType, error) {
	if t == "" {
		return model.ObjTypeAnalogValue, nil
	}
	objType := model.ObjectType(t)
	if _, ok := model.ObjectTypesMap[objType]; !ok {
		return "", errors.New("please provide a valid object type ie: analogInput or readCoil")
	}
	return objType, nil
}

func checkHistoryType(t string) (model.HistoryType, error) {
	if t == "" {
		return model.HistoryTypeInterval, nil
	}
	historyType := model.HistoryType(t)
	if _, ok := model.HistoryTypeMap[historyType]; !ok {
		return "", errors.New("please provide a valid history type ie: COV, INTERVAL or COV_AND_INTERVAL")
	}
	return historyType, nil
}

func checkHistoryCovType(t string) bool {
	if t == "" {
		return false
	}
	historyType := model.HistoryType(t)
	if _, ok := model.HistoryTypeCovMap[historyType]; !ok {
		return false
	}
	return true
}

func checkMemberState(t string) (model.MemberState, error) {
	if t == "" {
		return model.UnVerified, nil
	}
	memberState := model.MemberState(t)
	if _, ok := model.MemberStateMap[memberState]; !ok {
		return "", errors.New("please provide a valid member state ie: VERIFIED or UNVERIFIED")
	}
	return memberState, nil
}

func checkMemberPermission(t string) (model.MemberPermission, error) {
	if t == "" {
		return model.Read, nil
	}
	memberPermission := model.MemberPermission(t)
	if _, ok := model.MemberPermissionMap[memberPermission]; !ok {
		return "", errors.New("please provide a valid member permission ie: READ or WRITE")
	}
	return memberPermission, nil
}

func checkMemberDevicePlatform(t string) (model.MemberDevicePlatform, error) {
	if t == "" {
		return model.Android, nil
	}
	memberDevicePlatform := model.MemberDevicePlatform(t)
	if _, ok := model.MemberDevicePlatformMap[memberDevicePlatform]; !ok {
		return "", errors.New("please provide a valid device platform ie: ANDROID or IOS")
	}
	return memberDevicePlatform, nil
}

func (d *GormDatabase) deleteResponseBuilder(query *gorm.DB) (bool, error) {
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, gorm.ErrRecordNotFound
	} else {
		return true, nil
	}
}

func (d *GormDatabase) deleteResponse(query *gorm.DB) (*interfaces.Message, error) {
	msg := &interfaces.Message{
		Message: fmt.Sprintf("no record found, deleted count: %d", 0),
	}
	if query.Error != nil {
		return msg, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return msg, query.Error
	}
	msg.Message = fmt.Sprintf("deleted count: %d", query.RowsAffected)
	return msg, nil
}

func metaTagsArgsToKeyValues(metaTags string) [][]interface{} {
	mapMetaTags := map[string]string{}
	var keyValues [][]interface{}
	_ = json.Unmarshal([]byte(metaTags), &mapMetaTags)
	for k, v := range mapMetaTags {
		keyValues = append(keyValues, []interface{}{k, v})
	}
	if keyValues == nil {
		keyValues = append(keyValues, []interface{}{"", ""})
	}
	return keyValues
}

func validateName(name string) (string, error) {
	if name == "" {
		return "", errors.New(fmt.Sprintf("name cannot be empty"))
	}
	for _, char := range nameExcludeChars {
		if strings.Contains(name, string(char)) {
			return "", errors.New(fmt.Sprintf("name cannot contains: %c", char))
		}
	}
	for _, char := range name {
		if unicode.IsSymbol(char) {
			return "", errors.New(fmt.Sprintf("name cannot contains: %c", char))
		}
	}
	name = strings.TrimSpace(strings.Join(strings.Fields(name), " "))
	return name, nil
}

func marshalJson(jsonData datatypes.JSON) []byte {
	if jsonData == nil {
		jsonData = datatypes.JSON{}
	}
	mJsonData, _ := json.Marshal(jsonData)
	return mJsonData
}

func moduleNotFoundError(moduleName string) error {
	return errors.New(fmt.Sprintf("module with module name %s doesn't exist", moduleName))
}

func filterOutItem(slice []*string, item *string) []*string {
	filteredSlice := make([]*string, 0, len(slice))
	for _, s := range slice {
		if !reflect.DeepEqual(s, item) {
			filteredSlice = append(filteredSlice, s)
		}
	}
	return filteredSlice
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func (d *GormDatabase) getGlobalUUID() (string, error) {
	deviceInfo, err := d.RubixRegistry.GetDeviceInfo()
	if err != nil {
		return "", err
	}
	return deviceInfo.GlobalUUID, nil
}

func checkAlertStatus(s string) error {
	switch model.AlertStatus(s) {
	case model.AlertStatusActive:
		return nil
	case model.AlertStatusAcknowledged:
		return nil
	case model.AlertStatusClosed:
		return nil
	}
	return errors.New("invalid alert status, try active, acknowledged, closed")
}

func checkAlertSeverity(s string) error {
	switch model.AlertSeverity(s) {
	case model.AlertSeverityCrucial:
		return nil
	case model.AlertSeverityMinor:
		return nil
	case model.AlertSeverityInfo:
		return nil
	case model.AlertSeverityWarning:
		return nil
	}
	return errors.New("invalid alert status, try crucial, info, warning")
}

func checkAlertStatusClosed(s string) bool {
	return model.AlertStatus(s) == model.AlertStatusClosed
}

func checkAlertEntityType(s string) error {
	switch model.AlertEntityType(s) {
	case model.AlertEntityTypeGateway:
		return nil
	case model.AlertEntityTypeNetwork:
		return nil
	case model.AlertEntityTypeDevice:
		return nil
	case model.AlertEntityTypePoint:
		return nil
	case model.AlertEntityTypeService:
		return nil
	}
	return errors.New("invalid alert entity type, try gateway, network")
}

func checkAlertType(s string) error {
	switch model.AlertType(s) {
	case model.AlertTypePing:
		return nil
	case model.AlertTypeFault:
		return nil
	case model.AlertTypeThreshold:
		return nil
	case model.AlertTypeFlatLine:
		return nil
	}
	return errors.New("invalid alert type, try ping, threshold, fault")
}

func alertTypeTitle(s string) string {
	switch model.AlertType(s) {
	case model.AlertTypePing:
		return "Failed to ping the device"
	case model.AlertTypeFault:
		return "Fault"
	case model.AlertTypeThreshold:
		return "Out of range threshold"
	case model.AlertTypeFlatLine:
		return "Flat line"
	}
	return s
}

func checkAlertTarget(s string) error {
	switch model.AlertTarget(s) {
	case model.AlertTargetMobile:
		return nil
	case model.AlertTargetNone:
		return nil
	}
	return errors.New("invalid alert target, try mobile, none")
}
