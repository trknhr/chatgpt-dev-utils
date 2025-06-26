let ws = null;

// Check if the local CLI proxy is available before attempting to connect
function isServerRunning() {
  return fetch("http://localhost:32123/ping", { method: "GET" })
    .then(() => true)
    .catch(() => false);
}

function connectWebSocket() {
  if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
    return; // do nothing when it's already connected or connecting
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

  ws.addEventListener("error", () => {
    // Silently handle connection errors without logging
    ws = null;
  });
}

// Check WebSocket connection every second
setInterval(() => {
  if (!ws || ws.readyState === WebSocket.CLOSED) {
    isServerRunning().then((running) => {
      if (running) {
        console.log("🔁 Attempting to reconnect WebSocket...");
        connectWebSocket();
      } else {
        console.log("🚫 CLI proxy is not running");
      }
    });
  } else if (ws.readyState === WebSocket.OPEN) {
    ws.send("ping");
    console.log("📡 Sent ping to CLI proxy (keep-alive)");
  }
}, 1000);

// Create a heartbeat alarm to verify the extension is alive
chrome.runtime.onInstalled.addListener(() => {
  chrome.alarms.create('heartbeat', { periodInMinutes: 1 });
});

chrome.alarms.onAlarm.addListener((alarm) => {
  if (alarm.name === 'heartbeat') {
    console.log("I'm alive 🫀");
  }
});

// Open or reuse a ChatGPT tab and send the prompt
function openOrCreateChatGPTTab(prompt) {
  chrome.tabs.query({}, (tabs) => {
    const existingNewChatPage = tabs.find(tab =>
      tab.url && tab.url === "https://chatgpt.com" && tab.status === "complete"
    );

    if (existingNewChatPage) {
      console.log("🟢 Found existing ChatGPT tab:", existingNewChatPage.id);
      chrome.tabs.sendMessage(existingNewChatPage.id, { type: "chatgpt-prompt", prompt });
    } else {
      chrome.tabs.create({ url: "https://chatgpt.com" }, (tab) => {
        const tabId = tab.id;
        console.log("🆕 Created new ChatGPT tab:", tabId);

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
