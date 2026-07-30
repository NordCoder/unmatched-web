import assert from 'node:assert/strict';
import {readFile} from 'node:fs/promises';
import test from 'node:test';
import vm from 'node:vm';

const source=await readFile(new URL('./main.js',import.meta.url),'utf8');
const artSource=await readFile(new URL('./card-art.js',import.meta.url),'utf8');
const indexSource=await readFile(new URL('../static/index.html',import.meta.url),'utf8');
const mapAsset=await readFile(new URL('../static/sherwood-forest.svg',import.meta.url),'utf8');
const battlefieldManifest=JSON.parse(await readFile(new URL('../static/battlefields/sherwood-forest/manifest.json',import.meta.url),'utf8'));
const elements=new Map();
const context={
  console,
  document:{getElementById(id){if(!elements.has(id))elements.set(id,{});return elements.get(id)}},
  sessionStorage:{getItem(){return null},setItem(){},removeItem(){}},
  setInterval(){return 1},
  clearInterval(){},
};
vm.createContext(context);
vm.runInContext(artSource,context);
vm.runInContext(source,context);

test('command gate rejects duplicate in-flight submissions and resets',()=>{
  assert.equal(context.beginCommand(),true);
  assert.equal(context.beginCommand(),false);
  context.endCommand();
  assert.equal(context.beginCommand(),true);
  context.endCommand();
});

test('unchanged polling view is ignored while a new revision is applied',()=>{
  const current={match_id:'match-1',revision:7,phase:'active',viewing_player_id:'player-1'};
  assert.equal(context.shouldApplyView(current,{...current}),false);
  assert.equal(context.shouldApplyView(current,{...current,revision:8}),true);
  assert.equal(context.shouldApplyView(current,{...current},true),true);
});

test('dialog interaction survives no-op polling but not a legal-state revision change',()=>{
  const state={mode:'dialog',revision:7};
  const stable={revision:7,phase:'active',pending:null,current_player_id:'player-1',viewing_player_id:'player-1'};
  assert.equal(context.interactionCanContinue(state,stable),true);
  assert.equal(context.interactionCanContinue(state,{...stable,revision:8}),false);
  assert.equal(context.interactionCanContinue(state,{...stable,pending:{kind:'defense'}}),false);
});

test('legacy map fallback preserves the supplied battlefield aspect ratio',()=>{
  const dimensions=context.mapDimensions();
  assert.equal(dimensions.width,100);
  assert.ok(Math.abs(context.mapY(100)-(742/1337*100))<1e-9);
  assert.ok(dimensions.height<56);
});

test('Maneuver pending model separates map destinations from finish choice',()=>{
  const pending={kind:'maneuver_move',options:[
    {id:'maneuver:done',label:'Finish maneuver'},
    {id:'move-rh',fighter_id:'robin',destination:'s21',path:['s21']},
    {id:'stay-rh',fighter_id:'robin',destination:'s20',path:[]},
    {id:'move-outlaw',fighter_id:'outlaw',destination:'s14',path:['s14']},
  ]};
  const choose=context.pendingMapModel(pending,'');
  assert.deepEqual(Array.from(choose.fighters),['robin','outlaw']);
  assert.equal(choose.destinationOptions.length,0);
  assert.deepEqual(Array.from(choose.specialOptions.map(option=>option.id)),['maneuver:done']);
  const robin=context.pendingMapModel(pending,'robin');
  assert.deepEqual(Array.from(robin.destinationOptions.map(option=>option.destination)),['s21','s20']);
  assert.deepEqual(Array.from(robin.destinationOptions[1].path),[]);
});

test('attack map model uses server legal fighter and target domains',()=>{
  const legal={
    attack_cards_by_fighter:{robin:['card-1'],outlaw:[]},
    attack_targets_by_fighter:{robin:['bigfoot','jackalope'],outlaw:['bigfoot']},
  };
  const initial=context.attackMapModel(legal,'');
  assert.deepEqual(Array.from(initial.fighters),['robin']);
  const selected=context.attackMapModel(legal,'robin');
  assert.deepEqual(Array.from(selected.targetIDs),['bigfoot','jackalope']);
});

test('scheme movement uses authoritative server path and target domain',()=>{
  const action={fighters:[{
    fighter_id:'jackalope',
    destinations:[{id:'stay',destination:'s02',path:[]},{id:'move',destination:'s03',path:['s03']}],
    targets_by_destination:{s03:['bigfoot','robin']},
  }]};
  const model=context.schemeMapModel(action,'jackalope','s03');
  assert.deepEqual(Array.from(model.selectedDestination.path),['s03']);
  assert.deepEqual(Array.from(model.targetIDs),['bigfoot','robin']);
  assert.equal(context.schemeDestinationOption(context.schemeFighterAction(action,'jackalope'),'s99'),undefined);
});

test('Jackalope Horns target helper keeps target optional when server domain is empty',()=>{
  const payload={type:'scheme'};
  context.addOptionalHornsTarget(payload,null);
  assert.deepEqual(payload,{type:'scheme'});
  context.addOptionalHornsTarget(payload,{value:'bigfoot'});
  assert.equal(payload.target_id,'bigfoot');
  assert.throws(()=>context.addOptionalHornsTarget({}, {value:''}),/Choose a living fighter/);
});

test('graphical battlefield engine and calibrated surface are committed',()=>{
  assert.match(indexSource,/id="board-interaction"/);
  assert.match(indexSource,/viewBox="0 0 1337 742"/);
  assert.match(indexSource,/battlefield-renderer\.js/);
  assert.match(indexSource,/stage5-battlefield\.js/);
  assert.match(indexSource,/stage5-battlefield\.css/);
  assert.match(mapAsset,/Sherwood Forest battlefield/);
  assert.equal(battlefieldManifest.art.src,'/battlefields/sherwood-forest/board-1x.webp');
  assert.deepEqual(battlefieldManifest.art.variants.map(variant=>[variant.id,variant.src,variant.pixel_width,variant.pixel_height]),[
    ['1x','/battlefields/sherwood-forest/board-1x.webp',1337,742],
    ['2x','/battlefields/sherwood-forest/board-2x.webp',2674,1484],
  ]);
  assert.equal(battlefieldManifest.art.fallback_src,'/sherwood-forest.svg');
  assert.match(source,/MAP_ASSET = '\/sherwood-forest\.svg'/);
  assert.doesNotMatch(source,/function shortestPath/);
  assert.doesNotMatch(source,/function reachable/);
});
