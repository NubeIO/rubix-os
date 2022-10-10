package decoder

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
)

const ZHTPlLenStaticV1 = 97
const ZHTPlLenStaticV2 = 102 // 9500ms

const ZHTPlLenWriteV1 = 51
const ZHTPlLenWriteV2 = 66 // 7200ms

const ZHTPlLenPollV1 = 40
const ZHTPlLenPollV2 = 47 // 6200ms

type TZipHydrotapBase struct {
	CommonValues
}

type TZipHydrotapStatic struct {
	LoRaFirmwareMajor       uint8  `json:"lora_firmware_major"`
	LoRaFirmwareMinor       uint8  `json:"lora_firmware_minor"`
	LoRaBuildMajor          uint8  `json:"lora_build_major"`
	LoRaBuildMinor          uint8  `json:"lora_build_minor"`
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
	FilterLogDateUV         string `json:"filter_log_date_uv"`
	FilterLogLitresUV       int    `json:"filter_log_litres_uv"`
}

const ZipHTTimerLength = 7

type TZipHydrotapTimer struct {
	TimeStart   int  `json:"time_start"`
	TimeStop    int  `json:"time_stop"`
	EnableStart bool `json:"enable_start"`
	EnableStop  bool `json:"enable_stop"`
}

type TZipHydrotapWrite struct {
	Time                         int                                 `json:"time"`
	DispenseTimeBoiling          int                                 `json:"dispense_time_boiling"`
	DispenseTimeChilled          int                                 `json:"dispense_time_chilled"`
	DispenseTimeSparkling        int                                 `json:"dispense_time_sparkling"`
	TemperatureSPBoiling         float32                             `json:"temperature_sp_boiling"`
	TemperatureSPChilled         float32                             `json:"temperature_sp_chilled"`
	TemperatureSPSparkling       float32                             `json:"temperature_sp_sparkling"`
	SleepModeSetting             int                                 `json:"sleep_mode_setting"`
	FilterInfoLifeLitresInternal int                                 `json:"filter_info_life_litres_internal"`
	FilterInfoLifeMonthsInternal int                                 `json:"filter_info_life_months_internal"`
	FilterInfoLifeLitresExternal int                                 `json:"filter_info_life_litres_external"`
	FilterInfoLifeMonthsExternal int                                 `json:"filter_info_life_months_external"`
	SafetyAllowTapChanges        bool                                `json:"safety_allow_tap_changes"`
	SafetyLock                   bool                                `json:"safety_lock"`
	SafetyHotIsolation           bool                                `json:"safety_hot_isolation"`
	SecurityEnable               bool                                `json:"security_enable"`
	SecurityPin                  int                                 `json:"security_pin"`
	Timers                       [ZipHTTimerLength]TZipHydrotapTimer `json:"timers"`
	// Packet V2
	FilterInfoLifeLitresUV int `json:"filter_info_life_litres_uv"`
	FilterInfoLifeMonthsUV int `json:"filter_info_life_months_uv"`
	CO2LifeGrams           int `json:"co2_life_grams"`
	CO2LifeMonths          int `json:"co2_life_months"`
	CO2Pressure            int `json:"co2_pressure"`
	CO2TankCapacity        int `json:"co2_tank_capacity"`
	CO2AbsorptionRate      int `json:"co2_absorption_rate"`
	SparklingFlowRate      int `json:"sparkling_flow_rate"`
	SparklingFlushTime     int `json:"sparkling_flush_time"`
}

type TZipHydrotapPoll struct {
	Rebooted bool `json:"rebooted"`
	// StaticCOVFlag                     bool    `json:"static_cov_flag"`
	// WriteCOVFlag                      bool    `json:"write_cov_flag"`
	SleepModeStatus                   int8    `json:"sleep_mode_status"`
	TemperatureNTCBoiling             float32 `json:"temperature_ntc_boiling"`
	TemperatureNTCChilled             float32 `json:"temperature_ntc_chilled"`
	TemperatureNTCStream              float32 `json:"temperature_ntc_stream"`
	TemperatureNTCCondensor           float32 `json:"temperature_ntc_condensor"`
	UsageEnergyKWh                    float32 `json:"usage_energy_kwh"`
	UsageWaterDeltaDispensesBoiling   int     `json:"usage_water_delta_dispenses_boiling"`
	UsageWaterDeltaDispensesChilled   int     `json:"usage_water_delta_dispenses_chilled"`
	UsageWaterDeltaDispensesSparkling int     `json:"usage_water_delta_dispenses_sparkling"`
	UsageWaterDeltaLitresBoiling      float32 `json:"usage_water_delta_litres_boiling"`
	UsageWaterDeltaLitresChilled      float32 `json:"usage_water_delta_litres_chilled"`
	UsageWaterDeltaLitresSparkling    float32 `json:"usage_water_delta_litres_sparkling"`
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
	// Packet V2
	FilterInfoUsageLitresUV int  `json:"filter_info_usage_litres_uv"`
	FilterInfoUsageDaysUV   int  `json:"filter_info_usage_days_uv"`
	FilterWarningUV         bool `json:"filter_warning_uv"`
	CO2LowGasWarning        bool `json:"co2_low_gas_warning"`
	CO2UsageGrams           int  `json:"co2_usage_grams"`
	CO2UsageDays            int  `json:"co2_usage_days"`
}

type TZipHydrotapWriteOnly struct {
	Reboot            int `json:"reboot"`
	ResetFilter       int `json:"reset_filter"`
	RemoteCalibration int `json:"remote_calibration"`
	ResetEnergy       int `json:"reset_energy"`
}

type TZipHydrotapStaticFull struct {
	TZipHydrotapBase
	TZipHydrotapStatic
}

type TZipHydrotapWriteFull struct {
	TZipHydrotapBase
	TZipHydrotapWrite
}

type TZipHydrotapPollFull struct {
	TZipHydrotapBase
	TZipHydrotapPoll
}

type TZipHydrotapAll struct {
	TZipHydrotapBase
	TZipHydrotapWriteOnly
	TZipHydrotapWrite
	TZipHydrotapPoll
}

type TZHTPayloadType int

const (
	ErrorData = iota
	StaticData
	WriteData
	PollData
)

func DecodeZHT(data string, _ *LoRaDeviceDescription) (*CommonValues, interface{}) {
	bytes := getPayloadBytes(data)
	switch pl := getPayloadType(data); pl {
	// TODO: This should be meta data when it gets supported
	// case StaticData:
	//     payload := staticPayloadDecoder(bytes)
	//     payloadFull := TZipHydrotapStaticFull{TZipHydrotapStatic: payload}
	//     return &payloadFull.CommonValues, payloadFull
	case WriteData:
		payload := writePayloadDecoder(bytes)
		payloadFull := TZipHydrotapWriteFull{TZipHydrotapWrite: payload}
		return &payloadFull.CommonValues, payloadFull
	case PollData:
		payload := pollPayloadDecoder(bytes)
		payloadFull := TZipHydrotapPollFull{TZipHydrotapPoll: payload}
		return &payloadFull.CommonValues, payloadFull
	}

	return nil, nil
}

func getPayloadType(data string) TZHTPayloadType {
	plID, _ := strconv.ParseInt(data[14:16], 16, 0)
	return TZHTPayloadType(plID)
}

func CheckPayloadLengthZHT(data string) bool {
	payloadLength := len(data) - 10 // removed addr, nonce and MAC
	payloadLength /= 2
	payloadType := getPayloadType(data)
	dataLength, _ := strconv.ParseInt(data[12:14], 16, 0)

	if getPacketVersion(data) == 1 {
		return (payloadType == StaticData && dataLength == ZHTPlLenStaticV1 && payloadLength > ZHTPlLenStaticV1) ||
			(payloadType == WriteData && dataLength == ZHTPlLenWriteV1 && payloadLength > ZHTPlLenWriteV1) ||
			(payloadType == PollData && dataLength == ZHTPlLenPollV1 && payloadLength > ZHTPlLenPollV1)
	} else if getPacketVersion(data) == 2 {
		return (payloadType == StaticData && dataLength == ZHTPlLenStaticV2 && payloadLength > ZHTPlLenStaticV2) ||
			(payloadType == WriteData && dataLength == ZHTPlLenWriteV2 && payloadLength > ZHTPlLenWriteV2) ||
			(payloadType == PollData && dataLength == ZHTPlLenPollV2 && payloadLength > ZHTPlLenPollV2)
	}
	return false
}

func GetPointsStructZHT() interface{} {
	return TZipHydrotapAll{}
}

func getPayloadBytes(data string) []byte {
	length, _ := strconv.ParseInt(data[12:14], 16, 0)
	bytes, _ := hex.DecodeString(data[16 : 16+((length-1)*2)])
	return bytes
}

func getPacketVersion(data string) uint8 {
	v, _ := strconv.ParseInt(data[16:18], 16, 0)
	return uint8(v)
}

func bytesToString(bytes []byte) string {
	str := ""
	for _, b := range bytes {
		if b == 0 {
			break
		}
		str += string(b)
	}
	return str
}

func bytesToDate(bytes []byte) string {
	return fmt.Sprintf("%d/%d/%d", bytes[0], bytes[1], bytes[2])
}

func staticPayloadDecoder(data []byte) TZipHydrotapStatic {
	index := 1
	fwMa := data[index]
	index += 1
	fwMi := data[index]
	index += 1
	buildMa := data[index]
	index += 1
	buildMi := data[index]
	index += 1
	sn := bytesToString(data[index : index+15])
	index += 15
	mn := bytesToString(data[index : index+20])
	index += 20
	pn := bytesToString(data[index : index+20])
	index += 20
	fw := bytesToString(data[index : index+20])
	index += 20
	calDate := bytesToDate(data[index : index+3])
	index += 3
	f50lDate := bytesToDate(data[index : index+3])
	index += 3
	filtLogDateInt := bytesToDate(data[index : index+3])
	index += 3
	filtLogLitresInt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	filtLogDateExt := bytesToDate(data[index : index+3])
	index += 3
	filtLogLitresExt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2

	filtLogDateUV := ""
	filtLogLitresUV := 0
	if data[0] >= 2 {
		filtLogDateUV = bytesToDate(data[index : index+3])
		index += 3
		filtLogLitresUV = int(binary.LittleEndian.Uint16(data[index : index+2]))
		index += 2
	}

	return TZipHydrotapStatic{
		LoRaFirmwareMajor:       fwMa,
		LoRaFirmwareMinor:       fwMi,
		LoRaBuildMajor:          buildMa,
		LoRaBuildMinor:          buildMi,
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
		FilterLogDateUV:         filtLogDateUV,
		FilterLogLitresUV:       filtLogLitresUV,
	}
}

func writePayloadDecoder(data []byte) TZipHydrotapWrite {
	index := 1
	time := int(binary.LittleEndian.Uint32(data[index : index+4]))
	index += 4
	dispB := int(data[index])
	index += 1
	dispC := int(data[index])
	index += 1
	dispS := int(data[index])
	index += 1
	tempSpB := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	tempSpC := float32(int(data[index]))
	index += 1
	tempSpS := float32(int(data[index]))
	index += 1
	sm := int(data[index])
	index += 1
	filLyfLtrInt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	filLyfMnthInt := int(data[index])
	index += 1
	filLyfLtrExt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	filLyfMnthExt := int(data[index])
	index += 1
	sfTap := (data[index]>>2)&1 == 1
	sfL := (data[index]>>1)&1 == 1
	sfHi := (data[index]>>0)&1 == 1
	index += 1
	secUI16 := binary.LittleEndian.Uint16(data[index : index+2])
	secEn := secUI16 >= 10000
	secPin := int(secUI16 % 10000)
	index += 2

	var timers [ZipHTTimerLength]TZipHydrotapTimer
	var u16 uint16
	for i := 0; i < ZipHTTimerLength; i++ {
		u16 = binary.LittleEndian.Uint16(data[index : index+2])
		timers[i].TimeStart = int(u16 % 10000)
		timers[i].EnableStart = u16 >= 10000
		index += 2
		u16 = binary.LittleEndian.Uint16(data[index : index+2])
		timers[i].TimeStop = int(u16 % 10000)
		timers[i].EnableStop = u16 >= 10000
		index += 2
	}

	filLyfLtrUV := 0
	filLyfMnthUV := 0
	cO2LyfGrams := 0
	cO2LyfMnths := 0
	cO2Pressure := 0
	cO2TankCap := 0
	cO2AbsorpRate := 0
	sparklFlowRate := 0
	sparklFlushTime := 0
	if data[0] >= 2 {
		filLyfLtrUV = int(binary.LittleEndian.Uint16(data[index : index+2]))
		index += 2
		filLyfMnthUV = int(data[index])
		index += 1
		cO2LyfGrams = int(binary.LittleEndian.Uint16(data[index : index+2]))
		index += 2
		cO2LyfMnths = int(data[index])
		index += 1
		cO2Pressure = int(data[index])
		index += 1
		cO2TankCap = int(binary.LittleEndian.Uint16(data[index : index+2]))
		index += 2
		cO2AbsorpRate = int(binary.LittleEndian.Uint16(data[index : index+2]))
		index += 2
		sparklFlowRate = int(binary.LittleEndian.Uint16(data[index : index+2]))
		index += 2
		sparklFlushTime = int(binary.LittleEndian.Uint16(data[index : index+2]))
		index += 2
	}

	return TZipHydrotapWrite{
		Time:                         time,
		DispenseTimeBoiling:          dispB,
		DispenseTimeChilled:          dispC,
		DispenseTimeSparkling:        dispS,
		TemperatureSPBoiling:         tempSpB,
		TemperatureSPChilled:         tempSpC,
		TemperatureSPSparkling:       tempSpS,
		SleepModeSetting:             sm,
		FilterInfoLifeLitresInternal: filLyfLtrInt,
		FilterInfoLifeMonthsInternal: filLyfMnthInt,
		FilterInfoLifeLitresExternal: filLyfLtrExt,
		FilterInfoLifeMonthsExternal: filLyfMnthExt,
		SafetyAllowTapChanges:        sfTap,
		SafetyLock:                   sfL,
		SafetyHotIsolation:           sfHi,
		SecurityEnable:               secEn,
		SecurityPin:                  secPin,
		Timers:                       timers,
		// Pkt V2
		FilterInfoLifeLitresUV: filLyfLtrUV,
		FilterInfoLifeMonthsUV: filLyfMnthUV,
		CO2LifeGrams:           cO2LyfGrams,
		CO2LifeMonths:          cO2LyfMnths,
		CO2Pressure:            cO2Pressure,
		CO2TankCapacity:        cO2TankCap,
		CO2AbsorptionRate:      cO2AbsorpRate,
		SparklingFlowRate:      sparklFlowRate,
		SparklingFlushTime:     sparklFlushTime,
	}
}

func pollPayloadDecoder(data []byte) TZipHydrotapPoll {
	index := 1
	rebooted := (data[index]>>5)&1 == 1
	// sCov := (data[index]>>6)&1 == 1
	// wCov := (data[index]>>7)&1 == 1
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
	kwh := float32(binary.LittleEndian.Uint32(data[index:index+4])) * 0.1
	index += 4
	dltDispB := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	dltDispC := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	dltDispS := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	dltLtrB := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	dltLtrC := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	dltLtrS := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	warningIndex := index
	fltrWrnInt := (data[index]>>0)&1 == 1
	fltrWrnExt := (data[index]>>1)&1 == 1
	index += 1
	fltrNfoUseLtrInt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	fltrNfoUseDayInt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	fltrNfoUseLtrExt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	fltrNfoUseDayExt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2

	fltrNfoUseLtrUV := 0
	fltrNfoUseDayUV := 0
	fltrWrnUV := false
	cO2GasPressureWrn := false
	cO2UsgGrams := 0
	cO2UsgDays := 0
	if data[0] >= 2 {
		fltrNfoUseLtrUV = int(binary.LittleEndian.Uint16(data[index : index+2]))
		index += 2
		fltrNfoUseDayUV = int(binary.LittleEndian.Uint16(data[index : index+2]))
		index += 2

		fltrWrnUV = (data[warningIndex]>>2)&1 == 1
		cO2GasPressureWrn = (data[warningIndex]>>3)&1 == 1
		cO2UsgGrams = int(binary.LittleEndian.Uint16(data[index : index+2]))
		index += 2
		cO2UsgDays = int(data[index])
		index += 1
	}

	return TZipHydrotapPoll{
		Rebooted:                          rebooted,
		SleepModeStatus:                   sms,
		TemperatureNTCBoiling:             tempB,
		TemperatureNTCChilled:             tempC,
		TemperatureNTCStream:              tempS,
		TemperatureNTCCondensor:           tempCond,
		UsageEnergyKWh:                    kwh,
		UsageWaterDeltaDispensesBoiling:   dltDispB,
		UsageWaterDeltaDispensesChilled:   dltDispC,
		UsageWaterDeltaDispensesSparkling: dltDispS,
		UsageWaterDeltaLitresBoiling:      dltLtrB,
		UsageWaterDeltaLitresChilled:      dltLtrC,
		UsageWaterDeltaLitresSparkling:    dltLtrS,
		Fault1:                            f1,
		Fault2:                            f2,
		Fault3:                            f3,
		Fault4:                            f4,
		FilterWarningInternal:             fltrWrnInt,
		FilterWarningExternal:             fltrWrnExt,
		FilterInfoUsageLitresInternal:     fltrNfoUseLtrInt,
		FilterInfoUsageDaysInternal:       fltrNfoUseDayInt,
		FilterInfoUsageLitresExternal:     fltrNfoUseLtrExt,
		FilterInfoUsageDaysExternal:       fltrNfoUseDayExt,
		// Pkt V2
		FilterInfoUsageLitresUV: fltrNfoUseLtrUV,
		FilterInfoUsageDaysUV:   fltrNfoUseDayUV,
		FilterWarningUV:         fltrWrnUV,
		CO2LowGasWarning:        cO2GasPressureWrn,
		CO2UsageGrams:           cO2UsgGrams,
		CO2UsageDays:            cO2UsgDays,
	}
}
