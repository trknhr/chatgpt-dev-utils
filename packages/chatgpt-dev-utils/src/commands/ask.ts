import { askChatGPTViaExtension } from "../utils/chatgptViaExtension";
import ora from "ora";
import { logger } from "../utils/logger";
import { presets } from "../presets/index.js";

export type AskOptions = {
  preset?: string;
};

export async function runAsk(input: string, options: AskOptions) {
  let prompt: string;

  if (options.preset) {
    const presetFn = presets[options.preset];
    if (!presetFn) {
      logger.error(`Preset '${options.preset}' not found.`);
      process.exit(1);
    }
    logger.info(`Using preset: ${options.preset}`);
    prompt = await presetFn(input ?? "");
  } else if (input) {
    prompt = input;
  } else {
    logger.error("No input or preset provided.");
    process.exit(1);
  }

  try {
    await askChatGPTViaExtension(prompt);
    logger.success("✅ Prompt sent to ChatGPT");
  } catch (err) {
    logger.error("❌ Failed to get response from ChatGPT", err);
    process.exit(1);
  }
}

// export async function runAsk(input?: string, options: AskOptions) {
//   const spinner = ora("Asking ChatGPT...").start();

//   try {
//     await askChatGPTViaExtension(input);
//     spinner.succeed("Received response:");
//   } catch (err) {
//     spinner.fail("Failed to get response from ChatGPT");
//     logger.error(err);
//     process.exit(1);
//   }
// }