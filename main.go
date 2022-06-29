package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	"github.com/icemint0828/imgedit"
)

var window = js.Global()
var document = window.Get("document")

var ErrInvalidParam = errors.New("invalid parameter")

func main() {
	ch := make(chan interface{})
	window.Set("resize", js.FuncOf(resize))
	window.Set("trim", js.FuncOf(trim))
	window.Set("reverse", js.FuncOf(reverse))
	window.Set("grayscale", js.FuncOf(grayscale))
	window.Set("extension", js.FuncOf(extension))
	<-ch
}

func extension(_ js.Value, _ []js.Value) interface{} {
	elements := getElementsByName("extension")
	var extension imgedit.Extension
	for i := 0; i < elements.Length(); i++ {
		element := elements.Index(i)
		checked := element.Get("checked")
		if checked.Bool() {
			extension = imgedit.Extension(element.Get("value").String())
			break
		}
	}
	fmt.Println(extension)
	fileConvert(extension)
	return nil
}

func grayscale(_ js.Value, _ []js.Value) interface{} {
	f := func(bc imgedit.ByteConverter) error {
		bc.Grayscale()
		return nil
	}
	fileEdit(f)
	return nil
}

func reverse(_ js.Value, _ []js.Value) interface{} {
	f := func(bc imgedit.ByteConverter) error {
		verticalVal := getElementById("vertical").Get("checked")

		if verticalVal.Bool() {
			bc.ReverseY()
		} else {
			bc.ReverseX()
		}
		return nil
	}
	fileEdit(f)
	return nil
}

func trim(_ js.Value, _ []js.Value) interface{} {
	f := func(bc imgedit.ByteConverter) error {
		leftVal := getElementById("left").Get("value")
		topVal := getElementById("top").Get("value")
		widthVal := getElementById("trim-width").Get("value")
		heightVal := getElementById("trim-height").Get("value")

		left, err := strconv.Atoi(leftVal.String())
		if err != nil {
			return ErrInvalidParam
		}
		top, err := strconv.Atoi(topVal.String())
		if err != nil {
			return ErrInvalidParam
		}
		width, err := strconv.Atoi(widthVal.String())
		if err != nil {
			return ErrInvalidParam
		}
		height, err := strconv.Atoi(heightVal.String())
		if err != nil {
			return ErrInvalidParam
		}
		if !((0 <= left && left <= 5000) && (0 <= top && top <= 5000) && (0 <= width && width <= 5000) && (0 <= height && height <= 5000)) {
			return nil
		}
		bc.Trim(left, top, width, height)
		return nil
	}
	fileEdit(f)
	return nil
}

func resize(_ js.Value, _ []js.Value) interface{} {
	f := func(bc imgedit.ByteConverter) error {
		widthVal := getElementById("resize-width").Get("value")
		heightVal := getElementById("resize-height").Get("value")
		ratioVal := getElementById("resize-ratio").Get("value")

		ratio, err := strconv.ParseFloat(ratioVal.String(), 64)
		if err == nil {
			if !(0.01 <= ratio && ratio <= 10) {
				return nil
			}
			bc.ResizeRatio(ratio)
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
			bc.Resize(width, height)
		}
		return nil
	}
	fileEdit(f)
	return nil
}

func fileEdit(f func(bc imgedit.ByteConverter) error) {
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
			bc, format, err := imgedit.NewByteConverter(bytes.NewBuffer(src))
			if err != nil {
				time.Sleep(1 * time.Second)
				message.Set("innerHTML", "unsupported file")
				return nil
			}
			err = f(bc)
			if err != nil {
				time.Sleep(1 * time.Second)
				message.Set("innerHTML", err.Error())
				return nil
			}

			dstBuf := &bytes.Buffer{}
			_ = bc.WriteAs(dstBuf, format)
			fmt.Println(format)

			var dstData = window.Get("Uint8Array").New(dstBuf.Len())
			js.CopyBytesToJS(dstData, dstBuf.Bytes())
			window.Call("previewBlob", dstData.Get("buffer"), string(format))
			message.Set("innerHTML", "edit success")
			status.Set("innerHTML", "preview image")
			preview.Call("setAttribute", "data-state", "onPreview")
			return nil
		}))
	}()
	return
}

func fileConvert(dstFormat imgedit.Extension) {
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
			bc, _, err := imgedit.NewByteConverter(bytes.NewBuffer(src))
			if err != nil {
				time.Sleep(1 * time.Second)
				message.Set("innerHTML", "unsupported file")
				return nil
			}

			dstBuf := &bytes.Buffer{}
			_ = bc.WriteAs(dstBuf, dstFormat)

			var dstData = window.Get("Uint8Array").New(dstBuf.Len())
			js.CopyBytesToJS(dstData, dstBuf.Bytes())
			window.Call("previewBlob", dstData.Get("buffer"), string(dstFormat))
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

func getElementsByName(name string) js.Value {
	return document.Call("getElementsByName", name)
}
