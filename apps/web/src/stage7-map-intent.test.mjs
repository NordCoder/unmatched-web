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
    match_id:'match-1',
    revision:12,
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

function json(value){return JSON.parse(JSON.stringify(value))}

function continuation(){
  return ux.createManeuverContinuation(baseView(),'rh','s03',['ignored-client-path']);
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
    json(ux.activationFor(view,{mode:'idle',values:{}},'s01','rh')),
    {kind:'select-fighter',fighterID:'rh'},
  );
  assert.deepEqual(
    json(ux.activationFor(view,{mode:'fighter_selected',values:{fighter_id:'rh'}},'s01','rh')),
    {kind:'clear-fighter',fighterID:'rh'},
  );
  assert.deepEqual(
    json(ux.activationFor(view,{mode:'fighter_selected',values:{fighter_id:'rh'}},'s06','outlaw')),
    {kind:'select-fighter',fighterID:'outlaw'},
  );
});

test('board activation routes enemy and empty-space clicks without action buttons',()=>{
  const view=baseView();
  const interaction={mode:'fighter_selected',values:{fighter_id:'rh'}};
  assert.deepEqual(
    json(ux.activationFor(view,interaction,'s10','bf')),
    {kind:'attack',fighterID:'rh',targetID:'bf'},
  );
  assert.deepEqual(
    json(ux.activationFor(view,interaction,'s03','')),
    {kind:'maneuver',fighterID:'rh',destination:'s03',path:['s02','s03']},
  );
  assert.deepEqual(
    json(ux.activationFor(view,interaction,'s30','jack')),
    {kind:'none'},
  );
});

test('activation consumes only projected domains and never reconstructs graph legality',()=>{
  const view=baseView();
  view.edges=[{from:'s01',to:'s30'}];
  const interaction={mode:'fighter_selected',values:{fighter_id:'rh'}};
  assert.deepEqual(
    json(ux.activationFor(view,interaction,'s30','')),
    {kind:'none'},
  );
});

test('maneuver continuation captures authoritative pre-maneuver destination and path',()=>{
  assert.deepEqual(json(continuation()),{
    kind:'maneuver',
    match_id:'match-1',
    player_id:'p1',
    source_revision:12,
    fighter_id:'rh',
    intended_destination:'s03',
    intended_path:['s02','s03'],
    selected_destination:'s03',
    selected_path:['s02','s03'],
    choice_id:'',
  });
  assert.equal(ux.createManeuverContinuation(baseView(),'rh','s30',[]),null);
});

test('continuation survives draw and BOOST revisions',()=>{
  const next={
    ...baseView(),
    revision:13,
    pending:{kind:'maneuver_boost',owner_id:'p1',options:[{id:'skip'}]},
  };
  const result=ux.reconcileManeuverContinuation(continuation(),next);
  assert.equal(result.state,'boost');
  assert.equal(result.continuation.intended_destination,'s03');
});

test('post-BOOST reconciliation selects the same destination from the final server domain',()=>{
  const next={
    ...baseView(),
    revision:14,
    pending:{kind:'maneuver_move',owner_id:'p1',options:[
      {id:'finish',label:'Finish maneuver'},
      {id:'rh-s03-final',fighter_id:'rh',destination:'s03',path:['s09','s03']},
      {id:'rh-s08-final',fighter_id:'rh',destination:'s08',path:['s09','s08']},
      {id:'outlaw-s07',fighter_id:'outlaw',destination:'s07',path:['s07']},
    ]},
  };
  const result=ux.reconcileManeuverContinuation(continuation(),next);
  assert.equal(result.state,'destination');
  assert.equal(result.option.id,'rh-s03-final');
  assert.deepEqual(Array.from(result.continuation.selected_path),['s09','s03']);
  assert.equal(result.continuation.choice_id,'rh-s03-final');
});

test('player may replace the preselected destination only with another final server option',()=>{
  const pending={kind:'maneuver_move',owner_id:'p1',options:[
    {id:'rh-s03',fighter_id:'rh',destination:'s03',path:['s02','s03']},
    {id:'rh-s08',fighter_id:'rh',destination:'s08',path:['s02','s08']},
    {id:'outlaw-s08',fighter_id:'outlaw',destination:'s08',path:['s07','s08']},
  ]};
  const selected=ux.selectManeuverContinuationDestination(continuation(),pending,'s08');
  assert.equal(selected.option.id,'rh-s08');
  assert.equal(selected.continuation.selected_destination,'s08');
  assert.deepEqual(Array.from(selected.continuation.selected_path),['s02','s08']);
  assert.equal(ux.selectManeuverContinuationDestination(continuation(),pending,'s30'),null);
});

test('stale destination and unexpected identity or pending transitions clear continuation',()=>{
  const stale={
    ...baseView(),revision:14,
    pending:{kind:'maneuver_move',owner_id:'p1',options:[
      {id:'outlaw-s07',fighter_id:'outlaw',destination:'s07',path:['s07']},
    ]},
  };
  assert.equal(ux.reconcileManeuverContinuation(continuation(),stale).reason,'destination_unavailable');

  const otherMatch={...baseView(),match_id:'match-2'};
  assert.equal(ux.reconcileManeuverContinuation(continuation(),otherMatch).reason,'match_changed');

  const otherTurn={...baseView(),revision:13,current_player_id:'p2'};
  assert.equal(ux.reconcileManeuverContinuation(continuation(),otherTurn).reason,'turn_changed');

  const missingFighter={...baseView(),revision:13};
  missingFighter.fighters=missingFighter.fighters.filter(fighter=>fighter.id!=='rh');
  assert.equal(ux.reconcileManeuverContinuation(continuation(),missingFighter).reason,'fighter_unavailable');

  const combat={
    ...baseView(),revision:13,
    pending:{kind:'defense',owner_id:'p1',options:[]},
  };
  assert.equal(ux.reconcileManeuverContinuation(continuation(),combat).reason,'unexpected_pending_kind');
});
