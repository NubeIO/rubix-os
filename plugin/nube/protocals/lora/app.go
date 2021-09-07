package main

import "github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/tty"

func SerialOpenAndRead() {
	bb := new(tty.SerialSetting)
	bb.BaudRate = 9600
	aa := tty.New(bb)
	aa.NewSerialConnection()
	aa.Loop()
}


