import assert from 'node:assert/strict';
import {readFile} from 'node:fs/promises';
import test from 'node:test';
import vm from 'node:vm';

const source=await readFile(new URL('./battlefield-renderer.js',import.meta.url),'utf8');
const manifest=JSON.parse(await readFile(new URL('../static/battlefields/sherwood-forest/manifest.json',import.meta.url),'utf8'));
const context={console,URLSearchParams,location:{search:''}};
vm.createContext(context);
vm.runInContext(source,context);
const renderer=context.BattlefieldRenderer;
const spaceIDs=Array.from({length:30},(_,index)=>`s${String(index+1).padStart(2,'0')}`);

test('Sherwood Forest presentation has exact complete geometry',()=>{
  const result=renderer.validateManifest(manifest,spaceIDs);
  assert.deepEqual(Array.from(result.errors),[]);
  assert.equal(result.ok,true);
  assert.equal(Object.keys(manifest.spaces).length,30);
  assert.deepEqual(manifest.coordinate_space,{width:1337,height:742});
  assert.deepEqual(manifest.spaces.s01.shape,{type:'circle',cx:292.2,cy:144.6,r:62.8});
  assert.deepEqual(manifest.spaces.s30.shape,{type:'circle',cx:1062.6,cy:516.6,r:63.1});
});

test('manifest validation rejects missing, unknown and out-of-bounds geometry',()=>{
  const broken=structuredClone(manifest);
  delete broken.spaces.s01;
  broken.spaces.s99={shape:{type:'circle',cx:1,cy:1,r:20}};
  const result=renderer.validateManifest(broken,spaceIDs);
  assert.equal(result.ok,false);
  assert.ok(Array.from(result.errors).some(error=>error.includes('s01: presentation geometry is missing')));
  assert.ok(Array.from(result.errors).some(error=>error.includes('s99: presentation geometry has no runtime space')));
  assert.ok(Array.from(result.errors).some(error=>error.includes('s99: circle exceeds coordinate space')));
});

test('renderer supports circle, ellipse and polygon geometry',()=>{
  const sample={
    schema_version:1,
    battlefield_id:'shapes',
    coordinate_space:{width:200,height:100},
    art:{src:'/board.png'},
    spaces:{
      circle:{shape:{type:'circle',cx:30,cy:30,r:20}},
      ellipse:{shape:{type:'ellipse',cx:90,cy:30,rx:25,ry:15}},
      polygon:{shape:{type:'polygon',points:[[130,10],[180,10],[155,50]]}},
    },
  };
  assert.equal(renderer.validateManifest(sample,['circle','ellipse','polygon']).ok,true);
  assert.match(renderer.shapeMarkup(sample.spaces.circle.shape),/<circle /);
  assert.match(renderer.shapeMarkup(sample.spaces.ellipse.shape),/<ellipse /);
  assert.match(renderer.shapeMarkup(sample.spaces.polygon.shape),/<polygon /);
});

test('highlight geometry is inset while hit geometry expands',()=>{
  const circle={type:'circle',cx:100,cy:80,r:60};
  assert.equal(renderer.adjustShape(circle,-6).r,54);
  assert.equal(renderer.adjustShape(circle,9).r,69);
});

test('art variants select the smallest source that satisfies display density',()=>{
  const sample={art:{variants:[
    {id:'1x',src:'/board-1x.webp',pixel_width:1337,pixel_height:742},
    {id:'2x',src:'/board-2x.webp',pixel_width:2674,pixel_height:1484},
  ]}};
  assert.equal(renderer.selectArtVariant(sample,{displayWidth:900,devicePixelRatio:1}).id,'1x');
  assert.equal(renderer.selectArtVariant(sample,{displayWidth:900,devicePixelRatio:2}).id,'2x');
  assert.equal(renderer.selectArtVariant(sample,{displayWidth:1800,devicePixelRatio:2}).id,'2x');
});

test('calibrated renderer keeps art, large fighter art and overlay in one coordinate system',()=>{
  renderer.register(manifest);
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
    players:[{id:'p1',hero_id:'rh'},{id:'p2',hero_id:'bf'}],
    spaces:spaceIDs.map(id=>({
      id,
      zones:['orange'],
      ...(id==='s01'?{fighter:{id:'rh',definition_id:'robin-hood',name:'Robin Hood',health:13,owner_id:'p1',defeated:false}}:{}),
      ...(id==='s02'?{fighter:{id:'outlaw',definition_id:'outlaw',name:'Outlaw',health:1,owner_id:'p1',defeated:false}}:{}),
    })),
    edges:[{from:'s01',to:'s02'}],
  };
  renderer.render({
    svg,
    view,
    highlights:{
      destinations:new Map([['s02',{}]]),
      fighterCandidates:new Set(['rh']),
      targets:new Set(),
      selectedFighterID:'rh',
      selectedTargetID:'',
    },
    viewerID:'p1',
    initials:()=> 'RH',
    debug:true,
  });
  assert.equal(svg.attrs.get('viewBox'),'0 0 1337 742');
  assert.equal(svg.attrs.get('preserveAspectRatio'),'xMidYMid meet');
  assert.equal(svg.dataset.presentation,'calibrated');
  assert.match(svg.innerHTML,/unmatchedpicks\.com\/maps\/sherwoodforest\.webp/);
  assert.match(svg.innerHTML,/cx="292\.2" cy="144\.6" r="71\.8"/);
  assert.match(svg.innerHTML,/class="space-highlight"[^>]*r="56\.8"/);
  assert.match(svg.innerHTML,/class="fighter-piece hero-piece"/);
  assert.match(svg.innerHTML,/class="fighter-piece sidekick-piece"/);
  assert.match(svg.innerHTML,/href="\/fighters\/tokens\/robin-hood\.svg"/);
  assert.match(svg.innerHTML,/href="\/fighters\/tokens\/outlaw\.svg"/);
  assert.match(svg.innerHTML,/class="fighter-art"[^>]*width="102\./);
  assert.match(svg.innerHTML,/class="fighter-art"[^>]*width="97\./);
  assert.match(svg.innerHTML,/class="health-token"/);
  assert.match(svg.innerHTML,/class="debug-edge"/);
  assert.doesNotMatch(svg.innerHTML,/preserveAspectRatio="none"/);
});

test('fighter token paths are derived only from safe stable definitions',()=>{
  for(const definitionID of ['robin-hood','outlaw','bigfoot','jackalope']){
    assert.equal(renderer.fighterTokenHref({definition_id:definitionID}),`/fighters/tokens/${definitionID}.svg`);
  }
  assert.equal(renderer.fighterTokenHref({definition_id:'../robin-hood'}),'');
  assert.equal(renderer.fighterTokenHref({definition_id:'Robin Hood'}),'');
});

test('fallback presentation preserves legacy maps without a graphical manifest',()=>{
  const view={battlefield_id:'unknown-map',spaces:[{id:'a',x:50,y:50}]};
  const fallback=renderer.presentationFor(view);
  assert.equal(fallback.fallback,true);
  assert.equal(renderer.point(fallback,'a').x,668.5);
  assert.equal(renderer.point(fallback,'a').y,371);
});

test('debug mode is explicit',()=>{
  assert.equal(renderer.debugEnabled('?battlefieldDebug=1'),true);
  assert.equal(renderer.debugEnabled('?battlefieldDebug=0'),false);
});
