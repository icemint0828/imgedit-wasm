package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strconv"
	"syscall/js"
	"time"

	"github.com/icemint0828/imgedit"
)

var window = js.Global()
var document = window.Get("document")

func main() {
	ch := make(chan interface{})
	window.Set("resize", js.FuncOf(resize))
	window.Set("trim", js.FuncOf(trim))
	window.Set("grayscale", js.FuncOf(grayscale))
	window.Set("reverse", js.FuncOf(reverse))
	<-ch
}

func grayscale(_ js.Value, _ []js.Value) interface{} {
	f := func(srcImg image.Image) image.Image {
		c := imgedit.NewConverter(srcImg)
		c.Grayscale()
		return c.Convert()
	}
	fileEdit(f)
	return nil
}

func reverse(_ js.Value, _ []js.Value) interface{} {
	f := func(srcImg image.Image) image.Image {
		verticalVal := getElementById("vertical").Get("checked")

		c := imgedit.NewConverter(srcImg)
		if verticalVal.Bool() {
			c.ReverseY()
		} else {
			c.ReverseX()
		}
		return c.Convert()
	}
	fileEdit(f)
	return nil
}

func trim(_ js.Value, _ []js.Value) interface{} {
	f := func(srcImg image.Image) image.Image {
		leftVal := getElementById("left").Get("value")
		topVal := getElementById("top").Get("value")
		widthVal := getElementById("trim-width").Get("value")
		heightVal := getElementById("trim-height").Get("value")

		left, err := strconv.Atoi(leftVal.String())
		if err != nil {
			return nil
		}
		top, err := strconv.Atoi(topVal.String())
		if err != nil {
			return nil
		}
		width, err := strconv.Atoi(widthVal.String())
		if err != nil {
			return nil
		}
		height, err := strconv.Atoi(heightVal.String())
		if err != nil {
			return nil
		}
		if !((0 <= left && left <= 5000) && (0 <= top && top <= 5000) && (0 <= width && width <= 5000) && (0 <= height && height <= 5000)) {
			return nil
		}
		c := imgedit.NewConverter(srcImg)
		c.Trim(left, top, width, height)
		return c.Convert()
	}
	fileEdit(f)
	return nil
}

func resize(_ js.Value, _ []js.Value) interface{} {
	f := func(srcImg image.Image) image.Image {
		widthVal := getElementById("resize-width").Get("value")
		heightVal := getElementById("resize-height").Get("value")
		ratioVal := getElementById("resize-ratio").Get("value")

		ratio, err := strconv.ParseFloat(ratioVal.String(), 64)
		c := imgedit.NewConverter(srcImg)
		if err == nil {
			if !(0.01 <= ratio && ratio <= 10) {
				return nil
			}
			c.ResizeRatio(ratio)
		} else {
			width, err := strconv.Atoi(widthVal.String())
			if err != nil {
				return nil
			}
			height, err := strconv.Atoi(heightVal.String())
			if err != nil {
				return nil
			}
			if !((0 <= width && width <= 5000) && (0 <= height && height <= 5000)) {
				return nil
			}
			c.Resize(width, height)
		}
		return c.Convert()
	}
	fileEdit(f)
	return nil
}

func fileEdit(f func(srcImg image.Image) image.Image) {
	message := getElementById("error-message")
	message.Set("innerHTML", "image editing now")
	go func() {
		status := getElementById("image-status")
		fileInput := getElementById("file-input")
		preview := getElementById("preview")
		item := fileInput.Get("files").Call("item", 0)
		if item.IsNull() {
			time.Sleep(1 * time.Second)
			message.Set("innerHTML", "file not found")
			return
		}

		item.Call("arrayBuffer").Call("then", js.FuncOf(func(v js.Value, x []js.Value) any {
			srcData := window.Get("Uint8Array").New(x[0])
			src := make([]byte, srcData.Get("length").Int())
			js.CopyBytesToGo(src, srcData)
			srcImg, format, err := image.Decode(bytes.NewBuffer(src))
			if err != nil {
				time.Sleep(1 * time.Second)
				message.Set("innerHTML", "unsupported file")
				return nil
			}
			dstImg := f(srcImg)
			if dstImg == nil {
				time.Sleep(1 * time.Second)
				message.Set("innerHTML", "invalid parameter")
				return nil
			}

			dstBuf := &bytes.Buffer{}
			switch format {
			case "png":
				png.Encode(dstBuf, dstImg)
			case "jpeg":
				jpeg.Encode(dstBuf, dstImg, &jpeg.Options{Quality: 100})
			case "gif":
				gif.Encode(dstBuf, dstImg, &gif.Options{NumColors: 256})
			}
			var dstData = window.Get("Uint8Array").New(dstBuf.Len())
			js.CopyBytesToJS(dstData, dstBuf.Bytes())
			window.Call("previewBlob", dstData.Get("buffer"))
			message.Set("innerHTML", "edit success")
			status.Set("innerHTML", "preview image")
			preview.Call("setAttribute", "data-state", "onPreview")
			return nil
		}))
	}()
	return
}

func getElementById(id string) js.Value {
	return document.Call("getElementById", id)
}
