import { select, text } from "@clack/prompts";
import { runAsk } from "./commands/ask.js";
import { withProxy } from "./utils/ensureProxyRunning.js";
import { presets } from "./presets/index.js";

const main = withProxy(async () => {
  const args = process.argv.slice(2);
  const nonInteractive = args.includes("--nonInteractive");

  const presetArg = args.find(arg => arg.startsWith("--preset="))?.split("=")[1];
  const inputArg = args.find(arg => !arg.startsWith("--"));

  let preset: string | undefined;
  let input: string | undefined;

  if (nonInteractive) {
    if (!presetArg) {
      console.error("❌ --preset option is required in non-interactive mode.");
      process.exit(1);
    }
    preset = presetArg;
    input = inputArg;
  } else {
    preset = await select({
      message: "Choose a preset to use",
      options: Object.keys(presets).map((key) => ({ value: key, label: key }))
    });

    if (!preset) {
      console.error("❌ No preset selected.");
      process.exit(1);
    }

    input = await text({ message: "Enter input or file pattern (optional):" });
  }

  await runAsk(input || "", { preset });
});

main();
