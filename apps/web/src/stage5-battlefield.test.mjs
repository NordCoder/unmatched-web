import assert from 'node:assert/strict';
import {readFile} from 'node:fs/promises';
import test from 'node:test';
import vm from 'node:vm';

const rendererSource=await readFile(new URL('./battlefield-renderer.js',import.meta.url),'utf8');
const stage5Source=await readFile(new URL('./stage5-battlefield.js',import.meta.url),'utf8');
const manifest=JSON.parse(await readFile(new URL('../static/battlefields/sherwood-forest/manifest.json',import.meta.url),'utf8'));

const loading={hidden:true};
const loadingLabel={textContent:''};
const game={
  attributes:new Map(),
  setAttribute(name,value){this.attributes.set(name,value)},
  removeAttribute(name){this.attributes.delete(name)},
};
const svg={
  attrs:new Map(),
  style:{},
  dataset:{},
  setAttribute(name,value){this.attrs.set(name,value)},
  innerHTML:'',
  querySelectorAll(){return []},
};
class FakeImage{
  constructor(){this.complete=false;this.naturalWidth=0;this.onload=null;this.onerror=null}
  set src(value){
    this._src=value;
    this.complete=true;
    this.naturalWidth=2674;
    queueMicrotask(()=>this.onload?.());
  }
  get src(){return this._src}
  decode(){return Promise.resolve()}
}
let localRenderCount=0;
const view={
  match_id:'match-1',
  battlefield_id:'sherwood-forest',
  viewing_player_id:'p1',
  players:[{id:'p1',hero_id:'rh'}],
  fighters:[{id:'rh',definition_id:'robin-hood',name:'Robin Hood',health:13,owner_id:'p1'}],
  spaces:Object.keys(manifest.spaces).map(id=>({
    id,
    zones:['orange'],
    ...(id==='s01'?{fighter:{id:'rh',definition_id:'robin-hood',name:'Robin Hood',health:13,owner_id:'p1'}}:{}),
  })),
  edges:[],
};
const context={
  console,
  URLSearchParams,
  location:{search:''},
  innerWidth:1200,
  devicePixelRatio:2,
  setTimeout,
  clearTimeout,
  queueMicrotask,
  Image:FakeImage,
  session:null,
  view,
  game,
  document:{
    querySelector(selector){return selector==='.board-frame'?{clientWidth:1000}:null},
  },
  applyView(next){this.view=next;return true},
  enterGame(){this.entered=true},
  leave(){this.left=true},
  mapDimensions(){return {width:100,height:55}},
  mapY(value){return value},
  boardPoint(){return ''},
  renderBoard(){},
  renderLocalInteraction(){localRenderCount++},
  boardHighlights(){return {destinations:new Map(),fighterCandidates:new Set(),targets:new Set(),selectedFighterID:'',selectedTargetID:''}},
  initials(name){return name.slice(0,2)},
  $(id){
    if(id==='board')return svg;
    if(id==='battlefield-loading')return loading;
    if(id==='battlefield-loading-label')return loadingLabel;
    return {};
  },
};
vm.createContext(context);
vm.runInContext(rendererSource,context);
context.BattlefieldRenderer.register(structuredClone(manifest));
vm.runInContext(stage5Source,context);

test('Stage 6 keeps calibrated map dimensions and direct graphical tokens',()=>{
  const dimensions=context.mapDimensions();
  assert.equal(dimensions.width,1337);
  assert.equal(dimensions.height,742);
  assert.equal(context.mapY(50),371);
  assert.equal(context.boardPoint('s01'),'292.2,144.6');
  context.renderBoard();
  assert.equal(svg.attrs.get('viewBox'),'0 0 1337 742');
  assert.equal(svg.dataset.presentation,'calibrated');
  assert.match(svg.innerHTML,/data-space="s01"/);
  assert.match(svg.innerHTML,/href="\/fighters\/tokens\/robin-hood\.svg"/);
  assert.match(svg.innerHTML,/class="fighter-piece hero-piece"/);
});

test('Stage 6 configures only local density variants and role-based token ratios',()=>{
  const configured=context.stage5ConfigurePresentation(context.BattlefieldRenderer.get('sherwood-forest'));
  assert.equal(configured.art.variants.length,2);
  assert.equal(
    JSON.stringify(configured.art.variants.map(variant=>[variant.id,variant.src,variant.pixel_width,variant.pixel_height])),
    JSON.stringify([
      ['1x','/battlefields/sherwood-forest/board-1x.webp',1337,742],
      ['2x','/battlefields/sherwood-forest/board-2x.webp',2674,1484],
    ]),
  );
  assert.equal(configured.art.variants.some(variant=>/^https?:/.test(variant.src)),false);
  assert.equal(configured.defaults.hero_token_diameter_ratio,.82);
  assert.equal(configured.defaults.sidekick_token_diameter_ratio,.77);
  assert.equal(
    context.BattlefieldRenderer.fighterTokenHref({definition_id:'jackalope'}),
    '/fighters/tokens/jackalope.svg',
  );
});

test('Stage 6 loading overlay covers preload and clears after decoded assets are ready',async()=>{
  context.stage5ShowLoading('Preparing test battlefield…');
  assert.equal(loading.hidden,false);
  assert.equal(loadingLabel.textContent,'Preparing test battlefield…');
  assert.equal(game.attributes.get('aria-busy'),'true');

  await context.stage5PrepareAssets(view);
  assert.equal(loading.hidden,true);
  assert.equal(game.attributes.has('aria-busy'),false);
  const active=context.BattlefieldRenderer.get('sherwood-forest').art.active_variant;
  assert.equal(active.id,'2x');
  assert.equal(active.src,'/battlefields/sherwood-forest/board-2x.webp');
  assert.ok(localRenderCount>=1);
});

test('all four bundled fighter token SVGs contain valid embedded WebP payloads',async()=>{
  for(const definitionID of ['robin-hood','outlaw','bigfoot','jackalope']){
    const source=await readFile(new URL(`../static/fighters/tokens/${definitionID}.svg`,import.meta.url),'utf8');
    assert.match(source,/^<svg /);
    const encoded=source.match(/data:image\/webp;base64,([^"']+)/)?.[1];
    assert.ok(encoded,`${definitionID} must contain embedded WebP art`);
    const bytes=Buffer.from(encoded,'base64');
    assert.equal(bytes.subarray(0,4).toString('ascii'),'RIFF');
    assert.equal(bytes.subarray(8,12).toString('ascii'),'WEBP');
  }
});
