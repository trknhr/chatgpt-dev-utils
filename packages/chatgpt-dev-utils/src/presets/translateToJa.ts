import fg from "fast-glob";
import { readFile } from "fs/promises";
import { existsSync } from "fs";

export async function generateTranslatePrompt(input: string): Promise<string> {
  const paths = existsSync(input) ? [input] : await fg(input);

  if (paths.length === 0) {
    throw new Error(`❌ No files matched: ${input}`);
  }

  let combined = "";

  for (const file of paths) {
    try {
      const content = await readFile(file, "utf-8");
      combined += `\n--- ${file} ---\n${content}\n`;
    } catch {
      combined += `\n--- ${file} ---\n[Failed to read]\n`;
    }
  }

  return `Translate the following text to Japanese:\n${combined}`;
}
