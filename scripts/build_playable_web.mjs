import {mkdir,readFile,writeFile} from 'node:fs/promises';
import path from 'node:path';

const generated=[
  ['apps/web/src/card-art.js','apps/web/static/card-art.js'],
  ['apps/web/src/battlefield-renderer.js','apps/web/static/battlefield-renderer.js'],
  ['apps/web/src/main.js','apps/web/static/app.js'],
  ['apps/web/src/stage4.js','apps/web/static/stage4.js'],
  ['apps/web/src/stage5-battlefield.js','apps/web/static/stage5-battlefield.js'],
  ['apps/web/src/stage6.js','apps/web/static/stage6.js'],
  ['apps/web/src/stage7-map-intent.js','apps/web/static/stage7-map-intent.js'],
  ['apps/web/src/tools/battlefield-editor.js','apps/web/static/tools/battlefield-editor.js'],
];

for(const [sourcePath,targetPath] of generated){
  const source=await readFile(sourcePath,'utf8');
  if(process.argv.includes('--check')){
    const target=await readFile(targetPath,'utf8');
    if(source!==target)throw new Error(`${targetPath} is stale; run npm run build`);
  }else{
    await mkdir(path.dirname(targetPath),{recursive:true});
    await writeFile(targetPath,source);
  }
}

for(const filePath of [
  'apps/web/static/index.html',
  'apps/web/static/styles.css',
  'apps/web/static/stage5-battlefield.css',
  'apps/web/static/stage6.css',
  'apps/web/static/tools/battlefield-editor.html',
  'apps/web/static/tools/battlefield-editor.css',
]){
  const value=await readFile(filePath,'utf8');
  if(value.length<100)throw new Error(`${filePath} is unexpectedly empty`);
}

console.log('playable web build: PASS');
