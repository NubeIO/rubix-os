package main

type commandTopics string

const txSource = "ff"

const (
	topicCommandWhoIs commandTopics = "bacnet/cmd/whois"
	topicCommandRead  commandTopics = "bacnet/cmd/read_value"
	topicCommandWrite commandTopics = "bacnet/cmd/write_value"
)

const presentValue = 85
const presentValueName = "PV"

const priorityArray = 87
const priorityArrayName = "PRI"

const objectName = 77
const objectNameName = "NAME"

const objectType = 79
const objectTypeName = "TYPE"

const objectList = 76
const objectListName = "OBJECTS"

const analogueInput = 0
const analogueInputName = "AI"
const analogueOutput = 1
const analogueOutputName = "AO"
const analogueValue = 2
const analogueValueName = "AV"

const binaryInput = 3
const binaryInputName = "BI"
const binaryOutput = 4
const binaryOutputName = "BO"
const binaryValue = 5
const binaryValueName = "BV"

const multiStateInput = 13
const multiStateName = "MSI"
const multiStateOutput = 14
const multiStateOutputName = "MSO"
const multiStateValue = 19
const multiStateValueName = "MSV"
