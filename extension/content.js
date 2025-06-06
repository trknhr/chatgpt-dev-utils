function findInputBox() {
  return document.querySelector('div[contenteditable="true"]');
}

console.log("content.js is loadded")
chrome.runtime.onMessage.addListener((message) => {
  if (message.type === "chatgpt-prompt" && message.prompt) {
    console.log("🧠 content.js received prompt:", message.prompt);

    const inputBox = findInputBox();
    if (!inputBox) {
      console.warn("❗ ChatGPT input box not found.");
      return;
    }

    inputBox.focus();
    const html = message.prompt
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/\n/g, "<br>");
    console.log(html)
    inputBox.innerHTML = html

    // Important: fire `input` event
    inputBox.dispatchEvent(new InputEvent("input", { bubbles: true }));

    // After wainting 100ms, Fire Enter key event
    setTimeout(() => {
      inputBox.dispatchEvent(new KeyboardEvent("keydown", {
        bubbles: true,
        cancelable: true,
        key: "Enter",
        code: "Enter",
        keyCode: 13,
        which: 13
      }));
    }, 1000);
  }
});
