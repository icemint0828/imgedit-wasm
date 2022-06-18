package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strconv"
	"syscall/js"

	"github.com/icemint0828/imgedit"
)

// グローバルオブジェクト（Webブラウザはwindow）の取得
var window = js.Global()
var document = window.Get("document")

func main() {
	ch := make(chan interface{})
	window.Set("resize", js.FuncOf(resize))
	window.Set("grayscale", js.FuncOf(grayscale))
	window.Set("trim", js.FuncOf(grayscale))
	window.Set("reverse", js.FuncOf(grayscale))
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

func resize(_ js.Value, _ []js.Value) interface{} {
	f := func(srcImg image.Image) image.Image {
		c := imgedit.NewConverter(srcImg)

		widthVal := getElementById("width").Get("value")
		heightVal := getElementById("height").Get("value")

		if widthVal.IsNaN() || heightVal.IsNaN() {
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
		c.Resize(width, height)
		return c.Convert()
	}
	fileEdit(f)
	return nil
}

func fileEdit(f func(srcImg image.Image) image.Image) {
	fileInput := getElementById("file-input")
	message := getElementById("error-message")
	item := fileInput.Get("files").Call("item", 0)
	if item.IsNull() {
		message.Set("innerHTML", "file not found")
		return
	}

	item.Call("arrayBuffer").Call("then", js.FuncOf(func(v js.Value, x []js.Value) any {
		srcData := window.Get("Uint8Array").New(x[0])
		src := make([]byte, srcData.Get("length").Int())
		js.CopyBytesToGo(src, srcData)
		srcImg, fmt, err := image.Decode(bytes.NewBuffer(src))
		if err != nil {
			message.Set("innerHTML", " unsupported file")
			return nil
		}
		dstImg := f(srcImg)
		if dstImg == nil {
			message.Set("innerHTML", " invalid parameter")
			return nil
		}

		dstBuf := &bytes.Buffer{}
		switch fmt {
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
		return nil
	}))

	return
}

func getElementById(id string) js.Value {
	return document.Call("getElementById", id)
}
