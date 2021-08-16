const download = () => {
    fetch("/download", {
        method: "POST",
        headers: {
            "Content-Type": "text/plain"
        },
        body: new URLSearchParams(window.location.search).get("f")
    }).then((data) => {
        console.log(data)
        // Get text from promise
        data.text().then((url) => {
            // Download file
            const link = document.createElement("a");
            link.setAttribute("download", "");
            link.href = url;
            document.body.appendChild(link);
            link.click();
            link.remove();
        })
    }).catch(error => console.log(error));
}