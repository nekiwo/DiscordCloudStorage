const upload: Function = () => {
    const file: File = document.getElementById("input").files[0];
    const form: FormData = new FormData();
    form.append("file", file);

    fetch("/upload", {
        method: "POST",
        body: form
    }).then(
        response: Response => response.json()
    ).then((data) => {
        console.log(data);
    }).catch((err: Error)) => {
        console.log(err)
        alert("Error: " + err)
    })
};