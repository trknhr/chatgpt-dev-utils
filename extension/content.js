function findInputBox() {
  return document.querySelector('div[contenteditable="true"]');
}

function textToParagraphs(text) {
  // 改行で split → 各行を <p> に包む → 連結
  return text
    .split(/\n/)
    .map(line => `<p>${escapeHtml(line)}</p>`)
    .join('');
}

function escapeHtml(str) {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
}



console.log("content.js is loadded")
chrome.runtime.onMessage.addListener((message) => {
  if (message.type === "chatgpt-prompt" && message.prompt) {
    console.log("🧠 content.js received prompt:", message.prompt);

    waitForInputBox().then((inputBox) => {
      setTimeout(() => {
      inputBox.focus();
      const html = message.prompt
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/\n/g, "<br>");
        inputBox.innerHTML = textToParagraphs(message.prompt)
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
      })
      // Important: fire `input` event

      // After wainting 100ms, Fire Enter key event
      // setTimeout(() => {
      //   inputBox.dispatchEvent(new KeyboardEvent("keydown", {
      //     bubbles: true,
      //     cancelable: true,
      //     key: "Enter",
      //     code: "Enter",
      //     keyCode: 13,
      //     which: 13
      //   }));
      // }, 1000);
    })
  }
});

function waitForInputBox(retries = 10, delay = 500) {
  return new Promise((resolve, reject) => {
    const tryFind = () => {
      const inputBox = findInputBox();
      if (inputBox) {
        resolve(inputBox);
      } else if (retries > 0) {
        setTimeout(() => waitForInputBox(retries - 1, delay).then(resolve).catch(reject), delay);
      } else {
        reject(new Error("ChatGPT input box not found"));
      }
    };
    tryFind();
  });
}

async function insertTextWithNewlines(inputBox, text) {
  inputBox.focus();

  const lines = text.split('\n');
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];

    // insert text
    inputBox.dispatchEvent(new InputEvent("beforeinput", {
      bubbles: true,
      cancelable: true,
      inputType: "insertText",
      data: line,
    }));

    // insert line break except after the last line
    if (i < lines.length - 1) {
      inputBox.dispatchEvent(new KeyboardEvent("keydown", {
        bubbles: true,
        cancelable: true,
        key: "Enter",
        code: "Enter",
        keyCode: 13,
        which: 13,
      }));

      // 一部環境では keypress も必要な場合あり
      inputBox.dispatchEvent(new KeyboardEvent("keypress", {
        bubbles: true,
        cancelable: true,
        key: "Enter",
        code: "Enter",
        keyCode: 13,
        which: 13,
      }));

      inputBox.dispatchEvent(new KeyboardEvent("keyup", {
        bubbles: true,
        cancelable: true,
        key: "Enter",
        code: "Enter",
        keyCode: 13,
        which: 13,
      }));
    }
  }

  // 念のため最終的に input イベントを送って状態を確定
  inputBox.dispatchEvent(new Event("input", { bubbles: true }));
}
