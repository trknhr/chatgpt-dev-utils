import fs from 'fs';
import path from 'path';

const mode = process.argv[2] === 'prod' ? 'prod' : 'dev';
const manifestFile = `manifest.${mode}.json`;
const outputDir = 'dist';

// 1. distディレクトリ作成
fs.rmSync(outputDir, { recursive: true, force: true });
fs.mkdirSync(outputDir, { recursive: true });

// 2. manifestをコピー
fs.copyFileSync(
  path.resolve(manifestFile),
  path.resolve(outputDir, 'manifest.json')
);

// 3. その他ファイルコピー
const filesToCopy = ['background.js', 'content.js', 'popup.html'];
filesToCopy.forEach(file => {
  fs.copyFileSync(path.resolve(file), path.resolve(outputDir, file));
});

// 4. アイコンなど（任意）
fs.cpSync('icon', path.resolve(outputDir, 'icon'), { recursive: true });

console.log(`✅ Build complete for ${mode}`);
