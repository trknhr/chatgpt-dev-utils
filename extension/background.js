let ws = null;

const LOG_LEVEL = __LOG_LEVEL__; 

const LOG_LEVELS = {
  none: 0,
  error: 1,
  warn: 2,
  log: 3
};

function shouldLog(type) {
  return LOG_LEVELS[type] <= LOG_LEVELS[LOG_LEVEL];
}

// Utility function to add timestamps to logs
function logWithTimestamp(message, type = 'log') {
  if (!shouldLog(type)) return;

  const timestamp = new Date().toISOString();
  const formattedMessage = `[${timestamp}] ${message}`;
  
  switch (type) {
    case 'error':
      console.error(formattedMessage);
      break;
    case 'warn':
      console.warn(formattedMessage);
      break;
    default:
      console.log(formattedMessage);
  }
}

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
    logWithTimestamp("✅ WebSocket connected to CLI proxy");
  });

  ws.addEventListener("message", (event) => {
    logWithTimestamp("📬 Message received from CLI: " + event.data);
    try {
      const { type, prompt } = JSON.parse(event.data);
      if (type === "chatgpt-prompt") {
        logWithTimestamp("📨 Prompt received from CLI: " + prompt);
        openOrCreateChatGPTTab(prompt);
      }
    } catch (e) {
      logWithTimestamp("❌ Invalid WS message: " + e, 'error');
    }
  });

  ws.addEventListener("close", () => {
    logWithTimestamp("🔌 WebSocket disconnected", 'warn');
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
        logWithTimestamp("🔁 Attempting to reconnect WebSocket...");
        connectWebSocket();
      } else {
        logWithTimestamp("🚫 CLI proxy is not running");
      }
    });
  } else if (ws.readyState === WebSocket.OPEN) {
    ws.send("ping");
    logWithTimestamp("📡 Sent ping to CLI proxy (keep-alive)");
  }
}, 1000);

// Create a heartbeat alarm to verify the extension is alive
chrome.runtime.onInstalled.addListener(() => {
  chrome.alarms.create('heartbeat', { periodInMinutes: 1 });
});

chrome.alarms.onAlarm.addListener((alarm) => {
  if (alarm.name === 'heartbeat') {
    logWithTimestamp("I'm alive 🫀");
  }
});

// Open or reuse a ChatGPT tab and send the prompt
function openOrCreateChatGPTTab(prompt) {
  chrome.tabs.query({}, (tabs) => {
    const existingNewChatPage = tabs.find(tab =>
      tab.url && tab.url === "https://chatgpt.com" && tab.status === "complete"
    );

    if (existingNewChatPage) {
      logWithTimestamp("🟢 Found existing ChatGPT tab: " + existingNewChatPage.id);
      chrome.tabs.sendMessage(existingNewChatPage.id, { type: "chatgpt-prompt", prompt });
    } else {
      chrome.tabs.create({ url: "https://chatgpt.com" }, (tab) => {
        const tabId = tab.id;
        logWithTimestamp("🆕 Created new ChatGPT tab: " + tabId);

        const checkTabReady = (retries = 20) => {
          if (retries <= 0) {
            logWithTimestamp("⚠️ New ChatGPT tab did not load in time", 'warn');
            return;
          }

          chrome.tabs.get(tabId, (updatedTab) => {
            if (updatedTab.status === "complete") {
              logWithTimestamp("✅ ChatGPT tab is ready: " + updatedTab.id);
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
