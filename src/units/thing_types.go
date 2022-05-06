package units

type PointThingTypes struct {
	AHEnable                                       PointThing `json:"AH Enable"`
	AHUFanFault                                    PointThing `json:"AHU Fan Fault"`
	AHUFanSpeed                                    PointThing `json:"AHU Fan Speed"`
	AHUFanStartStop                                PointThing `json:"AHU Fan StartStop"`
	AHUFanStatus                                   PointThing `json:"AHU Fan Status"`
	AHUReturnAirCO2                                PointThing `json:"AHU Return Air CO2"`
	AHUReturnAirFanFault                           PointThing `json:"AHU Return Air Fan Fault"`
	AHUReturnAirFanSpeed                           PointThing `json:"AHU Return Air Fan Speed"`
	AHUReturnAirFanStartStop                       PointThing `json:"AHU Return Air Fan StartStop"`
	AHUReturnAirFanStatus                          PointThing `json:"AHU Return Air Fan Status"`
	AHUReturnAirHumidity                           PointThing `json:"AHU Return Air Humidity"`
	AHUReturnAirTemp                               PointThing `json:"AHU Return Air Temperature"`
	AHUReturnAirConstPressure                      PointThing `json:"AHU Return Air const Pressure"`
	AHUReturnAirConstPressureSetPoint              PointThing `json:"AHU Return Air const Pressure SetPoint"`
	AfterHoursActiveTime                           PointThing `json:"After Hours Active Time"`
	AfterHoursElapsedTime                          PointThing `json:"After Hours ElapsedTime"`
	AirOffCoilTemperature                          PointThing `json:"Air Off Coil Temperature"`
	Airflow                                        PointThing `json:"Airflow"`
	AirflowSetPoint                                PointThing `json:"Airflow SetPoint"`
	BoilerFault                                    PointThing `json:"Boiler Fault"`
	BoilerStartStop                                PointThing `json:"Boiler StartStop"`
	BoilerStatus                                   PointThing `json:"Boiler Status"`
	CarParkExhaustAirFanFault                      PointThing `json:"CarPark Exhaust Air Fan Fault"`
	CarParkExhaustAirFanSpeed                      PointThing `json:"CarPark Exhaust Air Fan Speed"`
	CarParkExhaustAirFanStartStop                  PointThing `json:"CarPark Exhaust Air Fan StartStop"`
	CarParkExhaustAirFanStatus                     PointThing `json:"CarPark Exhaust Air Fan Status"`
	CarParkExhaustAirConstPressure                 PointThing `json:"CarPark Exhaust Air const Pressure"`
	CarParkExhaustAirConstPressureSetPoint         PointThing `json:"CarPark Exhaust Air const Pressure SetPoint"`
	CarParkExhaustCarbonMonoxide                   PointThing `json:"CarPark Exhaust Carbon Monoxide"`
	CarParkExhaustCarbonMonoxideSetPoint           PointThing `json:"CarPark Exhaust Carbon Monoxide SetPoint"`
	CarParkSupplyAirFanFault                       PointThing `json:"CarPark Supply Air Fan Fault"`
	CarParkSupplyAirFanSpeed                       PointThing `json:"CarPark Supply Air Fan Speed"`
	CarParkSupplyAirFanStartStop                   PointThing `json:"CarPark Supply Air Fan StartStop"`
	CarParkSupplyAirFanStatus                      PointThing `json:"CarPark Supply Air Fan Status"`
	CarParkSupplyAirConstPressure                  PointThing `json:"CarPark Supply Air const Pressure"`
	CarParkSupplyAirConstPressureSetPoint          PointThing `json:"CarPark Supply Air const Pressure SetPoint"`
	ChilledWaterFlow                               PointThing `json:"Chilled Water Flow"`
	ChilledWaterFlowSetPoint                       PointThing `json:"Chilled Water Flow SetPoint"`
	ChilledWaterPumpDifferentialPressure           PointThing `json:"Chilled Water Pump Differential Pressure"`
	ChilledWaterPumpDifferentialPressureSetPoint   PointThing `json:"Chilled Water Pump Differential Pressure SetPoint"`
	ChilledWaterPumpFault                          PointThing `json:"Chilled Water Pump Fault"`
	ChilledWaterPumpSpeed                          PointThing `json:"Chilled Water Pump Speed"`
	ChilledWaterPumpStartStop                      PointThing `json:"Chilled Water Pump StartStop"`
	ChilledWaterPumpStatus                         PointThing `json:"Chilled Water Pump Status"`
	ChilledWaterReturnTemperature                  PointThing `json:"Chilled Water Return Temperature"`
	ChilledWaterSupplyTemperature                  PointThing `json:"Chilled Water Supply Temperature"`
	ChilledWaterSupplyTemperatureSetPoint          PointThing `json:"Chilled Water Supply Temperature SetPoint"`
	ChillerFault                                   PointThing `json:"Chiller Fault"`
	ChillerStartStop                               PointThing `json:"Chiller StartStop"`
	ChillerStatus                                  PointThing `json:"Chiller Status"`
	CommonBypassValve                              PointThing `json:"Common Bypass Valve"`
	CommonChilledWaterBypassValve                  PointThing `json:"Common Chilled Water Bypass Valve"`
	CommonChilledWaterDifferentialPressure         PointThing `json:"Common Chilled Water Differential Pressure"`
	CommonChilledWaterDifferentialPressureSetPoint PointThing `json:"Common Chilled Water Differential Pressure SetPoint"`
	CommonChilledWaterReturnTemperature            PointThing `json:"Common Chilled Water Return Temperature"`
	CommonChilledWaterSupplyTemperature            PointThing `json:"Common Chilled Water Supply Temperature"`
	CommonCondenserPressureDiff                    PointThing `json:"Common Condenser Pressure Diff"`
	CommonCondenserWaterReturnTemperature          PointThing `json:"Common Condenser Water Return Temperature"`
	CommonCondenserWaterSupplyTemperature          PointThing `json:"Common Condenser Water Supply Temperature"`
	CommonHotWaterReturnTemperature                PointThing `json:"Common Hot Water Return Temperature"`
	CommonHotWaterSupplyTemperature                PointThing `json:"Common Hot Water Supply Temperature"`
	CompressorStartStop                            PointThing `json:"Compressor StartStop"`
	CondenserWaterCall                             PointThing `json:"Condenser Water Call"`
	CondenserWaterFlow                             PointThing `json:"Condenser Water Flow"`
	CondenserWaterPumpDifferentialPressure         PointThing `json:"Condenser Water Pump Differential Pressure"`
	CondenserWaterPumpDifferentialPressureSetPoint PointThing `json:"Condenser Water Pump Differential Pressure SetPoint"`
	CondenserWaterPumpFault                        PointThing `json:"Condenser Water Pump Fault"`
	CondenserWaterPumpSpeed                        PointThing `json:"Condenser Water Pump Speed"`
	CondenserWaterPumpStartStop                    PointThing `json:"Condenser Water Pump StartStop"`
	CondenserWaterPumpStatus                       PointThing `json:"Condenser Water Pump Status"`
	CondenserWaterReturnTemperature                PointThing `json:"Condenser Water Return Temperature"`
	CondenserWaterSupplyTemperature                PointThing `json:"Condenser Water Supply Temperature"`
	CondenserWaterSupplyTemperatureSetPoint        PointThing `json:"Condenser Water Supply Temperature SetPoint"`
	CoolingCoilValveControl                        PointThing `json:"Cooling Coil Valve Control"`
	CoolingTowerFanSpeed                           PointThing `json:"Cooling Tower Fan Speed"`
	CoolingTowerFault                              PointThing `json:"Cooling Tower Fault"`
	CoolingTowerStartStop                          PointThing `json:"Cooling Tower StartStop"`
	CoolingTowerStatus                             PointThing `json:"Cooling Tower Status"`
	Damper                                         PointThing `json:"Damper"`
	DamperFeedback                                 PointThing `json:"Damper Feedback"`
	DewPoint                                       PointThing `json:"Dew Point"`
	DewPointSetPoint                               PointThing `json:"Dew Point SetPoint"`
	ElectricDuctHeaterStartStop                    PointThing `json:"Electric Duct Heater StartStop"`
	ElectricDuctHeaterStatus                       PointThing `json:"Electric Duct Heater Status"`
	EnergyTarget                                   PointThing `json:"Energy Target"`
	ExhaustAirDamper                               PointThing `json:"Exhaust Air Damper"`
	ExhaustAirFanFault                             PointThing `json:"Exhaust Air Fan Fault"`
	ExhaustAirFanStartStop                         PointThing `json:"Exhaust Air Fan StartStop"`
	ExhaustAirFanStatus                            PointThing `json:"Exhaust Air Fan Status"`
	ExhaustAirConstPressure                        PointThing `json:"Exhaust Air const Pressure"`
	FCUFanFault                                    PointThing `json:"FCU Fan Fault"`
	FCUFanSpeed                                    PointThing `json:"FCU Fan Speed"`
	FCUFanStartStop                                PointThing `json:"FCU Fan StartStop"`
	FCUFanStatus                                   PointThing `json:"FCU Fan Status"`
	FTUFanStartStop                                PointThing `json:"FTU Fan StartStop"`
	FilterDifferentialPressure                     PointThing `json:"Filter Differential Pressure"`
	FilterDifferentialPressureSetPoint             PointThing `json:"Filter Differential Pressure SetPoint"`
	FilterDifferentialPressureStatus               PointThing `json:"Filter Differential Pressure Status"`
	FilterStatus                                   PointThing `json:"Filter Status"`
	FloorOccupancy                                 PointThing `json:"Floor Occupancy"`
	GasTarget                                      PointThing `json:"Gas Target"`
	GeneralExhaustAirFanFault                      PointThing `json:"General Exhaust Air Fan Fault"`
	GeneralExhaustAirFanStartStop                  PointThing `json:"General Exhaust Air Fan StartStop"`
	GeneralExhaustAirFanStatus                     PointThing `json:"General Exhaust Air Fan Status"`
	GeneralExhaustAirConstPressure                 PointThing `json:"General Exhaust Air const Pressure"`
	HeatingCoilValveControl                        PointThing `json:"Heating Coil Valve Control"`
	HotWaterPumpDifferentialPressure               PointThing `json:"Hot Water Pump Differential Pressure"`
	HotWaterPumpDifferentialPressureSetPoint       PointThing `json:"Hot Water Pump Differential Pressure SetPoint"`
	HotWaterPumpFault                              PointThing `json:"Hot Water Pump Fault"`
	HotWaterPumpSpeed                              PointThing `json:"Hot Water Pump Speed"`
	HotWaterPumpStartStop                          PointThing `json:"Hot Water Pump StartStop"`
	HotWaterPumpStatus                             PointThing `json:"Hot Water Pump Status"`
	HotWaterReturnTemperature                      PointThing `json:"Hot Water Return Temperature"`
	HotWaterSupplyTemperature                      PointThing `json:"Hot Water Supply Temperature"`
	HotWaterSupplyTemperatureSetPoint              PointThing `json:"Hot Water Supply Temperature SetPoint"`
	IsolationAirDamper                             PointThing `json:"Isolation Air Damper"`
	KitchenExhaustAirFanFault                      PointThing `json:"Kitchen Exhaust Air Fan Fault"`
	KitchenExhaustAirFanStartStop                  PointThing `json:"Kitchen Exhaust Air Fan StartStop"`
	KitchenExhaustAirFanStatus                     PointThing `json:"Kitchen Exhaust Air Fan Status"`
	KitchenExhaustAirConstPressure                 PointThing `json:"Kitchen Exhaust Air const Pressure"`
	Load                                           PointThing `json:"Load"`
	MaximumAirflow                                 PointThing `json:"Maximum Airflow"`
	MinimumAirflow                                 PointThing `json:"Minimum Airflow"`
	MinimumOutsideAirDamperControl                 PointThing `json:"Minimum Outside Air Damper Control"`
	MinimumOutsideAirDamperFeedback                PointThing `json:"Minimum Outside Air Damper Feedback"`
	MixAirTemperature                              PointThing `json:"Mix Air Temperature"`
	OutsideAirDamperControl                        PointThing `json:"Outside Air Damper Control"`
	OutsideAirEnthalpy                             PointThing `json:"Outside Air Enthalpy"`
	OutsideAirFanFault                             PointThing `json:"Outside Air Fan Fault"`
	OutsideAirFanSpeed                             PointThing `json:"Outside Air Fan Speed"`
	OutsideAirFanStartStop                         PointThing `json:"Outside Air Fan StartStop"`
	OutsideAirFanStatus                            PointThing `json:"Outside Air Fan Status"`
	OutsideAirFlow                                 PointThing `json:"Outside Air Flow"`
	OutsideAirHumidity                             PointThing `json:"Outside Air Humidity"`
	OutsideAirRelativeHumidity                     PointThing `json:"Outside Air Relative Humidity"`
	OutsideAirTemperature                          PointThing `json:"Outside Air Temperature"`
	OutsideAirConstPressure                        PointThing `json:"Outside Air const Pressure"`
	OutsideAirConstPressureSetPoint                PointThing `json:"Outside Air const Pressure SetPoint"`
	OutsideReturnRoomAirTemperature                PointThing `json:"Outside Return Room Air Temperature"`
	PACFanFault                                    PointThing `json:"PAC Fan Fault"`
	PACFanStartStop                                PointThing `json:"PAC Fan StartStop"`
	PACFanStatus                                   PointThing `json:"PAC Fan Status"`
	PrimaryLoopValve                               PointThing `json:"Primary Loop Valve"`
	PrimaryReturnTemperature                       PointThing `json:"Primary Return Temperature"`
	PrimarySupplyTemperature                       PointThing `json:"Primary Supply Temperature"`
	RACo2                                          PointThing `json:"RACO2"`
	ReturnAirCO2                                   PointThing `json:"Return Air CO2"`
	ReturnAirDamper                                PointThing `json:"Return Air Damper"`
	ReturnAirDamperControl                         PointThing `json:"Return Air Damper Control"`
	ReturnAirEnthalpy                              PointThing `json:"Return Air Enthalpy"`
	ReturnAirFanFault                              PointThing `json:"Return Air Fan Fault"`
	ReturnAirFanSpeed                              PointThing `json:"Return Air Fan Speed"`
	ReturnAirFanStartStop                          PointThing `json:"Return Air Fan StartStop"`
	ReturnAirFanStatus                             PointThing `json:"Return Air Fan Status"`
	ReturnAirFlow                                  PointThing `json:"Return Air Flow"`
	ReturnAirHumidity                              PointThing `json:"Return Air Humidity"`
	ReturnAirSuctionPressure                       PointThing `json:"Return Air Suction Pressure"`
	ReturnAirTemp                                  PointThing `json:"Return Air Temperature"`
	ReturnAirTempSetPoint                          PointThing `json:"Return Air Temperature SetPoint"`
	ReturnAirConstPressure                         PointThing `json:"Return Air const Pressure"`
	ReturnAirConstPressureSetPoint                 PointThing `json:"Return Air const Pressure SetPoint"`
	ReturnRoomAirHumidity                          PointThing `json:"Return Room Air Humidity"`
	ReturnRoomAirTemperature                       PointThing `json:"Return Room Air Temperature"`
	ReturnRoomAirTemperatureSetPoint               PointThing `json:"Return Room Air Temperature SetPoint"`
	ReverseValveControl                            PointThing `json:"Reverse Valve Control"`
	RoomAirHumidity                                PointThing `json:"Room Air Humidity"`
	RoomAirTemperature                             PointThing `json:"Room Air Temperature"`
	RoomAirTemperatureSetPoint                     PointThing `json:"Room Air Temperature SetPoint"`
	SecondaryLoopValve                             PointThing `json:"Secondary Loop Valve"`
	SecondaryReturnTemperature                     PointThing `json:"Secondary Return Temperature"`
	SecondarySupplyTemperature                     PointThing `json:"Secondary Supply Temperature"`
	Speed                                          PointThing `json:"Speed"`
	SupplyAirDamper                                PointThing `json:"Supply Air Damper"`
	SupplyAirFanFault                              PointThing `json:"Supply Air Fan Fault"`
	SupplyAirFanSpeed                              PointThing `json:"Supply Air Fan Speed"`
	SupplyAirFanStartStop                          PointThing `json:"Supply Air Fan StartStop"`
	SupplyAirFanStatus                             PointThing `json:"Supply Air Fan Status"`
	SupplyAirFlow                                  PointThing `json:"Supply Air Flow"`
	SupplyAirHumidity                              PointThing `json:"Supply Air Humidity"`
	SupplyAirTemperature                           PointThing `json:"Supply Air Temperature"`
	SupplyAirTemperatureSetPoint                   PointThing `json:"Supply Air Temperature SetPoint"`
	SupplyAirConstPressure                         PointThing `json:"Supply Air const Pressure"`
	SupplyAirConstPressureSetPoint                 PointThing `json:"Supply Air const Pressure SetPoint"`
	SupplyReturnRoomAirHumidity                    PointThing `json:"Supply Return Room Air Humidity"`
	TerminalLoad                                   PointThing `json:"Terminal Load"`
	ToiletExhaustAirFanFault                       PointThing `json:"Toilet Exhaust Air Fan Fault"`
	ToiletExhaustAirFanStartStop                   PointThing `json:"Toilet Exhaust Air Fan StartStop"`
	ToiletExhaustAirFanStatus                      PointThing `json:"Toilet Exhaust Air Fan Status"`
	ToiletExhaustAirConstPressure                  PointThing `json:"Toilet Exhaust Air const Pressure"`
	TotalActiveEnergy                              PointThing `json:"Total Active Energy"`
	TotalActivePower                               PointThing `json:"Total Active Power"`
	TotalGasConsumption                            PointThing `json:"Total Gas Consumption"`
	TotalWaterConsumption                          PointThing `json:"Total Water Consumption"`
	VAVEnable                                      PointThing `json:"VAV Enable"`
	WaterTarget                                    PointThing `json:"Water Target"`
	ZoneAfterHoursActiveTime                       PointThing `json:"Zone After Hours Active Time"`
	ZoneOccupancy                                  PointThing `json:"Zone Occupancy"`
	ZoneRoomRelativeHumidity                       PointThing `json:"Zone Room Relative Humidity"`
	ZoneRoomTemperature                            PointThing `json:"Zone Room Temperature"`
	ZoneRoomTemperatureSetPoint                    PointThing `json:"Zone Room Temperature SetPoint"`
	ZoneTemperatureSetPoint                        PointThing `json:"Zone Temperature SetPoint"`
}

type PointThing struct {
	EquipType       string   `json:"equipType"`
	Kind            string   `json:"kind"`
	Tags            []string `json:"tags"`
	PointThingClass string   `json:"PointThingClass"`
	PointThingType  string   `json:"PointThingType"`
	UnitImperial    string   `json:"unitImperial"`
	Unit            string   `json:"unit"`
	UnitsTo         string   `json:"units_to"`
	Decimal         int      `json:"decimal"`
	Round           *float64 `json:"round"`
	InputMin        *float64 `json:"input_min"`
	InputMax        *float64 `json:"input_max"`
	ScaleMin        *float64 `json:"scale_min"`
	ScaleMax        *float64 `json:"scale_max"`
	Writeable       bool     `json:"writeable"`
	Invert          bool     `json:"invert"`
}

//type EquipTypeEnum string
//
//type KindEnum string
//
//const (
//	Boolean       KindEnum = "boolean"
//	Number        KindEnum = "number"
//	NumberBoolean KindEnum = "number boolean"
//	String        KindEnum = "string"
//)
//
//type Unt string
//
//const (
//	Empty       Unt = ""
//	Temperature Unt = "Temperature"
//)
