import {readFile,readdir,stat} from 'node:fs/promises';
import path from 'node:path';
import vm from 'node:vm';

const repositoryRoot=path.resolve(new URL('..',import.meta.url).pathname);
const staticRoot=path.join(repositoryRoot,'apps/web/static');
const rendererSource=await readFile(path.join(repositoryRoot,'apps/web/src/battlefield-renderer.js'),'utf8');
const context={URLSearchParams,location:{search:''}};
vm.createContext(context);
vm.runInContext(rendererSource,context);
const renderer=context.BattlefieldRenderer;

const EXPECTED_SHERWOOD_ART=[
  {id:'1x',src:'/battlefields/sherwood-forest/board-1x.webp',width:1337,height:742,minBytes:300_000},
  {id:'2x',src:'/battlefields/sherwood-forest/board-2x.webp',width:2674,height:1484,minBytes:400_000},
];

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

function rejectRemoteUrls(value,location){
  if(typeof value==='string'&&/^https?:\/\//i.test(value))throw new Error(`${location} must not use a remote URL: ${value}`);
  if(Array.isArray(value))value.forEach((entry,index)=>rejectRemoteUrls(entry,`${location}[${index}]`));
  if(value&&typeof value==='object'&&!Array.isArray(value)){
    for(const [key,entry] of Object.entries(value))rejectRemoteUrls(entry,`${location}.${key}`);
  }
}

function localArtPath(src){
  if(typeof src!=='string'||!src.startsWith('/')||/^https?:\/\//i.test(src))throw new Error(`battlefield art must be a local absolute path: ${src}`);
  const resolved=path.resolve(staticRoot,src.slice(1));
  if(resolved!==staticRoot&&!resolved.startsWith(`${staticRoot}${path.sep}`))throw new Error(`battlefield art escapes static root: ${src}`);
  return resolved;
}

function read24BitLittleEndian(bytes,offset){
  return bytes[offset]+(bytes[offset+1]<<8)+(bytes[offset+2]<<16);
}

function decodeWebPHeader(bytes,assetPath){
  if(bytes.length<32)throw new Error(`${assetPath}: WebP is too small to contain a valid RIFF header`);
  if(bytes.toString('ascii',0,4)!=='RIFF'||bytes.toString('ascii',8,12)!=='WEBP')throw new Error(`${assetPath}: invalid WebP RIFF signature`);
  const riffSize=bytes.readUInt32LE(4);
  if(riffSize+8>bytes.length)throw new Error(`${assetPath}: truncated RIFF payload`);
  let offset=12;
  let hasExtendedHeader=false;
  let hasImagePayload=false;
  let width=0;
  let height=0;
  while(offset+8<=bytes.length){
    const chunkType=bytes.toString('ascii',offset,offset+4);
    const chunkSize=bytes.readUInt32LE(offset+4);
    const dataOffset=offset+8;
    const dataEnd=dataOffset+chunkSize;
    if(dataEnd>bytes.length)throw new Error(`${assetPath}: truncated ${chunkType} chunk`);
    if(chunkType==='VP8X'){
      if(chunkSize<10)throw new Error(`${assetPath}: invalid VP8X header`);
      hasExtendedHeader=true;
      width=read24BitLittleEndian(bytes,dataOffset+4)+1;
      height=read24BitLittleEndian(bytes,dataOffset+7)+1;
    }
    if(chunkType==='VP8 '||chunkType==='VP8L')hasImagePayload=true;
    offset=dataEnd+(chunkSize%2);
  }
  if(offset!==bytes.length&&offset!==bytes.length-1)throw new Error(`${assetPath}: malformed RIFF chunk alignment`);
  if(!hasExtendedHeader||!hasImagePayload)throw new Error(`${assetPath}: missing decodable VP8/VP8L payload`);
  return {width,height};
}

async function validateWebPAsset(variant,manifestPath){
  const assetPath=localArtPath(variant.src);
  const info=await stat(assetPath);
  if(!info.isFile())throw new Error(`${manifestPath}: ${variant.src} is missing`);
  if(info.size<variant.minBytes)throw new Error(`${manifestPath}: ${variant.src} is unexpectedly small (${info.size} bytes)`);
  const dimensions=decodeWebPHeader(await readFile(assetPath),variant.src);
  if(dimensions.width!==variant.width||dimensions.height!==variant.height){
    throw new Error(`${manifestPath}: ${variant.src} decodes as ${dimensions.width}x${dimensions.height}, expected ${variant.width}x${variant.height}`);
  }
}

async function validateSherwoodAssets(manifest,manifestPath){
  const variants=renderer.artVariants(manifest.art);
  if(variants.length!==EXPECTED_SHERWOOD_ART.length)throw new Error(`${manifestPath}: exactly 1x and 2x local WebP variants are required`);
  for(const [index,expected] of EXPECTED_SHERWOOD_ART.entries()){
    const variant=variants[index];
    if(variant.id!==expected.id||variant.src!==expected.src||Number(variant.pixel_width)!==expected.width||Number(variant.pixel_height)!==expected.height){
      throw new Error(`${manifestPath}: variant ${index+1} does not match the required local Sherwood dimensions/path`);
    }
    await validateWebPAsset({...variant,...expected},manifestPath);
  }
  if(manifest.art.src!==EXPECTED_SHERWOOD_ART[0].src)throw new Error(`${manifestPath}: art.src must select the local 1x variant`);
  if(manifest.art.fallback_src!=='/sherwood-forest.svg')throw new Error(`${manifestPath}: fallback must remain /sherwood-forest.svg`);
  if(Number(variants.at(-1).pixel_width)<Number(manifest.coordinate_space.width)*2||Number(variants.at(-1).pixel_height)<Number(manifest.coordinate_space.height)*2){
    throw new Error(`${manifestPath}: highest-density variant must provide at least 2x coordinate-space pixels`);
  }
  const fallbackPath=localArtPath(manifest.art.fallback_src);
  const fallback=await readFile(fallbackPath,'utf8');
  const width=Number(fallback.match(/\bwidth="([0-9.]+)"/)?.[1]);
  const height=Number(fallback.match(/\bheight="([0-9.]+)"/)?.[1]);
  if(width!==Number(manifest.coordinate_space.width)||height!==Number(manifest.coordinate_space.height)){
    throw new Error(`${manifestPath}: fallback dimensions ${width}x${height} do not match coordinate space`);
  }
}

const manifestRoot=path.join(staticRoot,'battlefields');
const manifests=await walkManifests(manifestRoot);
if(!manifests.length)throw new Error('No battlefield presentation manifests found');

for(const manifestPath of manifests){
  const manifest=JSON.parse(await readFile(manifestPath,'utf8'));
  rejectRemoteUrls(manifest,manifestPath);
  const topologyPath=path.join(repositoryRoot,'docs/battlefields',`${manifest.battlefield_id}.yaml`);
  const topology=await readFile(topologyPath,'utf8');
  const runtimeIDs=yamlSpaceIDs(topology);
  const validation=renderer.validateManifest(manifest,runtimeIDs);
  if(!validation.ok)throw new Error(`${manifest.battlefield_id}: ${validation.errors.join('; ')}`);
  if(manifest.battlefield_id==='sherwood-forest')await validateSherwoodAssets(manifest,manifestPath);
  console.log(`${manifest.battlefield_id}: ${runtimeIDs.length} calibrated spaces, local 1x/2x WebP, local fallback`);
}

const serverSource=await readFile(path.join(repositoryRoot,'internal/playableslice/server/server.go'),'utf8');
if(!/img-src 'self' data:/.test(serverSource))throw new Error('server CSP must allow only local/data image sources');

console.log('battlefield presentations: PASS');
