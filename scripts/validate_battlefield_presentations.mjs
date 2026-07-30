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
  const spacesSection=source.match(/^spaces:\s*\n([\s\S]*?)^edges:\s*$/m)?.[1];
  if(!spacesSection)throw new Error('battlefield YAML is missing spaces/edges sections');
  return [...spacesSection.matchAll(/^\s*-\s*\{id:\s*([a-z0-9-]+),/gm)].map(match=>match[1]);
}

function localArtPath(src){
  if(!src.startsWith('/'))throw new Error(`battlefield art must be an absolute local path: ${src}`);
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

  const variants=renderer.artVariants(manifest.art);
  if(variants.length<2)throw new Error(`${manifest.battlefield_id}: local 1x and 2x art variants are required`);
  let previousWidth=0;
  for(const variant of variants){
    if(/^https?:\/\//.test(variant.src||''))throw new Error(`${manifest.battlefield_id}: runtime art variant ${variant.id||'unknown'} must not be remote`);
    const width=Number(variant.pixel_width);
    const height=Number(variant.pixel_height);
    if(!Number.isFinite(width)||!Number.isFinite(height)||width<=0||height<=0){
      throw new Error(`${manifest.battlefield_id}: variant ${variant.id||'unknown'} has invalid pixel dimensions`);
    }
    if(width<previousWidth)throw new Error(`${manifest.battlefield_id}: art variants must be ordered by density`);
    previousWidth=width;
    const artInfo=await stat(localArtPath(variant.src));
    if(!artInfo.isFile()||artInfo.size<100_000){
      throw new Error(`${manifest.battlefield_id}: local variant ${variant.id||'unknown'} is missing or too aggressively compressed`);
    }
  }
  const highest=variants[variants.length-1];
  if(Number(highest.pixel_width)<Number(manifest.coordinate_space.width)*2||Number(highest.pixel_height)<Number(manifest.coordinate_space.height)*2){
    throw new Error(`${manifest.battlefield_id}: highest-density variant must provide at least 2x coordinate-space pixels`);
  }

  const fallbackPath=localArtPath(manifest.art.fallback_src||manifest.art.src);
  const fallbackInfo=await stat(fallbackPath);
  if(!fallbackInfo.isFile())throw new Error(`${manifest.battlefield_id}: fallback art asset is missing`);
  const fallback=await readFile(fallbackPath,'utf8');
  const width=Number(fallback.match(/\bwidth="([0-9.]+)"/)?.[1]);
  const height=Number(fallback.match(/\bheight="([0-9.]+)"/)?.[1]);
  if(width!==Number(manifest.coordinate_space.width)||height!==Number(manifest.coordinate_space.height)){
    throw new Error(`${manifest.battlefield_id}: fallback dimensions ${width}x${height} do not match coordinate space`);
  }

  console.log(`${manifest.battlefield_id}: ${runtimeIDs.length} calibrated spaces, ${variants.length} local density variants, local fallback`);
}

console.log('battlefield presentations: PASS');
