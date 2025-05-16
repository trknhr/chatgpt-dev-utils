import open from "open";
import fetch from "node-fetch";

export async function askChatGPTViaExtension(prompt: string): Promise<void> {
  // optionally open ChatGPT tab if not open
  await open("https://chatgpt.com");

  const response = await fetch(`http://localhost:32123/chatgpt-prompt`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ type: "chatgpt-prompt", prompt })
  });

  if (!response.ok) {
    throw new Error(`Failed to contact extension server: ${response.statusText}`);
  }
}
