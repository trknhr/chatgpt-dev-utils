/// <reference types="vitest" />
import { describe, test, expect, vi, beforeEach } from "vitest";
import { askChatGPTViaExtension } from "./chatgptViaExtension.js";

vi.mock("open", () => ({
  default: vi.fn()
}));

vi.mock("node-fetch", () => ({
  default: vi.fn()
}));

const mockOpen = (await vi.importMock("open")).default as vi.Mock;
const mockFetch = (await vi.importMock("node-fetch")).default as vi.Mock;

describe("askChatGPTViaExtension", () => {
  beforeEach(() => {
    mockOpen.mockReset();
    mockFetch.mockReset();
  });

  test("opens ChatGPT tab", async () => {
    mockFetch.mockResolvedValue({ ok: true });

    await askChatGPTViaExtension("hello");

    expect(mockOpen).toHaveBeenCalledWith("https://chat.openai.com");
  });

  test("sends prompt to extension server", async () => {
    mockFetch.mockResolvedValue({ ok: true });

    const prompt = "generate commit message";
    await askChatGPTViaExtension(prompt);

    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:32123/chatgpt-prompt",
      expect.objectContaining({
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ type: "chatgpt-prompt", prompt })
      })
    );
  });

  test("throws error if fetch fails", async () => {
    mockFetch.mockResolvedValue({ ok: false, statusText: "Service Unavailable" });

    await expect(askChatGPTViaExtension("fail")).rejects.toThrow(
      "Failed to contact extension server: Service Unavailable"
    );
  });
});
