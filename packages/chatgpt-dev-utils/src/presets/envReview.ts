import { readFile } from "fs/promises";

export async function generateEnvReviewPrompt(): Promise<string> {
  let out = "";

  try {
    out += "--- docker-compose.yml ---\n";
    out += (await readFile("docker-compose.yml", "utf-8"))
      .split("\n")
      .slice(0, 100)
      .join("\n") + "\n\n";
  } catch {}

  try {
    out += "--- .env (sanitized) ---\n";
    out += (await readFile(".env", "utf-8"))
      .split("\n")
      .filter((line) => !/SECRET|TOKEN|PASSWORD/i.test(line))
      .slice(0, 50)
      .join("\n") + "\n\n";
  } catch {}

  try {
    out += "--- package.json ---\n";
    out += (await readFile("package.json", "utf-8"))
      .split("\n")
      .slice(0, 100)
      .join("\n") + "\n";
  } catch {}

  return `Please analyze the following project setup:\n\n${out}`;
}
