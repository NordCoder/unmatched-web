import {readFile,writeFile} from 'node:fs/promises';
const generated=[
  ['apps/web/src/main.js','apps/web/static/app.js'],
  ['apps/web/src/stage4.js','apps/web/static/stage4.js'],
];
for(const [sourcePath,targetPath] of generated){
  const source=await readFile(sourcePath,'utf8');
  if(process.argv.includes('--check')){
    const target=await readFile(targetPath,'utf8');
    if(source!==target)throw new Error(`${targetPath} is stale; run npm run build`);
  }else{
    await writeFile(targetPath,source);
  }
}
for(const path of ['apps/web/static/index.html','apps/web/static/styles.css']){
  const value=await readFile(path,'utf8');
  if(value.length<100)throw new Error(`${path} is unexpectedly empty`);
}
console.log('playable web build: PASS');
