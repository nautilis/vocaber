const host = "http://localhost:5000";

function appendList(value, id, times) {
    let review = document.getElementById("review");
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
    list.appendChild(span);
    list.appendChild(button);
    list.appendChild(span2);
    list.appendChild(span1);
    review.appendChild(list);
}
function handleButtonClick(element) {
    let buttonId = element.id;
    let id = buttonId.split("_")[1];
    console.log(id);
    knownIt(id);
};

function knownIt(id){
    let xhr = new XMLHttpRequest();
    xhr.open("POST", host +"/known_it", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.onreadystatechange = () => {
        if(xhr.readyState == 4){
            let resp = JSON.parse(xhr.responseText);
            let result = resp.result;
            console.log(resp);
            if(result == "failed"){
                alert("update failed");
            }else{
                let timespan = document.getElementById("span_" + id);
                timespan.innerHTML =parseInt(timespan.innerHTML) + 1;
            }
            
        }
    }
    xhr.send("itemid="+id +"&token=19959995");
}

(function () {
    let xhr = new XMLHttpRequest();
    xhr.open("GET", host + "/get_not_master", true);
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4) {
            let resp = JSON.parse(xhr.responseText);
            let items = resp.items;
            console.log(items);
            for (let i = 0; i < items.length; i++) {
                appendList(items[i].value, items[i].id, items[i].knownit);
            }
            let known_btn = document.getElementsByClassName("known_btn");

            for(let i=0;i< known_btn.length; i++){
                known_btn[i].onclick = function(){
                    handleButtonClick(this);
                } 
            }
        }
    }
    xhr.send();
})();

