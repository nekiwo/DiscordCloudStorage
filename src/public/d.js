const dltitle = document.getElementById("dltitle")

const download = () => {
    dltitle.innerHTML = "Searching for your file..."

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
            if (url != "null") {
                // Download file
                dltitle.innerHTML = "File downloading!";

                const link = document.createElement("a");
                link.setAttribute("download", "");
                link.href = url;
                document.body.appendChild(link);
                link.click();
                link.remove();
            } else {
                // File not found
                dltitle.innerHTML = "File not found :/";
            }
            
        })
    }).catch(error => console.log(error));
}