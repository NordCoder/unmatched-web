import assert from 'node:assert/strict';
import test from 'node:test';
import vm from 'node:vm';
import {readFile} from 'node:fs/promises';

const source=await readFile(new URL('./stage7-map-intent.js',import.meta.url),'utf8');
const context={console};
context.globalThis=context;
vm.createContext(context);
vm.runInContext(source,context);
const ux=context.Stage7MapIntent;

function baseView(){
  return {
    phase:'active',
    current_player_id:'p1',
    viewing_player_id:'p1',
    fighters:[
      {id:'rh',owner_id:'p1',defeated:false},
      {id:'outlaw',owner_id:'p1',defeated:false},
      {id:'bf',owner_id:'p2',defeated:false},
      {id:'jack',owner_id:'p2',defeated:false},
    ],
    legal:{
      can_maneuver:true,
      maneuver_action:{
        base_movement:2,
        destinations_by_fighter:{
          rh:[
            {destination:'s02',path:['s02']},
            {destination:'s03',path:['s02','s03']},
          ],
          outlaw:[{destination:'s07',path:['s07']}],
        },
      },
      attack_cards_by_fighter:{rh:['attack-1'],outlaw:['attack-2']},
      attack_targets_by_fighter:{rh:['bf'],outlaw:['jack']},
    },
  };
}

test('fighter intent combines server-projected movement and attack domains',()=>{
  const intent=ux.fighterIntentModel(baseView(),'rh');
  assert.equal(intent.fighterID,'rh');
  assert.equal(intent.baseMovement,2);
  assert.deepEqual(Array.from(intent.destinations,option=>option.destination),['s02','s03']);
  assert.deepEqual(Array.from(intent.targetIDs),['bf']);
});

test('fighter intent is unavailable outside the private active action window',()=>{
  const opponent=baseView();
  opponent.viewing_player_id='p2';
  assert.equal(ux.fighterIntentModel(opponent,'rh').fighterID,'');

  const pending=baseView();
  pending.pending={kind:'defense'};
  assert.equal(ux.fighterIntentModel(pending,'rh').destinations.length,0);

  const waiting=baseView();
  waiting.current_player_id='p2';
  assert.equal(ux.fighterIntentModel(waiting,'rh').targetIDs.length,0);
});

test('board activation selects and clears friendly fighters directly',()=>{
  const view=baseView();
  assert.deepEqual(
    JSON.parse(JSON.stringify(ux.activationFor(view,{mode:'idle',values:{}},'s01','rh'))),
    {kind:'select-fighter',fighterID:'rh'},
  );
  assert.deepEqual(
    JSON.parse(JSON.stringify(ux.activationFor(view,{mode:'fighter_selected',values:{fighter_id:'rh'}},'s01','rh'))),
    {kind:'clear-fighter',fighterID:'rh'},
  );
  assert.deepEqual(
    JSON.parse(JSON.stringify(ux.activationFor(view,{mode:'fighter_selected',values:{fighter_id:'rh'}},'s06','outlaw'))),
    {kind:'select-fighter',fighterID:'outlaw'},
  );
});

test('board activation routes enemy and empty-space clicks without action buttons',()=>{
  const view=baseView();
  const interaction={mode:'fighter_selected',values:{fighter_id:'rh'}};
  assert.deepEqual(
    JSON.parse(JSON.stringify(ux.activationFor(view,interaction,'s10','bf'))),
    {kind:'attack',fighterID:'rh',targetID:'bf'},
  );
  assert.deepEqual(
    JSON.parse(JSON.stringify(ux.activationFor(view,interaction,'s03',''))),
    {kind:'maneuver',fighterID:'rh',destination:'s03',path:['s02','s03']},
  );
  assert.deepEqual(
    JSON.parse(JSON.stringify(ux.activationFor(view,interaction,'s30','jack'))),
    {kind:'none'},
  );
});

test('activation consumes only projected domains and never reconstructs graph legality',()=>{
  const view=baseView();
  view.edges=[{from:'s01',to:'s30'}];
  const interaction={mode:'fighter_selected',values:{fighter_id:'rh'}};
  assert.deepEqual(
    JSON.parse(JSON.stringify(ux.activationFor(view,interaction,'s30',''))),
    {kind:'none'},
  );
});
