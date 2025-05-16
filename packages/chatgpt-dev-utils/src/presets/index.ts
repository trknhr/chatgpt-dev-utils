import { generateEnvReviewPrompt } from "./envReview";
import { generateTranslatePrompt } from "./translateToJa";
import { generateCommitMsgPrompt } from "./commitMsg";

export const presets: Record<string, (input: string) => Promise<string>> = {
  envReview: generateEnvReviewPrompt,
  explain: async (input: string) => {
    return `Explain what this code does:\n\n${input}`;
  },
  translateToJa: generateTranslatePrompt,
  commitMsg: generateCommitMsgPrompt,
};
