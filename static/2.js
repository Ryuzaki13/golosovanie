console.log("obj", obj);

let sum = obj.votes.map(x => x.votes).reduce((a, b) => a + b, 0)

let html = ""
html += `<form action="/choice" method="post">`
if (obj.winner == 0) {
    html += `Голосование идёт<br>`
} else {
    html += `Голосование окончено<br>`
}
html += `Дата окончания: ${obj.date}<br>`
for (let i = 0; i < obj.votes.length; i++) {
    html += `
    <input type="radio" id="input${i}" name="vote" value="${obj.votes[i].song_id}">
    <label for="input${i}">
        <a href="${obj.votes[i].url}">${obj.votes[i].name}</a>
    </label>(голосов: ${obj.votes[i].votes}, ${Math.round(obj.votes[i].votes / sum * 100)}%)<br>
    `
}

html += `<button type="submit">Голосовать</button>`

html += `</form>`

html += `<a href="/history">История голосований</a>`

document.body.innerHTML = html

try { document.querySelector(`input[value='${obj.user_choice}']`).checked = true } catch (error) { }

document.querySelector(`form`).addEventListener("submit", (e) => {
    e.preventDefault()
    console.log(e);
    let form_data = new FormData(e.target)
    console.log(form_data.get("vote"));

    // let response = await fetch('/article/formdata/post/user', {
    //     method: 'POST',
    //     body: new FormData(formElem)
    // });
    // let result = await response.json();
})

