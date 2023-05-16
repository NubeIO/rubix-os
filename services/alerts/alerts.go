package alerts

import (
	"errors"
)

type AlertStatus string
type AlertEntityType string
type AlertType string
type AlertSeverity string

const (
	Active       AlertStatus = "active"
	Acknowledged AlertStatus = "acknowledged"
	Closed       AlertStatus = "closed"
)

const (
	Gateway AlertEntityType = "gateway"
	Network AlertEntityType = "network"
	Device  AlertEntityType = "device"
	Point   AlertEntityType = "point"
	Service AlertEntityType = "service"
)

const (
	Ping      AlertType = "ping"
	Fault     AlertType = "fault"
	Threshold AlertType = "threshold"
	FlatLine  AlertType = "flat-line"
)

const (
	Crucial AlertSeverity = "crucial"
	Minor   AlertSeverity = "minor"
	Info    AlertSeverity = "info"
	Warning AlertSeverity = "warning"
)

func CheckStatus(s string) error {
	switch AlertStatus(s) {
	case Active:
		return nil
	case Acknowledged:
		return nil
	case Closed:
		return nil
	}
	return errors.New("invalid alert status, try active, acknowledged, closed")
}

func CheckSeverity(s string) error {
	switch AlertSeverity(s) {
	case Crucial:
		return nil
	case Minor:
		return nil
	case Info:
		return nil
	case Warning:
		return nil
	}
	return errors.New("invalid alert status, try crucial, info, warning")
}

func CheckStatusClosed(s string) bool {
	return AlertStatus(s) == Closed
}

func CheckEntityType(s string) error {
	switch AlertEntityType(s) {
	case Gateway:
		return nil
	case Network:
		return nil
	case Device:
		return nil
	case Point:
		return nil
	case Service:
		return nil
	}
	return errors.New("invalid alert entity type, try gateway, network")
}

func CheckAlertType(s string) error {
	switch AlertType(s) {
	case Ping:
		return nil
	case Fault:
		return nil
	case Threshold:
		return nil
	case FlatLine:
		return nil
	}
	return errors.New("invalid alert type, try ping, threshold, fault")
}

func AlertTypeMessage(s string) string {
	switch AlertType(s) {
	case Ping:
		return "failed to ping the device"
	case Fault:
		return ""
	case Threshold:
		return "out of range threshold"
	case FlatLine:
		return ""
	}
	return ""
}
