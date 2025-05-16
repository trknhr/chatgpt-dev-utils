let ws = null;
let reconnectAttempts = 0;

function connectWebSocket() {
  ws = new WebSocket("ws://localhost:32123");

  ws.addEventListener("open", () => {
    console.log("✅ WebSocket connected to CLI proxy");
    reconnectAttempts = 0;
  });

  ws.addEventListener("message", (event) => {
    try {
      const { type, prompt } = JSON.parse(event.data);
      if (type === "chatgpt-prompt") {
        console.log("📨 Prompt received from CLI:", prompt);
        waitForChatGPTTab().then((tab) => {
          console.log("active tabID is ", tab.id)
          chrome.tabs.sendMessage(tab.id, { type: "chatgpt-prompt", prompt });
        });
      }
    } catch (e) {
      console.error("❌ Invalid WS message:", e);
    }
  });

  ws.addEventListener("close", () => {
    console.warn("🔌 WebSocket disconnected");
    attemptReconnect();
  });

  ws.addEventListener("error", (err) => {
    console.error("❌ WebSocket error", err);
    ws.close();
  });
}

function attemptReconnect() {
  reconnectAttempts++;
  const timeout = Math.min(5000, 500 * reconnectAttempts);
  setTimeout(() => {
    console.log("🔁 Reconnecting WebSocket...");
    connectWebSocket();
  }, timeout);
}

connectWebSocket();

setInterval(() => {
  if (ws?.readyState === WebSocket.OPEN) {
    ws.send("ping");
    console.log("📡 Sent ping to CLI proxy (keep-alive)");
  }
}, 30000);

// Reuse waitForChatGPTTab() as before
function waitForChatGPTTab(maxRetries = 20, interval = 500) {
  return new Promise((resolve, reject) => {
    const check = (remaining) => {
      chrome.tabs.query({}, (tabs) => {
        const chatTab = tabs.find(tab =>{
            return tab.url && tab.url === "https://chatgpt.com/" && tab.status === "complete"
        })
        if (chatTab) return resolve(chatTab);
        if (remaining <= 0) return reject("ChatGPT tab not found");
        setTimeout(() => check(remaining - 1), interval);
      });
    };
    check(maxRetries);
  });
}
