import fetch from "node-fetch";
import { spawn } from "child_process";
import path from "path";
import { fileURLToPath } from "url";
import { logger } from "./logger";

const __dirname = path.dirname(fileURLToPath(import.meta.url));

async function ensureProxyRunning(): Promise<void> {
  const url = "http://localhost:32123/ping";

  try {
    const res = await fetch(url);
    if (res.ok) return; 
  } catch(e) {
    // Not running, will start
    logger.error("failed to launch server", e)
  }

  const proxyScript = path.resolve(__dirname, "../../extension-proxy/server.ts");
  const tsxPath = path.resolve(__dirname, "../../../node_modules/.bin/tsx");


  const child = spawn(tsxPath, [proxyScript], {
    stdio: "ignore",
    detached: true
  });

  child.unref();

  // Optional: wait briefly to give server time to boot
  await new Promise(resolve => setTimeout(resolve, 300));
}

export function withProxy(fn: (...args: any[]) => any) {
  return async (...args: any[]) => {
    try {
      await ensureProxyRunning();
      return await fn(...args);
    } catch (err) {
      console.error("❌ Proxy startup failed:", err);
      process.exit(1);
    }
  };
}

