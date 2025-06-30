import fs from 'fs';
import path from 'path';

const mode = process.argv[2] === 'prod' ? 'prod' : 'dev';
const outputDir = 'dist';
fs.rmSync(outputDir, { recursive: true, force: true });
fs.mkdirSync(outputDir, { recursive: true });

fs.copyFileSync('manifest.json', path.join(outputDir, 'manifest.json'));

const replaceInFile = (filePath, placeholder, value) => {
  const content = fs.readFileSync(filePath, 'utf-8');
  const replaced = content.replace(placeholder, JSON.stringify(value));
  fs.writeFileSync(path.join(outputDir, path.basename(filePath)), replaced);
};

replaceInFile('background.js', '__LOG_LEVEL__', mode === 'prod' ? 'none' : 'log');
fs.copyFileSync('content.js', path.join(outputDir, 'content.js'));
fs.cpSync('icon', path.join(outputDir, 'icon'), { recursive: true });

console.log(`✅ Build complete for ${mode}`);
