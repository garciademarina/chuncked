package stream

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/hybridgroup/mjpeg"
	"github.com/pbnjay/pixfont"
	"github.com/pkg/errors"
)

// Stream ...
type Stream struct {
	name   string
	url    string
	stream *mjpeg.Stream
}

// New ...
func New(name string, url string) *Stream {
	return &Stream{
		name: name,
		url:  url,
	}
}

// CaptureFrame ...
func (c *Stream) CaptureFrame(ch chan image.Image) error {
	var resp *http.Response

	tr := &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 20,
	}
	client := http.Client{
		Transport: tr,
	}

	resp, err := client.Get(c.url)
	if err != nil {
		return errors.Wrap(err, "CaptureFrame")
	}

	// extract boundary
	_, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	mr := multipart.NewReader(resp.Body, params["boundary"])
	for part, err := mr.NextPart(); err == nil; part, err = mr.NextPart() {
		value, err := ioutil.ReadAll(part)
		if err != nil {
			return errors.Wrap(err, "CaptureFrame (readBytes)")
		}

		img, _, err := image.Decode(bytes.NewReader(value))
		if err != nil {
			fmt.Println(errors.Wrap(err, "CaptureFrame (lost frame)"))
		} else {

			date := time.Now().Format("2 Jan 2006 15:04:05")

			// add title to image
			m := image.NewRGBA(image.Rect(0, 0, 320, 240))
			draw.Draw(m, m.Bounds(), img, image.Point{0, 0}, draw.Src)
			addLabel(m, 20, 10, fmt.Sprintf("%s %s", c.name, date))

			ch <- m
		}
	}
	return nil

}

func addLabel(img *image.RGBA, x, y int, label string) {
	pixfont.DrawString(img, x, y, label, color.White)
}
