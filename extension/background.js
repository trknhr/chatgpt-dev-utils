// let ws = null;
// let reconnectAttempts = 0;

// function connectWebSocket() {
//   ws = new WebSocket("ws://localhost:32123");

//   ws.addEventListener("open", () => {
//     console.log("✅ WebSocket connected to CLI proxy");
//     reconnectAttempts = 0;
//   });

//   ws.addEventListener("message", (event) => {
//     try {
//       const { type, prompt } = JSON.parse(event.data);
//       if (type === "chatgpt-prompt") {
//         console.log("📨 Prompt received from CLI:", prompt);
//         waitForChatGPTTab().then((tab) => {
//           console.log("active tabID is ", tab.id)
//           chrome.tabs.sendMessage(tab.id, { type: "chatgpt-prompt", prompt });
//         });
//       }
//     } catch (e) {
//       console.error("❌ Invalid WS message:", e);
//     }
//   });

//   ws.addEventListener("close", () => {
//     console.warn("🔌 WebSocket disconnected");
//     attemptReconnect();
//   });

//   ws.addEventListener("error", (err) => {
//     console.error("❌ WebSocket error", err);
//     ws.close();
//   });
// }

// function attemptReconnect() {
//   reconnectAttempts++;
//   const timeout = Math.min(5000, 500 * reconnectAttempts);
//   setTimeout(() => {
//     console.log("🔁 Reconnecting WebSocket...");
//     connectWebSocket();
//   }, timeout);
// }

// connectWebSocket();

// setInterval(() => {
//   if (ws?.readyState === WebSocket.OPEN) {
//     ws.send("ping");
//     console.log("📡 Sent ping to CLI proxy (keep-alive)");
//   }
// }, 30000);

// // Reuse waitForChatGPTTab() as before
// function waitForChatGPTTab(maxRetries = 20, interval = 500) {
//   return new Promise((resolve, reject) => {
//     const check = (remaining) => {
//       chrome.tabs.query({}, (tabs) => {
//         const chatTab = tabs.find(tab =>{
//             return tab.url && tab.url === "https://chatgpt.com/" && tab.status === "complete"
//         })
//         if (chatTab) return resolve(chatTab);
//         if (remaining <= 0) return reject("ChatGPT tab not found");
//         setTimeout(() => check(remaining - 1), interval);
//       });
//     };
//     check(maxRetries);
//   });
// }

// chrome.runtime.onInstalled.addListener(() => {
//   chrome.alarms.create('heartbeat', { periodInMinutes: 1 });
// });

// chrome.alarms.onAlarm.addListener((alarm) => {
//   if (alarm.name === 'heartbeat') {
//     console.log("I'm alive 🫀");
//   }
// });

let ws = null;

function connectWebSocket() {
  if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
    return; // すでに接続中または接続済みなら何もしない
  }

  ws = new WebSocket("ws://localhost:32123/ws");

  ws.addEventListener("open", () => {
    console.log("✅ WebSocket connected to CLI proxy");
  });

  ws.addEventListener("message", (event) => {
    console.log("📬 Message received from CLI:", event.data);
    try {
      const { type, prompt } = JSON.parse(event.data);
      if (type === "chatgpt-prompt") {
        console.log("📨 Prompt received from CLI:", prompt);
        openOrCreateChatGPTTab(prompt);
        // waitForChatGPTTab().then((tab) => {
        //   console.log("active tabID is ", tab.id);
        //   chrome.tabs.sendMessage(tab.id, { type: "chatgpt-prompt", prompt });
        // });
      }
    } catch (e) {
      console.error("❌ Invalid WS message:", e);
    }
  });

  ws.addEventListener("close", () => {
    console.warn("🔌 WebSocket disconnected");
    ws = null;
  });

  ws.addEventListener("error", (err) => {
    ws.close();
    ws = null;
  });
}

// 1秒ごとに接続を確認して未接続なら再接続
setInterval(() => {
  if (!ws || ws.readyState === WebSocket.CLOSED) {
    console.log("🔁 Attempting to reconnect WebSocket...");
    connectWebSocket();
  } else if (ws.readyState === WebSocket.OPEN) {
    ws.send("ping");
    console.log("📡 Sent ping to CLI proxy (keep-alive)");
  }
}, 1000);

// Reuse waitForChatGPTTab() as before
function waitForChatGPTTab(maxRetries = 20, interval = 500) {
  return new Promise((resolve, reject) => {
    const check = (remaining) => {
      chrome.tabs.query({}, (tabs) => {
        const chatTab = tabs.find(tab =>
          tab.url && tab.url === "https://chatgpt.com/" && tab.status === "complete"
        );
        if (chatTab) return resolve(chatTab);
        if (remaining <= 0) return reject("ChatGPT tab not found");
        setTimeout(() => check(remaining - 1), interval);
      });
    };
    check(maxRetries);
  });
}

chrome.runtime.onInstalled.addListener(() => {
  chrome.alarms.create('heartbeat', { periodInMinutes: 1 });
});

chrome.alarms.onAlarm.addListener((alarm) => {
  if (alarm.name === 'heartbeat') {
    console.log("I'm alive 🫀");
  }
});

function openOrCreateChatGPTTab(prompt) {
  // 既存タブがあれば使い、なければ新規作成
  chrome.tabs.query({}, (tabs) => {
    const existingTab = tabs.find(tab =>
      tab.url && tab.url.startsWith("https://chatgpt.com/") && tab.status === "complete"
    );

    if (existingTab) {
      console.log("🟢 Found existing ChatGPT tab:", existingTab.id);
      chrome.tabs.sendMessage(existingTab.id, { type: "chatgpt-prompt", prompt });
    } else {
      // 新しいタブを作成して、読み込み完了を待つ
      chrome.tabs.create({ url: "https://chatgpt.com/" }, (tab) => {
        const tabId = tab.id;
        console.log("🆕 Created new ChatGPT tab:", tabId);

        // ポーリングして読み込み完了を待つ
        const checkTabReady = (retries = 20) => {
          if (retries <= 0) {
            console.warn("⚠️ New ChatGPT tab did not load in time");
            return;
          }

          chrome.tabs.get(tabId, (updatedTab) => {
            if (updatedTab.status === "complete") {
              console.log("✅ ChatGPT tab is ready:", updatedTab.id);
              chrome.tabs.sendMessage(updatedTab.id, { type: "chatgpt-prompt", prompt });
            } else {
              setTimeout(() => checkTabReady(retries - 1), 500);
            }
          });
        };

        checkTabReady();
      });
    }
  });
}
