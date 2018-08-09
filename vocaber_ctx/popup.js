const host = "http://127.0.0.1:5000";
const token = "19959995";

let newVocab = document.getElementById('newVocab');
newVocab.focus();
newVocab.onchange = function (element) {
    console.log(newVocab.value);
    if (newVocab.value.length != 0) {
        waiting();
        hello(newVocab.value);
        translate(newVocab.value);
    }
};
function waiting(){
    let translate = document.getElementById("translation");
    translate.setAttribute("style", "display: inline-block");
    translate.innerHTML = "waiting translate..."; 
}

function translate(value){
    var xhr = new XMLHttpRequest();
    xhr.open("GET", "http://47.90.206.255:8000/get_translate/" + value, true);
    xhr.onreadystatechange = function(){
        if(xhr.readyState == 4){
            let resp = JSON.parse(xhr.responseText);
            let translate = document.getElementById("translation");
            translate.setAttribute("style", "display: inline-block");
            translate.innerHTML = resp.result;
        }
    }
    xhr.send();
}

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
    xhr.send("item=" + value + "&token=" + token);
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