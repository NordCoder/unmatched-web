import assert from 'node:assert/strict';
import {readFile} from 'node:fs/promises';
import test from 'node:test';
import vm from 'node:vm';

const source=await readFile(new URL('./main.js',import.meta.url),'utf8');
const elements=new Map();
const context={
  console,
  document:{getElementById(id){if(!elements.has(id))elements.set(id,{});return elements.get(id)}},
  sessionStorage:{getItem(){return null},setItem(){},removeItem(){}},
  setInterval(){return 1},
  clearInterval(){},
};
vm.createContext(context);
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

test('Jackalope Horns omits target when the selector is empty',()=>{
  const payload={type:'scheme'};
  context.addOptionalHornsTarget(payload,null);
  assert.deepEqual(payload,{type:'scheme'});
});

test('Jackalope Horns includes a selected target and rejects an empty select',()=>{
  const payload={type:'scheme'};
  context.addOptionalHornsTarget(payload,{value:'player-2-bigfoot'});
  assert.equal(payload.target_id,'player-2-bigfoot');
  assert.throws(()=>context.addOptionalHornsTarget({}, {value:''}),/Choose a living fighter/);
});

test('scheme movement uses authoritative server path and target domain',()=>{
  const action={fighters:[{
    fighter_id:'player-2-jackalope',
    destinations:[{id:'stay',destination:'s02',path:[]},{id:'move',destination:'s03',path:['s03']}],
    targets_by_destination:{s03:['player-2-bigfoot','player-1-robin-hood']},
  }]};
  const fighter=context.schemeFighterAction(action,'player-2-jackalope');
  const destination=context.schemeDestinationOption(fighter,'s03');
  assert.deepEqual(Array.from(destination.path),['s03']);
  assert.deepEqual(Array.from(context.schemeTargetIDs(fighter,'s03')),['player-2-bigfoot','player-1-robin-hood']);
  assert.equal(context.schemeDestinationOption(fighter,'s99'),undefined);
});
