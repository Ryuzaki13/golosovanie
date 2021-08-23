(async () => {
    let data = await fetch("./api/user/select-info", { "method": "POST" }).then(x => x.json())

    if (data.Error) {
        // document.body.innerText = "Пользователь не авторизирован"
    } else {
        if (data.Data.Branch == "44b72cd889a44234a8a7a49750fedf1bc5a654d8b06257ec5866711be0be6286") {
            await fetch("./current", { "method": "POST", "body": data.Data.Phone }).then(x => x.text()).then(x => document.body.innerHTML = x)
        } else {
            // document.body.innerText = "Только для учащихся АиВТ"
        }
    }
})()
