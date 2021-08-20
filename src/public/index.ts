let DownloadURL: string = "null"

const LinkBtn: HTMLElement = document.getElementById("linkbtn")

const upload: Function = () => {
    const input: HTMLInputElement = <HTMLInputElement>document.getElementById("input")
    const file: File = input.files[0];

    if (input.value != "") {
        const form: FormData = new FormData();
        form.append("file", file);

        document.getElementById("upload").innerHTML = "Uploading..."

        fetch("/upload", {
            method: "POST",
            body: form
        }).then((data: Response) => {
            console.log(data)
            data.text().then((url: string) => {
                console.log(url)
                DownloadURL = url

                document.getElementById("upload").innerHTML = "Uploaded!"
                document.getElementById("uploadbtn").style.display = "none"
                LinkBtn.style.display = "block"
            })
        }).catch((err: Error) => {
            console.log(err)
            alert("Error: " + err)
        })
    }
}

const CopyLink: Function = () => {
    if (DownloadURL != "null") {
        navigator.clipboard.writeText(document.location.origin + DownloadURL).then(() => {
            LinkBtn.innerHTML = "Link copied!"
            setTimeout(() => {
                LinkBtn.innerHTML = "Copy link"
            }, 2000)
        }, () => {
            LinkBtn.innerHTML = "Couldn't copy D:"
            setTimeout(() => {
                LinkBtn.innerHTML = "Copy link"
            }, 2000)
        })
    }
}