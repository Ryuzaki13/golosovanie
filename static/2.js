console.log("obj", obj);

let sum = obj.votes.map(x => x.votes).reduce((a, b) => a + b, 0)

let html = ""

if (obj.winner == 0) {
    html += `Голосование идёт<br>`
} else {
    html += `Голосование окончено<br>`
}
html += `Дата окончания: ${obj.date}<br>`
for (let i = 0; i < obj.votes.length; i++) {
    html += `
    <input type="radio" id="input${i}" value="${obj.votes[i].song_id}">
    <label for="input${i}">
        <a href="${obj.votes[i].url}">${obj.votes[i].name}</a>(голосов: ${obj.votes[i].votes}, ${Math.round(obj.votes[i].votes / sum * 100)}%)<br>
    </label>
    `
}

if (obj.winner == 0) {
    html += `<button id="g">Голосовать</button><br>`
}

if (obj.user_student == false){
    html += `<a href="/edit-songs">Изменить список песен</a>`
}

document.body.innerHTML = html

try { document.querySelector(`input[value='${obj.user_choice}']`).checked = true } catch (error) { }

document.querySelector(`#g`).addEventListener("click", () => {
    if (![...document.querySelectorAll("input")].filter(x => x.checked)[0]) {
        fetch('/choice', {
            method: 'POST',
            body: [...document.querySelectorAll("input")].filter(x => x.checked)[0].value
        })
    }
})


