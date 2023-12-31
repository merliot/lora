package main

import (
	"bytes"
	"fmt"
	"machine"
	"time"

	"github.com/merliot/lora/lorae5"
)

var sig = []byte("hello")

func main() {
	var i = 0
	time.Sleep(2 * time.Second)
	lora := lorae5.New(machine.UART1, machine.GPIO4, machine.GPIO5, 9600)
	lora.Init()
	for {
		pkts := lora.Rx(2000)
		for _, pkt := range pkts {
			println("received", string(pkt))
			if bytes.Equal(pkt, sig) {
				msg := fmt.Sprintf("%s %d", string(pkt), i)
				i += 1
				err := lora.Tx([]byte(msg), 1000)
				if err != nil {
					println("tx error", err.Error())
				} else {
					println("sent", msg)
				}
			}
		}
	}
}
