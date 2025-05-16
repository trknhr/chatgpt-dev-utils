import { WebSocketServer } from "ws";
import http from "http";
// import { logger } from "shared/logger";

const PORT = 32123;
let wsClient: any = null;

const server = http.createServer((req, res) => {
    if (req.method === "GET" && req.url === "/ping") {
        res.writeHead(200).end("pong");
        return;
    }
    if (req.method === "POST" && req.url === "/chatgpt-prompt") {
        let body = "";
        req.on("data", chunk => (body += chunk));
        req.on("end", () => {
            try {
                const { prompt } = JSON.parse(body);
                console.log("📩 Received prompt:", prompt);

                if (wsClient && wsClient.readyState === wsClient.OPEN) {
                    wsClient.send(JSON.stringify({ type: "chatgpt-prompt", prompt }));
                    res.writeHead(200).end("Sent to extension");
                } else {
                    console.warn("❗ No WebSocket connection to extension");
                    res.writeHead(500).end("Extension not connected");
                }
            } catch (err) {
                res.writeHead(500).end("Invalid JSON");
            }
        });
    } else {
        res.writeHead(404).end("Not found");
    }
});

const wss = new WebSocketServer({ server });

wss.on("connection", ws => {
  console.log("🔌 Extension connected via WebSocket");
  wsClient = ws;

  ws.on("message", message => {
    if (message.toString() === "ping") {
      ws.send("pong");
    }
  });

  ws.on("close", () => {
    console.log("❎ Extension WebSocket disconnected");
    wsClient = null;
  });
});

server.listen(PORT, () => {
    console.log(`🚀 ChatGPT Proxy (WS+HTTP) running at http://localhost:${PORT}`);
});
