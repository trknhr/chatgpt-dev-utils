import { afterEach, expect, test, vi,describe } from "vitest";
import { runAsk } from "./ask.js";
import * as chatgpt from "../utils/chatgptViaExtension.js";

vi.spyOn(chatgpt, "askChatGPTViaExtension").mockResolvedValue(undefined);

describe("ask", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  test("runAsk sends prompt to extension", async () => {
    const spy = vi.spyOn(chatgpt, "askChatGPTViaExtension");
    await runAsk("Hello?");
    expect(spy).toHaveBeenCalledWith("Hello?");
  });
})
