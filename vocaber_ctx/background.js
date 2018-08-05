chrome.runtime.onInstalled.addListener(function () {
    chrome.storage.sync.set({ color: "#3aa757" }, function () {
        console.log("The color is green");
    });
    chrome.contextMenus.create({
        "id": "Vocaber",
        "title": "vocaber",
    })

});
chrome.commands.onCommand.addListener(command => {
    if (command === 'toggle-vober') {
        console.log("toggle");
        chrome.browserAction.setPopup({"popup": "popup.html"});
    }
});
