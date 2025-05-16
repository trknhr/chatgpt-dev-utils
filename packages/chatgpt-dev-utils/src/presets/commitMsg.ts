import { execSync } from "child_process";

export async function generateCommitMsgPrompt(): Promise<string> {
  let diff = "";

  try {
    diff = execSync("git diff --cached", { encoding: "utf-8" });
  } catch (err) {
    throw new Error("Failed to get git diff: " + err);
  }

  if (!diff.trim()) {
    throw new Error("No staged changes found. Please stage your changes with `git add`.");
  }

  return `Please generate a conventional commit message in English for the following Git diff.

- The message should be clear and descriptive.
- Include short explanation bullets if needed.
- It doesn't have to be a single line.
- Make plain text

Git Diff:
${diff}`;
}
