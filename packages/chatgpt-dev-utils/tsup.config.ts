import { defineConfig } from "tsup";

export default defineConfig({
  entry: ["src/chatgpt-dev-utils.ts"],
  format: ["esm"], // ← CJSからESMに戻す
  target: "es2020",
  outDir: "dist",
  clean: true,
  external: ["fast-glob"], // ← commander も外に出す！
  banner: {
    js: "#!/usr/bin/env node"
  }
});
