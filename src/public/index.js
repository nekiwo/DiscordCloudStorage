const upload = () => {
    const file = document.getElementById("input").files[0];
    const FormData = new FormData();
    FormData.append("file", file);

    fetch("/upload", {
        method: "POST",
        body: FormData
    }).then(
        response => response.json()
    ).then(
        data => {
            console.log(data);
        }
    ).catch(error => console.log(error));
};