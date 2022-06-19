// see:http://bashalog.c-brains.jp/20/03/30-170110.php
window.onload = function () {
    let dropZone = document.getElementById('drop-zone')
    let preview = document.getElementById('preview')
    let fileInput = document.getElementById('file-input')
    let imageStatus = document.getElementById("image-status")
    let errorMessage = document.getElementById("error-message")

    dropZone.addEventListener('dragover', function (e) {
        e.stopPropagation()
        e.preventDefault()
        this.style.background = '#e1e7f0'
    }, false)

    dropZone.addEventListener('dragleave', function (e) {
        e.stopPropagation()
        e.preventDefault()
        this.style.background = '#ffffff'
    }, false)

    fileInput.addEventListener('change', function () {
        previewFile(this.files[0])
        imageStatus.innerHTML = "<h2>upload image</h2>"
        errorMessage.innerHTML = ""
    })

    dropZone.addEventListener('drop', function (e) {
        e.stopPropagation()
        e.preventDefault()
        this.style.background = '#ffffff'
        let files = e.dataTransfer.files // get drop file
        if (files.length > 1) return alert('you can only one file to upload ')
        fileInput.files = files //inputのvalueをドラッグしたファイルに置き換える。
        previewFile(files[0])
        imageStatus.innerHTML = "<h2>upload image</h2>"
        errorMessage.innerHTML = ""
    }, false)

    function previewFile(file) {
        let fr = new FileReader()
        fr.readAsDataURL(file)
        fr.onload = function () {
            let img = document.createElement('img')
            let info = document.getElementById('size-info')
            img.setAttribute('src', String(fr.result))
            img.setAttribute('class', 'image')
            img.setAttribute('id', 'preview-image')
            preview.innerHTML = ""
            preview.appendChild(img).onload = function () {
                let previewImg = document.getElementById('preview-image')
                info.innerHTML = previewImg.naturalWidth + "×" + previewImg.naturalHeight
            }
        }
    }
}