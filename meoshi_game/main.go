package main

import (
	"machine"
	"time"

	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
)

type WS2812B struct {
	Pin machine.Pin
	ws  *piolib.WS2812B
}

func NewWS2812B(pin machine.Pin) *WS2812B {
	s, _ := pio.PIO0.ClaimStateMachine()
	ws, _ := piolib.NewWS2812B(s, pin)
	ws.EnableDMA(true)
	return &WS2812B{
		ws: ws,
	}
}

func (ws *WS2812B) WriteRaw(rawGRB []uint32) error {
	return ws.ws.WriteRaw(rawGRB)
}

func main() {
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 2.8 * machine.MHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})

	const WHITE uint32 = 0xFFFFFFFF
	const RED uint32 = 0x00FF00FF
	const KEY_LENGTH = 12

	var initColor = make([]uint32, 12)
	for i := 0; i < KEY_LENGTH; i++ {
		initColor[i] = WHITE
	}

	colPins := []machine.Pin{
		machine.GPIO5,
		machine.GPIO6,
		machine.GPIO7,
		machine.GPIO8,
	}

	rowPins := []machine.Pin{
		machine.GPIO9,
		machine.GPIO10,
		machine.GPIO11,
	}

	for _, c := range colPins {
		c.Configure(machine.PinConfig{Mode: machine.PinOutput})
		c.Low()
	}

	for _, r := range rowPins {
		r.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	ws := NewWS2812B(machine.GPIO1)
	ws.WriteRaw(initColor)

	keyMap := []int{0, 3, 6, 9, 10, 11, 8, 5, 2, 1}

	for {
		for _, keyValue := range keyMap {
			// COL1
			colPins[0].High()
			colPins[1].Low()
			colPins[2].Low()
			colPins[3].Low()
			time.Sleep((1 * time.Millisecond))
			if rowPins[0].Get() {
				if initColor[0] == RED {
					break
				}
			}
			if rowPins[1].Get() {
				if initColor[1] == RED {
					break
				}
			}
			if rowPins[2].Get() {
				if initColor[2] == RED {
					break
				}
			}

			// COL2
			colPins[0].Low()
			colPins[1].High()
			colPins[2].Low()
			colPins[3].Low()
			time.Sleep((1 * time.Millisecond))
			if rowPins[0].Get() {
				if initColor[3] == RED {
					break
				}
			}
			if rowPins[1].Get() {
				if initColor[4] == RED {
					break
				}
			}
			if rowPins[2].Get() {
				if initColor[5] == RED {
					break
				}
			}

			// COL3
			colPins[0].Low()
			colPins[1].Low()
			colPins[2].High()
			colPins[3].Low()
			time.Sleep((1 * time.Millisecond))
			if rowPins[0].Get() {
				if initColor[6] == RED {
					break
				}
			}
			if rowPins[1].Get() {
				if initColor[7] == RED {
					break
				}
			}
			if rowPins[2].Get() {
				if initColor[8] == RED {
					break
				}
			}

			// COL4
			colPins[0].Low()
			colPins[1].Low()
			colPins[2].Low()
			colPins[3].High()
			time.Sleep((1 * time.Millisecond))
			if rowPins[0].Get() {
				if initColor[9] == RED {
					break
				}
			}
			if rowPins[1].Get() {
				if initColor[10] == RED {
					break
				}
			}
			if rowPins[2].Get() {
				if initColor[11] == RED {
					break
				}
			}

			for _, keyValue2 := range keyMap {
				initColor[keyValue2] = WHITE
			}
			initColor[keyValue] = RED
			ws.WriteRaw(initColor)
			time.Sleep(500 * time.Millisecond)
		}
		// TODO: レベルを上げる or ゲームオーバー処理
	}
}
