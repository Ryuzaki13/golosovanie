(async () => {
    let data = await fetch("./api/user/select-info", { "method": "POST" }).then(x => x.json())
    if (data.Error) {
        // document.body.innerText = "Пользователь не авторизирован"
    } else {
        if (data.Data.Branch == "44b72cd889a44234a8a7a49750fedf1bc5a654d8b06257ec5866711be0be6286" || data.Data.isStudent == false) {
            await fetch("./vote", { "method": "POST", "body": data.Data.Phone + "-" + data.Data.isStudent }).then(x => x.text()).then(x => {
                document.body.innerHTML = x;
                [...document.body.querySelectorAll("script")].forEach(oldScript => {
                    let newScript = document.createElement("script");
                    [...oldScript.attributes].forEach(attr => newScript.setAttribute(attr.name, attr.value));
                    newScript.appendChild(document.createTextNode(oldScript.innerHTML));
                    oldScript.parentNode.replaceChild(newScript, oldScript);
                });
            })
        } else {
            // document.body.innerText = "Только для студентов АиВТ"
        }
    }
})()
