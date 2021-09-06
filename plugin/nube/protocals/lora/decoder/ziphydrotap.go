package decoder

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strconv"
)

type TZipHydrotapBase struct {
	CommonValues
	PayloadType string `json:"payload_type"`
}

type TZipHydrotapStatic struct {
	TZipHydrotapBase
	SerialNumber            string `json:"serial_number"`
	ModelNumber             string `json:"model_number"`
	ProductNumber           string `json:"product_number"`
	FirmwareVersion         string `json:"firmware_version"`
	CalibrationDate         string `json:"calibration_date"`
	First50LitresData       string `json:"first_50_litres_data"`
	FilterLogDateInternal   string `json:"filter_log_date_internal"`
	FilterLogLitresInternal int    `json:"filter_log_litres_internal"`
	FilterLogDateExternal   string `json:"filter_log_date_external"`
	FilterLogLitresExternal int    `json:"filter_log_litres_external"`
}

type TZipHydrotapWrite struct {
	TZipHydrotapBase
	Time                   string  `json:"time"`
	DispenseTimeBoiling    int     `json:"dispense_time_boiling"`
	DispenseTimeChilled    int     `json:"dispense_time_chilled"`
	DispenseTimeSparkling  int     `json:"dispense_time_sparkling"`
	TemperatureSPBoiling   float32 `json:"temperature_sp_boiling"`
	TemperatureSPChilled   float32 `json:"temperature_sp_chilled"`
	TemperatureSPSparkling float32 `json:"temperature_sp_sparkling"`
	// TODO: timers
	SleepModeSetting              int8   `json:"sleep_mode_setting"`
	FilterInfoWriteLitresInternal int    `json:"filter_info_write_litres_internal"`
	FilterInfoWriteMonthsInternal int    `json:"filter_info_write_months_internal"`
	FilterInfoWriteLitresExternal int    `json:"filter_info_write_litres_external"`
	FilterInfoWriteMonthsExternal int    `json:"filter_info_write_months_external"`
	SafetyAllowTapChanges         bool   `json:"safety_allow_tap_changes"`
	SafetyLock                    bool   `json:"safety_lock"`
	SafetyHotIsolation            bool   `json:"safety_hot_isolation"`
	SecurityEnable                bool   `json:"security_enable"`
	SecurityPin                   string `json:"security_pin"`
}

type TZipHydrotapPoll struct {
	TZipHydrotapBase
	StaticCOVFlag                     bool    `json:"static_cov_flag"`
	WriteCOVFlag                      bool    `json:"write_cov_flag"`
	SleepModeStatus                   int8    `json:"sleep_mode_status"`
	TemperatureNTCBoiling             float32 `json:"temperature_ntc_boiling"`
	TemperatureNTCChilled             float32 `json:"temperature_ntc_chilled"`
	TemperatureNTCStream              float32 `json:"temperature_ntc_stream"`
	TemperatureNTCCondensor           float32 `json:"temperature_ntc_condensor"`
	UsageEnergyKWh                    float32 `json:"usage_energy_kwh"`
	UsageWaterDeltaDispensesBoiling   int     `json:"usage_water_delta_dispenses_boiling"`
	UsageWaterDeltaDispensesChilled   int     `json:"usage_water_delta_dispenses_chilled"`
	UsageWaterDeltaDispensesSparkling int     `json:"usage_water_delta_dispenses_sparkling"`
	UsageWaterDeltaLitresBoiling      int     `json:"usage_water_delta_litres_boiling"`
	UsageWaterDeltaLitresChilled      int     `json:"usage_water_delta_litres_chilled"`
	UsageWaterDeltaLitresSparkling    int     `json:"usage_water_delta_litres_sparkling"`
	Fault1                            uint8   `json:"fault_1"`
	Fault2                            uint8   `json:"fault_2"`
	Fault3                            uint8   `json:"fault_3"`
	Fault4                            uint8   `json:"fault_4"`
	FilterWarningInternal             bool    `json:"filter_warning_internal"`
	FilterWarningExternal             bool    `json:"filter_warning_external"`
	FilterInfoUsageLitresInternal     int     `json:"filter_info_usage_litres_internal"`
	FilterInfoUsageDaysInternal       int     `json:"filter_info_usage_days_internal"`
	FilterInfoUsageLitresExternal     int     `json:"filter_info_usage_litres_external"`
	FilterInfoUsageDaysExternal       int     `json:"filter_info_usage_days_external"`
}

type TZHTPayloadType int

const (
	ErrorData = iota
	StaticData
	WriteData
	PollData
)

const ZHT_HEX_STR_DATA_START = 14

func ZipHydrotap(data string, sensor TSensorType) (TZipHydrotapBase, interface{}) {
	bytes := ZHtGetPayloadBytes(data)
	log.Printf("ZHT BYTES: %d", len(bytes))
	switch pl := ZHtGetPayloadType(data); pl {
	case StaticData:
		payloadFull := ZHtStaticPayloadDecoder(bytes)
		commonData := Common(data, sensor)
		payloadFull.CommonValues = commonData
		payloadFull.PayloadType = "static"
		return payloadFull.TZipHydrotapBase, payloadFull
	case WriteData:
		payloadFull := ZHtWritePayloadDecoder(bytes)
		commonData := Common(data, sensor)
		payloadFull.CommonValues = commonData
		payloadFull.PayloadType = "write"
		return payloadFull.TZipHydrotapBase, payloadFull
	case PollData:
		payloadFull := ZHtPollPayloadDecoder(bytes)
		commonData := Common(data, sensor)
		payloadFull.CommonValues = commonData
		payloadFull.PayloadType = "poll"
		return payloadFull.TZipHydrotapBase, payloadFull
	}

	return TZipHydrotapBase{}, nil
}

func ZHtGetPayloadType(data string) TZHTPayloadType {
	plId, _ := strconv.ParseInt(data[14:16], 16, 0)
	return TZHTPayloadType(plId)
}

func ZHtCheckPayloadLength(data string) bool {
	payloadLength := len(data) - 10 // removed addr, nonce and MAC
	payloadLength /= 2
	payloadType := ZHtGetPayloadType(data)
	dataLength, _ := strconv.ParseInt(data[12:14], 16, 0)
	log.Printf("ZHT data_length: %d\n", dataLength)

	return (payloadType == StaticData && dataLength == 93 && payloadLength > 93) ||
		(payloadType == WriteData && dataLength == 23 && payloadLength > 23) ||
		(payloadType == PollData && dataLength == 40 && payloadLength > 40)
}

func ZHtGetPayloadBytes(data string) []byte {
	length, _ := strconv.ParseInt(data[12:14], 16, 0)
	bytes, _ := hex.DecodeString(data[14 : 14+length*2])
	return bytes
}

func ZHtBytesToString(bytes []byte) string {
	str := ""
	for _, b := range bytes {
		if b == 0 {
			break
		}
		str += string(b)
	}
	return str
}

func ZHtBytesToDate(bytes []byte) string {
	return fmt.Sprintf("%d/%d/%d", bytes[0], bytes[1], bytes[2])
}

func ZHtStaticPayloadDecoder(data []byte) TZipHydrotapStatic {
	index := 1
	sn := ZHtBytesToString(data[index : index+15])
	index += 15
	mn := ZHtBytesToString(data[index : index+20])
	index += 20
	pn := ZHtBytesToString(data[index : index+20])
	index += 20
	fw := ZHtBytesToString(data[index : index+20])
	index += 20
	calDate := ZHtBytesToDate(data[index : index+3])
	index += 3
	f50lDate := ZHtBytesToDate(data[index : index+3])
	index += 3
	filtLogDateInt := ZHtBytesToDate(data[index : index+3])
	index += 3
	filtLogLitresInt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	filtLogDateExt := ZHtBytesToDate(data[index : index+3])
	index += 3
	filtLogLitresExt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	return TZipHydrotapStatic{
		SerialNumber:            sn,
		ModelNumber:             mn,
		ProductNumber:           pn,
		FirmwareVersion:         fw,
		CalibrationDate:         calDate,
		First50LitresData:       f50lDate,
		FilterLogDateInternal:   filtLogDateInt,
		FilterLogLitresInternal: filtLogLitresInt,
		FilterLogDateExternal:   filtLogDateExt,
		FilterLogLitresExternal: filtLogLitresExt,
	}
}

func ZHtWritePayloadDecoder(data []byte) TZipHydrotapWrite {
	index := 1
	time := fmt.Sprintf("%d", binary.LittleEndian.Uint32(data[index:index+4]))
	index += 4
	dispB := int(data[index])
	index += 1
	dispC := int(data[index])
	index += 1
	dispS := int(data[index])
	index += 1
	tempSpB := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	tempSpC := int(data[index])
	index += 1
	tempSpS := int(data[index])
	index += 1
	// TODO: timers
	sm := int8(data[index])
	index += 1
	filLyfLtrInt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	filLyfMnthInt := int(data[index])
	index += 1
	filLyfLtrExt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	filLyfMnthExt := int(data[index])
	index += 1
	sfTap := (data[index]>>2)&1 == 1
	sfL := (data[index]>>1)&1 == 1
	sfHi := (data[index]>>0)&1 == 1
	index += 1
	secUI16 := binary.LittleEndian.Uint16(data[index : index+2])
	secEn := (secUI16>>15)&1 == 1
	secPin := fmt.Sprintf("%.4d", secUI16&0x7FFF)
	return TZipHydrotapWrite{
		Time:                          time,
		DispenseTimeBoiling:           dispB,
		DispenseTimeChilled:           dispC,
		DispenseTimeSparkling:         dispS,
		TemperatureSPBoiling:          tempSpB,
		TemperatureSPChilled:          float32(tempSpC),
		TemperatureSPSparkling:        float32(tempSpS),
		SleepModeSetting:              sm,
		FilterInfoWriteLitresInternal: int(filLyfLtrInt),
		FilterInfoWriteMonthsInternal: filLyfMnthInt,
		FilterInfoWriteLitresExternal: int(filLyfLtrExt),
		FilterInfoWriteMonthsExternal: filLyfMnthExt,
		SafetyAllowTapChanges:         sfTap,
		SafetyLock:                    sfL,
		SafetyHotIsolation:            sfHi,
		SecurityEnable:                secEn,
		SecurityPin:                   secPin,
	}
}

func ZHtPollPayloadDecoder(data []byte) TZipHydrotapPoll {
	index := 1
	sCov := (data[index]>>6)&1 == 1
	wCov := (data[index]>>7)&1 == 1
	sms := int8((data[index]) & 0x3F)
	index += 1
	tempB := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	tempC := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	tempS := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	tempCond := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	f1 := data[index]
	index += 1
	f2 := data[index]
	index += 1
	f3 := data[index]
	index += 1
	f4 := data[index]
	index += 1
	kwh := math.Float32frombits(binary.LittleEndian.Uint32(data[index : index+4]))
	index += 4
	dltDispB := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	dltDispC := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	dltDispS := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	dltLtrB := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	dltLtrC := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	dltLtrS := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	fltrWrnInt := (data[index]>>1)&1 == 1
	fltrWrnExt := (data[index]>>0)&1 == 1
	index += 1
	fltrNfoUseLtrInt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	fltrNfoUseDayInt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	fltrNfoUseLtrExt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	fltrNfoUseDayExt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2

	return TZipHydrotapPoll{
		StaticCOVFlag:                     sCov,
		WriteCOVFlag:                      wCov,
		SleepModeStatus:                   sms,
		TemperatureNTCBoiling:             tempB,
		TemperatureNTCChilled:             tempC,
		TemperatureNTCStream:              tempS,
		TemperatureNTCCondensor:           tempCond,
		UsageEnergyKWh:                    kwh,
		UsageWaterDeltaDispensesBoiling:   int(dltDispB),
		UsageWaterDeltaDispensesChilled:   int(dltDispC),
		UsageWaterDeltaDispensesSparkling: int(dltDispS),
		UsageWaterDeltaLitresBoiling:      int(dltLtrB),
		UsageWaterDeltaLitresChilled:      int(dltLtrC),
		UsageWaterDeltaLitresSparkling:    int(dltLtrS),
		Fault1:                            f1,
		Fault2:                            f2,
		Fault3:                            f3,
		Fault4:                            f4,
		FilterWarningInternal:             fltrWrnInt,
		FilterWarningExternal:             fltrWrnExt,
		FilterInfoUsageLitresInternal:     int(fltrNfoUseLtrInt),
		FilterInfoUsageDaysInternal:       int(fltrNfoUseDayInt),
		FilterInfoUsageLitresExternal:     int(fltrNfoUseLtrExt),
		FilterInfoUsageDaysExternal:       int(fltrNfoUseDayExt),
	}
}
