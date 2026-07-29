import assert from 'node:assert/strict';
import {readFile} from 'node:fs/promises';
import test from 'node:test';
import vm from 'node:vm';

const rendererSource=await readFile(new URL('./battlefield-renderer.js',import.meta.url),'utf8');
const stage5Source=await readFile(new URL('./stage5-battlefield.js',import.meta.url),'utf8');
const manifest=JSON.parse(await readFile(new URL('../static/battlefields/sherwood-forest/manifest.json',import.meta.url),'utf8'));
const svg={
  attrs:new Map(),
  style:{},
  dataset:{},
  setAttribute(name,value){this.attrs.set(name,value)},
  innerHTML:'',
};
const view={
  battlefield_id:'sherwood-forest',
  viewing_player_id:'p1',
  spaces:Object.keys(manifest.spaces).map(id=>({id,zones:['orange']})),
  edges:[],
};
const context={
  console,
  URLSearchParams,
  location:{search:''},
  view,
  applyView(next){this.view=next;return true},
  mapDimensions(){return {width:100,height:55}},
  mapY(value){return value},
  boardPoint(){return ''},
  renderBoard(){},
  renderLocalInteraction(){},
  boardHighlights(){return {destinations:new Map(),fighterCandidates:new Set(),targets:new Set(),selectedFighterID:'',selectedTargetID:''}},
  initials(name){return name.slice(0,2)},
  $(id){return id==='board'?svg:{}},
};
vm.createContext(context);
vm.runInContext(rendererSource,context);
context.BattlefieldRenderer.register(manifest);
vm.runInContext(stage5Source,context);

test('Stage 5 overrides legacy percentage dimensions with manifest dimensions',()=>{
  const dimensions=context.mapDimensions();
  assert.equal(dimensions.width,1337);
  assert.equal(dimensions.height,742);
  assert.equal(context.mapY(50),371);
  assert.equal(context.boardPoint('s01'),'292.2,144.6');
});

test('Stage 5 renderer uses calibrated geometry',()=>{
  context.renderBoard();
  assert.equal(svg.attrs.get('viewBox'),'0 0 1337 742');
  assert.equal(svg.dataset.presentation,'calibrated');
  assert.match(svg.innerHTML,/data-space="s01"/);
  assert.match(svg.innerHTML,/cx="292\.2" cy="144\.6" r="71\.8"/);
});

test('fighter token assets are selected from stable fighter definition identities',()=>{
  for(const definitionID of ['robin-hood','outlaw','bigfoot','jackalope']){
    assert.equal(
      context.stage5FighterTokenHref({definition_id:definitionID}),
      `/fighters/tokens/${definitionID}.svg`,
    );
  }
  assert.equal(context.stage5FighterTokenHref({definition_id:'../robin-hood'}),'');
  assert.equal(context.stage5FighterTokenHref({definition_id:'Robin Hood'}),'');
  assert.equal(context.stage5FighterTokenHref({}),'');
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
