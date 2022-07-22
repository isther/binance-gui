package utils

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/gen2brain/malgo"
	"github.com/hajimehoshi/go-mp3"
	"github.com/isther/binanceGui/console"
)

func PlayAMusic(music string) {
	file, err := os.Open(music)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}
	defer file.Close()

	var reader io.Reader
	var channels, sampleRate uint32

	m, err := mp3.NewDecoder(file)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}

	reader = m
	channels = 2
	sampleRate = uint32(m.SampleRate())

	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = channels
	deviceConfig.SampleRate = sampleRate
	deviceConfig.Alsa.NoMMap = 1

	// This is the function that's used for sending more data to the device for playback.
	onSamples := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		io.ReadFull(reader, pOutputSample)
	}

	deviceCallbacks := malgo.DeviceCallbacks{
		Data: onSamples,
	}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}
	defer device.Uninit()

	err = device.Start()
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}

	time.Sleep(3 * time.Second)
}

func Float64ToStringLen3(f float64) string {
	var s = fmt.Sprintf("%.8f", f)
	if f >= 1.0 {
		if f >= 100.0 {
			s = s[:3]
		} else {
			s = s[:4]
		}
	} else {
		for i := 0; i < len(s); i++ {
			if s[i] == '.' {
				continue
			}

			if s[i] != '0' {
				if i+3 <= len(s) {
					s = s[:i+3]
				}
				break
			}

		}

		var pos = len(s) - 1
		for ; pos > 0; pos-- {
			if s[pos] != '0' {
				break
			}
		}
		s = s[:pos+1]
	}
	if s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}

func correction(val float64, size string) string {
	var (
		oneIdx    = strings.Index(size, "1")
		pointIdx  = strings.Index(size, ".")
		precision int
		resStr    string
	)
	if oneIdx < pointIdx {
		precision = oneIdx - pointIdx + 1
		resStr = fmt.Sprintf("%.8f", RoundLower(val, precision))
		pointIdx = strings.Index(resStr, ".")
		resStr = resStr[:pointIdx]
		return resStr
	}

	precision = oneIdx - pointIdx
	resStr = fmt.Sprintf("%.8f", RoundLower(val, precision))

	return resStr
}

func RoundLower(val float64, precision int) float64 {
	if precision == 0 {
		return math.Round(val)
	}

	p := math.Pow10(precision)
	if precision < 0 {
		return math.Floor(val*p) * math.Pow10(-precision)
	}

	return math.Floor(val*p) / p
}

func RoundUpper(val float64, precision int) float64 {
	if precision == 0 {
		return math.Round(val)
	}

	p := math.Pow10(precision)
	if precision < 0 {
		return math.Floor(val*p+0.5) * math.Pow10(-precision)
	}

	return math.Floor(val*p+0.5) / p
}
