const host = "http://127.0.0.1:5000";

let newVocab = document.getElementById('newVocab');
newVocab.onchange = function (element) {
    console.log(newVocab.value);
    if (newVocab.value.length != 0) {
        hello(newVocab.value);
    }
};

function hello(value) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", host + "/item", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4) {
            // JSON.parse does not evaluate the attacker's scripts.
            let resp = JSON.parse(xhr.responseText);
            let count = document.getElementById('count');
            count.innerHTML = resp.result;
        }
    }
    xhr.send("item=" + value + "&token=19959995");
}

let review = document.getElementById("review");
review.onclick = function () {
    window.open("review.html");
};

let seeall = document.getElementById("seeall");
seeall.onclick = () => {
    window.open("see-all.html");
}


(function () {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", host + "/get_today_count", true);
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4) {
            let resp = JSON.parse(xhr.responseText);
            let count = document.getElementById('count');
            count.innerHTML = resp.result;
        }
    }
    xhr.send();
})();