const download = () => {
    fetch("/download", {
        method: "POST",
        headers: {
            "Content-Type": "text/plain"
        },
        body: new URLSearchParams(window.location.search).get("f")
    }).then((data) => {

    }).catch(error => console.log(error));
};