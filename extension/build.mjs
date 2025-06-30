import fs from 'fs';
import path from 'path';

const mode = process.argv[2] === 'prod' ? 'prod' : 'dev';
const manifestFile = mode === 'dev' ? `manifest.dev.json` : `manifest.json`;
const outputDir = 'dist';

fs.rmSync(outputDir, { recursive: true, force: true });
fs.mkdirSync(outputDir, { recursive: true });

fs.copyFileSync(
  path.resolve(manifestFile),
  path.resolve(outputDir, 'manifest.json')
);

const filesToCopy = ['background.js', 'content.js'];
filesToCopy.forEach(file => {
  fs.copyFileSync(path.resolve(file), path.resolve(outputDir, file));
});

fs.cpSync('icon', path.resolve(outputDir, 'icon'), { recursive: true });

console.log(`✅ Build complete for ${mode}`);
