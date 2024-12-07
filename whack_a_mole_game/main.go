package main

import (
	"image/color"
	"machine"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
	"tinygo.org/x/tinyfont/gophers"
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

func displayCharacter(
	display ssd1306.Device,
	characterColor color.RGBA,
	character string,
	x int16,
	y int16,
) {
	display.ClearDisplay()
	tinyfont.WriteLine(&display, &freemono.Bold9pt7b, x, y, character, characterColor)
	display.Display()
}

func main() {
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 2.8 * machine.MHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})

	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address:  0x3C,
		Width:    128,
		Height:   64,
		Rotation: drivers.Rotation180,
	})
	display.ClearDisplay()
	time.Sleep(50 * time.Millisecond)

	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	tinyfont.WriteLine(&display, &gophers.Regular32pt, 5, 50, "ABCEF", white)
	display.Display()

	time.Sleep(2 * time.Second)

	const RED uint32 = 0x00FF00FF
	const WHITE uint32 = 0xFFFFFFFF

	colors := []uint32{
		RED, RED, RED, RED,
		RED, RED, RED, RED,
		RED, RED, RED, RED,
	}

	correctColor := []uint32{
		WHITE, WHITE, WHITE, WHITE,
		WHITE, WHITE, WHITE, WHITE,
		WHITE, WHITE, WHITE, WHITE,
	}

	// NOTE: 一旦白色にする
	ws := NewWS2812B(machine.GPIO1)
	ws.WriteRaw(correctColor)

	// NOTE: ランダムに白にする位置を決める
	whiteColorIdxs := mapset.NewSet[int]()
	for whiteColorIdxs.Cardinality() < 6 {
		randomNumber := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomIdx := randomNumber.Intn(len(colors))

		whiteColorIdxs.Add(randomIdx)
	}

	whiteColorIdxSlice := whiteColorIdxs.ToSlice()
	for _, val := range whiteColorIdxSlice {
		colors[val] = WHITE
	}

	enc := encoders.NewQuadratureViaInterrupt(
		machine.GPIO3,
		machine.GPIO4,
	)

	enc.Configure(encoders.QuadratureConfig{
		Precision: 4,
	})

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

	for _, c := range rowPins {
		c.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	// NOTE: カウントダウン表示
	tick := time.Tick(1 * time.Second)
	for countdown := 3; countdown > 0; countdown-- {
		display.ClearDisplay()
		tinyfont.WriteLine(&display, &freemono.Bold24pt7b, 50, 50, strconv.Itoa(countdown), white)
		display.Display()
		<-tick
	}

	display.ClearDisplay()
	tinyfont.WriteLine(&display, &freemono.Bold24pt7b, 35, 50, "Go!!!", white)
	display.Display()

	clearTime := time.Now()

	for {
		ws.WriteRaw(colors)

		// COL1
		colPins[0].High()
		colPins[1].Low()
		colPins[2].Low()
		colPins[3].Low()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			var colorCode0 uint32 = WHITE
			if colors[0] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[0] = colorCode0
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}
		if rowPins[1].Get() {
			var colorCode1 uint32 = WHITE
			if colors[1] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[1] = colorCode1
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}
		if rowPins[2].Get() {
			var colorCode2 uint32 = WHITE
			if colors[2] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[2] = colorCode2
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}

		// COL2
		colPins[0].Low()
		colPins[1].High()
		colPins[2].Low()
		colPins[3].Low()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			var colorCode3 uint32 = WHITE
			if colors[3] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[3] = colorCode3
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}
		if rowPins[1].Get() {
			var colorCode4 uint32 = WHITE
			if colors[4] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[4] = colorCode4
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}
		if rowPins[2].Get() {
			var colorCode5 uint32 = WHITE
			if colors[5] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[5] = colorCode5
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}

		// COL3
		colPins[0].Low()
		colPins[1].Low()
		colPins[2].High()
		colPins[3].Low()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			var colorCode6 uint32 = WHITE
			if colors[6] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[6] = colorCode6
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}
		if rowPins[1].Get() {
			var colorCode7 uint32 = WHITE
			if colors[7] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[7] = colorCode7
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}
		if rowPins[2].Get() {
			var colorCode8 uint32 = WHITE
			if colors[8] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[8] = colorCode8
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}

		// COL4
		colPins[0].Low()
		colPins[1].Low()
		colPins[2].Low()
		colPins[3].High()
		time.Sleep(1 * time.Millisecond)

		if rowPins[0].Get() {
			var colorCode9 uint32 = WHITE
			if colors[9] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[9] = colorCode9
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}
		if rowPins[1].Get() {
			var colorCode10 uint32 = WHITE
			if colors[10] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[10] = colorCode10
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}
		if rowPins[2].Get() {
			var colorCode11 uint32 = WHITE
			if colors[11] == WHITE {
				displayCharacter(display, white, "GameOver...", 5, 45)
				break
			}
			colors[11] = colorCode11
			ws.WriteRaw(colors)
			time.Sleep(100 * time.Millisecond)
		}

		// NOTE: 全て白にできていればクリア
		if reflect.DeepEqual(colors, correctColor) {
			display.ClearDisplay()
			tinyfont.WriteLine(&display, &freemono.Bold9pt7b, 5, 45, "GameClear!!", white)
			tinyfont.WriteLine(
				&display,
				&freemono.Bold9pt7b,
				5,
				15,
				time.Since(clearTime).Truncate(time.Millisecond).String(),
				white,
			)
			display.Display()
			break
		}
	}
}
