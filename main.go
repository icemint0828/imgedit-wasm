package main

import (
	"bytes"
	"errors"
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
	window.Set("filter", js.FuncOf(filter))
	window.Set("extension", js.FuncOf(extension))
	window.Set("tile", js.FuncOf(tile))
	<-ch
}

func tile(_ js.Value, _ []js.Value) interface{} {
	f := func(bc imgedit.ByteConverter) error {
		xLengthInput := getElementById("x-length")
		yLengthInput := getElementById("y-length")

		xLength, err := strconv.Atoi(xLengthInput.Get("value").String())
		if err != nil {
			return nil
		}
		yLength, err := strconv.Atoi(yLengthInput.Get("value").String())
		if err != nil {
			return nil
		}
		xLengthMax, err := strconv.Atoi(xLengthInput.Get("max").String())
		if err != nil {
			return nil
		}
		yLengthMax, err := strconv.Atoi(yLengthInput.Get("max").String())
		if err != nil {
			return nil
		}
		if !((1 <= xLength && xLength <= xLengthMax) && (1 <= yLength && yLength <= yLengthMax)) {
			return errors.New("[ERR]Enter a value between 1 and " + strconv.Itoa(xLengthMax) + " for rows and between 1 and " + strconv.Itoa(yLengthMax) + " for cols.")
		}
		bc.Tile(xLength, yLength)
		return nil
	}
	fileEdit(f, "")
	return nil
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
	fileEdit(nil, extension)
	return nil
}

func filter(_ js.Value, _ []js.Value) interface{} {
	f := func(bc imgedit.ByteConverter) error {
		elements := getElementsByName("filter")
		var val string
		for i := 0; i < elements.Length(); i++ {
			element := elements.Index(i)
			checked := element.Get("checked")
			if checked.Bool() {
				val = element.Get("value").String()
				break
			}
		}
		var model imgedit.FilterModel
		switch val {
		case "gray":
			model = imgedit.GrayModel
		case "sepia":
			model = imgedit.SepiaModel
		}
		bc.Filter(model)
		return nil
	}
	fileEdit(f, "")
	return nil
}

func reverse(_ js.Value, _ []js.Value) interface{} {
	f := func(bc imgedit.ByteConverter) error {
		verticalVal := getElementById("vertical").Get("checked")
		bc.Reverse(!verticalVal.Bool())
		return nil
	}
	fileEdit(f, "")
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
	fileEdit(f, "")
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
	fileEdit(f, "")
	return nil
}

func fileEdit(f func(bc imgedit.ByteConverter) error, dstFormat imgedit.Extension) {
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
			if f != nil {
				err = f(bc)
				if err != nil {
					time.Sleep(1 * time.Second)
					message.Set("innerHTML", err.Error())
					return nil
				}
			}

			if dstFormat != "" {
				format = dstFormat
			}

			dstBuf := &bytes.Buffer{}
			_ = bc.WriteAs(dstBuf, format)

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

func getElementById(id string) js.Value {
	return document.Call("getElementById", id)
}

func getElementsByName(name string) js.Value {
	return document.Call("getElementsByName", name)
}
