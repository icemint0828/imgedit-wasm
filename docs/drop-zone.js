// see:http://bashalog.c-brains.jp/20/03/30-170110.php
window.onload = function () {
    let dropZone = document.getElementById('drop-zone')
    let preview = document.getElementById('preview')
    let fileInput = document.getElementById('file-input')
    let imageStatus = document.getElementById("image-status")
    let errorMessage = document.getElementById("error-message")
    let info = document.getElementById('size-info')
    let fileReset = document.getElementById("file-reset")
    let download = document.getElementById("download")

    dropZone.addEventListener('dragover', function (e) {
        e.stopPropagation()
        e.preventDefault()
    }, false)

    dropZone.addEventListener('dragleave', function (e) {
        e.stopPropagation()
        e.preventDefault()
    }, false)

    fileInput.addEventListener('change', function () {
        if (this.files.length === 1) {
            previewFile(this.files[0])
            imageStatus.innerHTML = "upload image"
            errorMessage.innerHTML = ""
            preview.dataset.state = statusUpload
            return
        }
        preview.innerHTML = ""
        imageStatus.innerHTML = ""
        errorMessage.innerHTML = ""
        info.innerHTML = ""
        preview.dataset.state = statusNone
        toggleButtonVisible(false)
        delete preview.dataset.originFormat
    })

    dropZone.addEventListener('drop', function (e) {
        e.stopPropagation()
        e.preventDefault()
        let files = e.dataTransfer.files // get drop file
        if (files.length > 1) return alert('you can only one file to upload ')
        fileInput.files = files //inputのvalueをドラッグしたファイルに置き換える。
        previewFile(files[0])
        imageStatus.innerHTML = "upload image"
        errorMessage.innerHTML = ""
        preview.dataset.state = statusUpload
    }, false)

    function previewFile(file) {
        let fr = new FileReader()
        fr.readAsDataURL(file)
        fr.onload = function () {
            let img = document.createElement('img')
            img.setAttribute('src', String(fr.result))
            img.setAttribute('class', 'image')
            img.setAttribute('id', 'preview-image')
            preview.dataset.originFormat = String(fr.result.slice(fr.result.indexOf("/") + 1, fr.result.indexOf(";")))
            preview.innerHTML = ""
            preview.appendChild(img).onload = function () {
                let previewImg = document.getElementById('preview-image')
                info.innerHTML = "("+ previewImg.naturalWidth + "×" + previewImg.naturalHeight + ")"
                let tileInputX = document.getElementById('x-length')
                let tileInputY = document.getElementById('y-length')
                tileInputX.setAttribute("max", String(Math.max( 1, Math.trunc(5000 / previewImg.naturalWidth))))
                tileInputY.setAttribute("max", String(Math.max( 1, Math.trunc(5000 / previewImg.naturalHeight))))
                toggleButtonVisible(true)
            }
        }
    }

    function toggleButtonVisible(isShow) {
        if (isShow) {
            fileReset.classList.remove("hide")
            fileReset.classList.add("show")
            download.classList.remove("hide")
            download.classList.add("show")
        } else {
            fileReset.classList.remove("show")
            fileReset.classList.add("hide")
            download.classList.remove("show")
            download.classList.add("hide")
        }
    }
}
