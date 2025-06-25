let ws = null;

function connectWebSocket() {
  if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
    return; // do nothing when it's connected already
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

chrome.runtime.onInstalled.addListener(() => {
  chrome.alarms.create('heartbeat', { periodInMinutes: 1 });
});

chrome.alarms.onAlarm.addListener((alarm) => {
  if (alarm.name === 'heartbeat') {
    console.log("I'm alive 🫀");
  }
});

function openOrCreateChatGPTTab(prompt) {
  // If a ChatGPT tab exists, reuse it; otherwise, create a new one
  chrome.tabs.query({}, (tabs) => {
    const existingNewChatPage = tabs.find(tab =>
      tab.url && tab.url === "https://chatgpt.com" && tab.status === "complete"
    );

    if (existingNewChatPage) {
      console.log("🟢 Found existing ChatGPT tab:", existingNewChatPage.id);
      chrome.tabs.sendMessage(existingNewChatPage.id, { type: "chatgpt-prompt", prompt });
    } else {
      // 新しいタブを作成して、読み込み完了を待つ
      chrome.tabs.create({ url: "https://chatgpt.com" }, (tab) => {
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
