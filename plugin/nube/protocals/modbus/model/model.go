package modmodel

import (
	"github.com/NubeIO/null"
)

type Priority struct {
	P1  null.Float `json:"_1,omitempty"` //would be better if we stored the TS and where it was written from, for example from a Remote Producer
	P2  null.Float `json:"_2,omitempty"`
	P3  null.Float `json:"_3,omitempty"`
	P4  null.Float `json:"_4,omitempty"`
	P5  null.Float `json:"_5,omitempty"`
	P6  null.Float `json:"_6,omitempty"`
	P7  null.Float `json:"_7,omitempty"`
	P8  null.Float `json:"_8,omitempty"`
	P9  null.Float `json:"_9,omitempty"`
	P10 null.Float `json:"_10,omitempty"`
	P11 null.Float `json:"_11,omitempty"`
	P12 null.Float `json:"_12,omitempty"`
	P13 null.Float `json:"_13,omitempty"`
	P14 null.Float `json:"_14,omitempty"`
	P15 null.Float `json:"_15,omitempty"`
	P16 null.Float `json:"_16"` //removed and added to the point to save one DB write

}
