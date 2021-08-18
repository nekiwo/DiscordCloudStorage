let DownloadURL: string = "null"

const upload: Function = () => {
    const input: HTMLInputElement = <HTMLInputElement>document.getElementById("input")
    const file: File = input.files[0];
    const form: FormData = new FormData();
    form.append("file", file);

    fetch("/upload", {
        method: "POST",
        body: form
    }).then((data: Response) => {
        data.text().then((url: string) => {
            DownloadURL = url
        })
    }).catch((err: Error) => {
        console.log(err)
        alert("Error: " + err)
    })
}

const CopyLink: Function = () => {
    const LinkBtn: HTMLElement = document.getElementById("linkbtn")

    if (DownloadURL != "null") {
        navigator.clipboard.writeText(DownloadURL).then(() => {
            LinkBtn.innerHTML = "Copied!"
            setTimeout(() => {
                LinkBtn.innerHTML = "Copy"
            }, 2000)
        }, () => {
            LinkBtn.innerHTML = "Couldn't copy D:"
            setTimeout(() => {
                LinkBtn.innerHTML = "Copy"
            }, 2000)
        })
    }
}