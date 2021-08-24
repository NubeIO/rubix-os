package model

type AlertCategory struct {
	Type string //loss of data, offline

}

type AlertType struct {
	Type string //point, device

}

//Alert alerts TODO add in later
type Alert struct {
	Type string //cov, interval, cov_interval
	Duration int

}
type T struct {
	Timestamp int64  `json:"timestamp"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	AlertType string `json:"alert_type"`
	Priority  string `json:"priority"`
}