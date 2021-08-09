const upload = () => {
    const file = document.getElementById("input").files[0];
    const form = new FormData();
    form.append("file", file);

    fetch("/upload", {
        method: "POST",
        body: form
    }).then(
        response => response.json()
    ).then(
        data => {
            console.log(data);
        }
    ).catch(error => console.log(error));
};