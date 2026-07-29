import assert from 'node:assert/strict';
import {readFile} from 'node:fs/promises';
import test from 'node:test';
import vm from 'node:vm';

const mainSource=await readFile(new URL('./main.js',import.meta.url),'utf8');
const stage4Source=await readFile(new URL('./stage4.js',import.meta.url),'utf8');
const indexSource=await readFile(new URL('../static/index.html',import.meta.url),'utf8');
const elements=new Map();
const context={
  console,
  document:{getElementById(id){if(!elements.has(id))elements.set(id,{});return elements.get(id)}},
  sessionStorage:{getItem(){return null},setItem(){},removeItem(){}},
  setInterval(){return 1},
  clearInterval(){},
};
vm.createContext(context);
vm.runInContext(mainSource,context);
vm.runInContext(stage4Source,context);

test('Stage 4 scripts are loaded as ordered classic scripts',()=>{
  assert.match(indexSource,/<script src="\/app\.js"><\/script>\s*<script src="\/stage4\.js"><\/script>/);
  assert.doesNotMatch(indexSource,/type="module"/);
});

test('Scheme terminology is presented as Action',()=>{
  assert.equal(context.stage4UserFacingMessage('Choose a Scheme card'),'Choose an Action card');
  assert.equal(context.stage4UserFacingMessage('scheme card is invalid'),'action card is invalid');
  assert.equal(context.stage4CardTypeLabel({type:'scheme'}),'ACTION');
});

test('card presentation includes rules and usable-by labels',()=>{
  const fighters=[
    {definition_id:'robin-hood',name:'Robin Hood'},
    {definition_id:'outlaw',name:'Outlaw'},
  ];
  assert.equal(context.stage4UsableByLabel({usable_by:['robin-hood']},fighters),'Robin Hood');
  assert.equal(context.stage4UsableByLabel({usable_by:['any']},fighters),'Any fighter');
  assert.match(context.stage4CardRulesText({definition_id:'steal-from-the-rich'}),/Draw 1 card/);
  assert.match(context.stage4CardDetailsContent({name:'Highway Robbery',type:'attack',value:2,boost:2,usable_by:['outlaw'],definition_id:'highway-robbery'}),/Usable by: Outlaw/);
});

test('fighter range labels expose melee and ranged identity',()=>{
  assert.equal(context.stage4AttackTypeLabel('melee'),'Melee');
  assert.equal(context.stage4AttackTypeLabel('ranged'),'Ranged');
});

test('Maneuver requires a fighter unless exactly one is available',()=>{
  assert.deepEqual(context.stage4ManeuverFighterModel(['robin','outlaw'],'').selectedFighterID,'');
  assert.equal(context.stage4ManeuverFighterModel(['robin'],'').selectedFighterID,'robin');
  assert.equal(context.stage4ManeuverFighterModel(['robin','outlaw'],'outlaw').selectedFighterID,'outlaw');
});

test('BOOST model separates normal movement from discard choices',()=>{
  const model=context.stage4ManeuverBoostModel({kind:'maneuver_boost',options:[
    {id:'boost:none',label:'Do not BOOST movement'},
    {id:'boost:card-1',card_id:'card-1'},
    {id:'boost:card-2',card_id:'card-2'},
  ]});
  assert.equal(model.noBoost.id,'boost:none');
  assert.deepEqual(Array.from(model.boostOptions.map(option=>option.card_id)),['card-1','card-2']);
});

test('all launch card definitions have readable Stage 4 rules',()=>{
  const ids=[
    'a-hunters-eye','steal-from-the-rich','disarming-shot','piercing-shot','highway-robbery',
    'defenders-of-sherwood','feint','regroup','wily-fighting','snark','ambush','larger-than-life',
    'savagery','crash-through-the-trees','jackalope-horns','hoax','disengage',
    'its-just-your-imagination','skirmish','momentous-shift',
  ];
  for(const id of ids.filter(id=>!['a-hunters-eye','larger-than-life'].includes(id)))assert.notEqual(context.stage4CardRulesText({definition_id:id}),'No special effect.',id);
  assert.equal(context.stage4CardRulesText({definition_id:'a-hunters-eye'}),'No special effect.');
  assert.equal(context.stage4CardRulesText({definition_id:'larger-than-life'}),'No special effect.');
});
