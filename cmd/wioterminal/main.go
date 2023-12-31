package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"github.com/merliot/lora/lorae5"
	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/tinyfont/proggy"
	"tinygo.org/x/tinyterm"
)

var sigs = []string{"hello", "jello", "marty"}

var (
	display = ili9341.NewSPI(
		machine.SPI3,
		machine.LCD_DC,
		machine.LCD_SS_PIN,
		machine.LCD_RESET,
	)

	backlight = machine.LCD_BACKLIGHT

	terminal = tinyterm.NewTerminal(display)

	black = color.RGBA{0, 0, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	green = color.RGBA{0, 255, 0, 255}

	font = &proggy.TinySZ8pt7b
)

func main() {
	time.Sleep(2 * time.Second)
	println("running...")

	machine.SPI3.Configure(machine.SPIConfig{
		SCK:       machine.LCD_SCK_PIN,
		SDO:       machine.LCD_SDO_PIN,
		SDI:       machine.LCD_SDI_PIN,
		Frequency: 40000000,
	})

	display.Configure(ili9341.Config{})
	display.FillScreen(black)

	backlight.Configure(machine.PinConfig{machine.PinOutput})
	backlight.High()

	terminal.Configure(&tinyterm.Config{
		Font:       font,
		FontHeight: 10,
		FontOffset: 6,
	})

	lora := lorae5.New(machine.UART4, machine.D0, machine.D1, 9600)
	lora.Init()
	for {
		for _, sig := range sigs {
			time.Sleep(5 * time.Second)
			err := lora.Tx([]byte(sig), 1000)
			if err != nil {
				println("tx error", err.Error())
				fmt.Fprintf(terminal, "tx error %s\n", err)
				continue
			}
			println("sent", sig)
			fmt.Fprintf(terminal, "sent %s\n", sig)
			pkts := lora.Rx(2000)
			for _, pkt := range pkts {
				println("received", string(pkt))
				fmt.Fprintf(terminal, "received %s\n", string(pkt))
			}
		}
	}
}
