const DLTitle = document.getElementById("dltitle");
const download = () => {
    DLTitle.innerHTML = "Searching for your file...";
    fetch("/download", {
        method: "POST",
        headers: {
            "Content-Type": "text/plain"
        },
        body: new URLSearchParams(window.location.search).get("f")
    }).then((data) => {
        // Get text from promise
        data.text().then((url) => {
            if (url != "null") {
                // Status change
                DLTitle.innerHTML = "File downloading!";
                // Download file
                const link = document.createElement("a");
                link.setAttribute("download", "");
                link.href = url;
                document.body.appendChild(link);
                link.click();
                link.remove();
            }
            else {
                // File not found
                DLTitle.innerHTML = "File not found :/";
            }
        });
    }).catch((err) => {
        console.error(err);
        DLTitle.innerHTML = "Error: " + err;
    });
};
