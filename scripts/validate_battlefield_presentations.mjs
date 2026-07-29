import {readFile,readdir,stat} from 'node:fs/promises';
import path from 'node:path';
import vm from 'node:vm';

const repositoryRoot=path.resolve(new URL('..',import.meta.url).pathname);
const rendererSource=await readFile(path.join(repositoryRoot,'apps/web/src/battlefield-renderer.js'),'utf8');
const context={URLSearchParams,location:{search:''}};
vm.createContext(context);
vm.runInContext(rendererSource,context);
const renderer=context.BattlefieldRenderer;

async function walkManifests(directory){
  const result=[];
  for(const entry of await readdir(directory,{withFileTypes:true})){
    const full=path.join(directory,entry.name);
    if(entry.isDirectory())result.push(...await walkManifests(full));
    else if(entry.name==='manifest.json')result.push(full);
  }
  return result;
}

function yamlSpaceIDs(source){
  return [...source.matchAll(/-\s*\{id:\s*([a-z0-9-]+),/g)].map(match=>match[1]);
}

function localArtPath(src){
  if(!src.startsWith('/'))throw new Error(`art.src must be an absolute local path: ${src}`);
  return path.join(repositoryRoot,'apps/web/static',src.slice(1));
}

const manifestRoot=path.join(repositoryRoot,'apps/web/static/battlefields');
const manifests=await walkManifests(manifestRoot);
if(!manifests.length)throw new Error('No battlefield presentation manifests found');

for(const manifestPath of manifests){
  const manifest=JSON.parse(await readFile(manifestPath,'utf8'));
  const topologyPath=path.join(repositoryRoot,'docs/battlefields',`${manifest.battlefield_id}.yaml`);
  const topology=await readFile(topologyPath,'utf8');
  const runtimeIDs=yamlSpaceIDs(topology);
  const validation=renderer.validateManifest(manifest,runtimeIDs);
  if(!validation.ok)throw new Error(`${manifest.battlefield_id}: ${validation.errors.join('; ')}`);

  const artPath=localArtPath(manifest.art.src);
  const artInfo=await stat(artPath);
  if(!artInfo.isFile()||artInfo.size<10_000)throw new Error(`${manifest.battlefield_id}: art asset is missing or unexpectedly small`);

  const art=await readFile(artPath,'utf8');
  const width=Number(art.match(/\bwidth="([0-9.]+)"/)?.[1]);
  const height=Number(art.match(/\bheight="([0-9.]+)"/)?.[1]);
  if(width!==Number(manifest.coordinate_space.width)||height!==Number(manifest.coordinate_space.height)){
    throw new Error(`${manifest.battlefield_id}: art dimensions ${width}x${height} do not match coordinate space`);
  }
  if(/preserveAspectRatio="none"/.test(art)){
    throw new Error(`${manifest.battlefield_id}: art must not use preserveAspectRatio=none`);
  }

  console.log(`${manifest.battlefield_id}: ${runtimeIDs.length} calibrated spaces, ${artInfo.size} byte art asset`);
}

console.log('battlefield presentations: PASS');
