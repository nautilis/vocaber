const host = "http://127.0.0.1:5000";
const token = "19959995";

function appendList(value, id, times) {
    let todayItems = document.getElementById("today-items");
    let list = document.createElement("li");
    let button = document.createElement("button");
    button.setAttribute("class", "known_btn");
    button.setAttribute("id", "button_" + id);
    button.appendChild(document.createTextNode("knew it"));
    let span1 = document.createElement("span");
    span1.appendChild(document.createTextNode(" times "));
    let span2 = document.createElement("span");
    span2.setAttribute("id", "span_" + id);
    span2.appendChild(document.createTextNode(times));
    let span = document.createElement("span");
    span.appendChild(document.createTextNode(value));
    span.setAttribute("class", "word");
    span.setAttribute("id", "value_" +id);
    let deleteButton = document.createElement("button");
    deleteButton.appendChild(document.createTextNode("delete"))
    deleteButton.setAttribute("id", "del_btn_" + id);
    deleteButton.setAttribute("class", "del_btn");
    let transButton = document.createElement("button");
    transButton.appendChild(document.createTextNode("translate"))
    transButton.setAttribute("id", "trans_btn_" + id);
    transButton.setAttribute("class", "trans_btn");

    list.appendChild(span);
    list.appendChild(button);
    list.appendChild(span2);
    list.appendChild(span1);
    list.appendChild(transButton);
    list.appendChild(deleteButton);
    list.setAttribute("id", "list_" + id);

    todayItems.appendChild(list);
}

function handleButtonClick(element) {
    let buttonId = element.id;
    let id = buttonId.split("_")[1];
    console.log(id);
    knownIt(id);
    element.disabled = true;
};

function handleTranslate(element){
    element.disabled = true;
    let buttonId = element.id;
    let id = buttonId.split("_")[2];
    console.log("translate...   " + id + "translate....");
    let valueSpan = document.getElementById("value_" + id);
    let value = valueSpan.innerText;
    console.log(value)

    translate(element, id, value);
}

function translate(element, id, value){
    let xhr = new XMLHttpRequest();
    xhr.open("GET", host +"/get_translate/" + value, true);
    xhr.onreadystatechange = function(){
        if(xhr.readyState == 4){
            let resp = JSON.parse(xhr.responseText);
            let list = document.getElementById("list_" + id);
            let span = document.createElement("span");
            span.appendChild(document.createTextNode(resp.result));
            span.setAttribute("class", "trans_result");
            list.appendChild(span);
            element.disabled = false;
        }
    }
    xhr.send();
}




function knownIt(id) {
    let xhr = new XMLHttpRequest();
    xhr.open("POST", host + "/known_it", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.onreadystatechange = () => {
        if (xhr.readyState == 4) {
            let resp = JSON.parse(xhr.responseText);
            let result = resp.result;
            console.log(resp);
            if (result == "failed") {
                alert("update failed");
            } else {
                let timespan = document.getElementById("span_" + id);
                timespan.innerHTML = parseInt(timespan.innerHTML) + 1;
            }

        }
    }
    xhr.send("itemid=" + id + "&token=" + token);
}

function handleDelClick(element) {
    let r = confirm("Are you sure to delete it?");
    if (r == true) {
        let btnId = element.id;
        let id = btnId.split("_")[2]
        delete_item(id);
    }
}

function delete_item(id) {
    let xhr = new XMLHttpRequest();
    xhr.open("POST", host + "/delete_item", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.onreadystatechange = () => {
        if (xhr.readyState == 4) {
            let resp = JSON.parse(xhr.responseText);
            let result = resp.result;
            console.log(resp);
            if (result == "failed") {
                alert("delete failed");
            } else {
                let list = document.getElementById("list_" + id);
                list.parentElement.removeChild(list);
            }
        }
    }
    xhr.send("itemid=" + id + "&token=" + token);
}


(function () {
    let xhr = new XMLHttpRequest();
    xhr.open("GET", host + "/items_by_subday?subday=0", true);
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4) {
            let resp = JSON.parse(xhr.responseText);
            let items = resp.items;
            console.log(items);
            for (let i = 0; i < items.length; i++) {
                appendList(items[i].value, items[i].id, items[i].knownit);
            }
            let known_btn = document.getElementsByClassName("known_btn");

            for (let i = 0; i < known_btn.length; i++) {
                known_btn[i].onclick = function () {
                    handleButtonClick(this);
                }
            }

            let del_btn = document.getElementsByClassName("del_btn");

            for (let i = 0; i < del_btn.length; i++) {
                del_btn[i].onclick = function () {
                    handleDelClick(this);
                }
            }

            let trans_btn = document.getElementsByClassName("trans_btn");

            for(let i=0;i< trans_btn.length;i++){
                trans_btn[i].onclick = function(){
                    handleTranslate(this);
                }
            }


        }
    }
    xhr.send();
})();