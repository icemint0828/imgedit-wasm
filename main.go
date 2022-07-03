package main

import (
	"bytes"
	"errors"
	"math"
	"strconv"
	"syscall/js"
	"time"

	"github.com/icemint0828/imgedit"
)

const (
	MaxPixels = 5000
	MinPixels = 0
	MaxRatio  = 10
	MinRatio  = 0.1
	StepInt   = 1
	StepFloat = 0.1
)

var window = js.Global()
var document = window.Get("document")

var ErrInvalidParam = errors.New("invalid parameter")

var MaxRowsCount = 100
var MaxColsCount = 100
var MinCount = 1

func main() {
	ch := make(chan interface{})
	window.Set("resize", js.FuncOf(resize))
	window.Set("trim", js.FuncOf(trim))
	window.Set("reverse", js.FuncOf(reverse))
	window.Set("filter", js.FuncOf(filter))
	window.Set("extension", js.FuncOf(extension))
	window.Set("tile", js.FuncOf(tile))
	window.Set("setValidValues", js.FuncOf(setValidValues))
	window.Set("adjustmentTile", js.FuncOf(adjustmentTile))
	<-ch
}

func setValidValues(_ js.Value, _ []js.Value) interface{} {
	// resize
	setValidPixels("resize-width")
	setValidPixels("resize-height")
	setValidRatio("resize-ratio")

	// trim
	setValidPixels("left")
	setValidPixels("top")
	setValidPixels("trim-height")
	setValidPixels("trim-width")

	// tile
	setValidCount("cols")
	setValidCount("rows")
	return nil
}

func adjustmentTile(_ js.Value, _ []js.Value) interface{} {
	previewImg := getElementById("preview-image")

	naturalWidth := previewImg.Get("naturalWidth").Int()
	naturalHeight := previewImg.Get("naturalHeight").Int()
	MaxColsCount = int(math.Max(float64(MinCount), math.Trunc(float64(MaxPixels)/float64(naturalWidth))))
	MaxRowsCount = int(math.Max(float64(MinCount), math.Trunc(float64(MaxPixels)/float64(naturalHeight))))

	setValidCount("cols")
	setValidCount("rows")
	return nil
}

func setValidPixels(id string) {
	v := getElementById(id)
	v.Call("setAttribute", "max", MaxPixels)
	v.Call("setAttribute", "min", MinPixels)
	v.Call("setAttribute", "step", StepInt)
}

func setValidRatio(id string) {
	v := getElementById(id)
	v.Call("setAttribute", "max", MaxRatio)
	v.Call("setAttribute", "min", MinRatio)
	v.Call("setAttribute", "step", StepFloat)
}

func setValidCount(id string) {
	v := getElementById(id)
	var MaxCount int
	if id == "rows" {
		MaxCount = MaxRowsCount
	} else {
		MaxCount = MaxColsCount
	}
	v.Call("setAttribute", "max", MaxCount)
	v.Call("setAttribute", "min", MinCount)
	v.Call("setAttribute", "step", StepInt)
}

func tile(_ js.Value, _ []js.Value) interface{} {
	f := func(bc imgedit.ByteConverter) error {
		colsInput := getElementById("cols")
		rowsInput := getElementById("rows")

		cols, err := strconv.Atoi(colsInput.Get("value").String())
		if err != nil {
			return nil
		}
		rows, err := strconv.Atoi(rowsInput.Get("value").String())
		if err != nil {
			return nil
		}
		colsMax, err := strconv.Atoi(colsInput.Get("max").String())
		if err != nil {
			return nil
		}
		rowsMax, err := strconv.Atoi(rowsInput.Get("max").String())
		if err != nil {
			return nil
		}
		if !((MinCount <= cols && cols <= MaxColsCount) && (MinCount <= rows && rows <= MaxRowsCount)) {
			return errors.New("[ERR]Enter a value between 1 and " + strconv.Itoa(colsMax) + " for rows and between 1 and " + strconv.Itoa(rowsMax) + " for cols.")
		}
		bc.Tile(cols, rows)
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
		if !((MinPixels <= left && left <= MaxPixels) &&
			(MinPixels <= top && top <= MaxPixels) &&
			(MinPixels <= width && width <= MaxPixels) &&
			(MinPixels <= height && height <= MaxPixels)) {
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
			if !(MinRatio <= ratio && ratio <= MaxRatio) {
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
			if !((MinPixels <= width && width <= MaxPixels) && (MinPixels <= height && height <= MaxPixels)) {
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
