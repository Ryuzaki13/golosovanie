async function main() {

    let data = await fetch("./api/user/select-info", { "method": "POST" }).then(x => x.json())

    if (data.Error) {
        console.log(data);
    }

    let id = data.Data.Phone


}
main()
