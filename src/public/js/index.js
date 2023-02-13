let DownloadURL = "null";
const LinkBtn = document.getElementById("linkbtn");
const upload = () => {
    const input = document.getElementById("input");
    const file = input.files[0];
    if (input.value != "") {
        const form = new FormData();
        form.append("file", file);
        document.getElementById("upload").innerHTML = "Uploading...";
        fetch("/upload", {
            method: "POST",
            body: form
        }).then((data) => {
            console.log(data);
            data.text().then((url) => {
                console.log(url);
                DownloadURL = url;
                document.getElementById("upload").innerHTML = "Uploaded!";
                document.getElementById("uploadbtn").style.display = "none";
                LinkBtn.style.display = "block";
            });
        }).catch((err) => {
            console.log(err);
            alert("Error: " + err);
        });
    }
};
const CopyLink = () => {
    if (DownloadURL != "null") {
        navigator.clipboard.writeText(document.location.origin + DownloadURL).then(() => {
            LinkBtn.innerHTML = "Link copied!";
            setTimeout(() => {
                LinkBtn.innerHTML = "Copy link";
            }, 2000);
        }, () => {
            LinkBtn.innerHTML = "Couldn't copy D:";
            setTimeout(() => {
                LinkBtn.innerHTML = "Copy link";
            }, 2000);
        });
    }
};
