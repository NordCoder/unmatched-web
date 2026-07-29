import {readFile,writeFile} from 'node:fs/promises';
const source=await readFile('apps/web/src/main.js','utf8');
const targetPath='apps/web/static/app.js';
if(process.argv.includes('--check')){
  const target=await readFile(targetPath,'utf8');
  if(source!==target)throw new Error('apps/web/static/app.js is stale; run npm run build');
}else{
  await writeFile(targetPath,source);
}
for(const path of ['apps/web/static/index.html','apps/web/static/styles.css']){
  const value=await readFile(path,'utf8');
  if(value.length<100)throw new Error(`${path} is unexpectedly empty`);
}
console.log('playable web build: PASS');
