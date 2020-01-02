package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/garciademarina/chuncked/pkg/stream"
	"github.com/hybridgroup/mjpeg"
)

func main() {
	var listChannels []chan image.Image
	c := make(chan image.Image)
	listChannels = append(listChannels, c)

	// Setup our Ctrl+C handler
	SetupCloseHandler(listChannels)

	s1 := stream.New(
		"Cam1",
		"http://192.168.1.169:81/stream",
	)

	go func(ch chan image.Image) {
		for {
			err := s1.CaptureFrame(ch)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Reconnecting")
		}
	}(c)

	root := mjpeg.NewStream()
	go combineCamera(root, c)

	mux := http.NewServeMux()
	mux.Handle("/", root)

	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func combineCamera(stream *mjpeg.Stream, ch chan image.Image) {
	fmt.Println("Waiting images")
	frameDuration := 50 * time.Millisecond
	ticker := time.NewTicker(frameDuration)
	tickerDot := time.NewTicker(time.Second)
	var drawDot bool
	for {
		rgba := <-ch
		select {
		case <-tickerDot.C:
			drawDot = !drawDot

		case <-ticker.C:
			m := image.NewRGBA(image.Rect(0, 0, 320, 240))
			draw.Draw(m, m.Bounds(), rgba, image.Point{0, 0}, draw.Src)
			if drawDot {
				m.Set(10, 12, color.RGBA{255, 0, 0, 255})
				m.Set(10, 13, color.RGBA{255, 0, 0, 255})
				m.Set(10, 14, color.RGBA{255, 0, 0, 255})
				m.Set(10, 15, color.RGBA{255, 0, 0, 255})

				m.Set(11, 12, color.RGBA{255, 0, 0, 255})
				m.Set(11, 13, color.RGBA{255, 0, 0, 255})
				m.Set(11, 14, color.RGBA{255, 0, 0, 255})
				m.Set(11, 15, color.RGBA{255, 0, 0, 255})

				m.Set(12, 12, color.RGBA{255, 0, 0, 255})
				m.Set(12, 13, color.RGBA{255, 0, 0, 255})
				m.Set(12, 14, color.RGBA{255, 0, 0, 255})
				m.Set(12, 15, color.RGBA{255, 0, 0, 255})

				m.Set(13, 12, color.RGBA{255, 0, 0, 255})
				m.Set(13, 13, color.RGBA{255, 0, 0, 255})
				m.Set(13, 14, color.RGBA{255, 0, 0, 255})
				m.Set(13, 15, color.RGBA{255, 0, 0, 255})
			}

			buf := new(bytes.Buffer)
			err := jpeg.Encode(buf, m, nil)
			if err == nil {
				sendImg := buf.Bytes()
				stream.UpdateJPEG(sendImg)
			}
		default:
		}
	}
}

// SetupCloseHandler ...
func SetupCloseHandler(listChannels [](chan image.Image)) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		for _, currentChannel := range listChannels {
			fmt.Println("\r- closing channel")
			close(currentChannel)
		}
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}
