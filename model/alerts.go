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
