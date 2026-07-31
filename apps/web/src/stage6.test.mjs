import assert from 'node:assert/strict';
import test from 'node:test';
import vm from 'node:vm';
import {readFile} from 'node:fs/promises';

const source=await readFile(new URL('./stage6.js',import.meta.url),'utf8');
const context={console,setTimeout,clearTimeout,URLSearchParams};
context.globalThis=context;
vm.createContext(context);
vm.runInContext(source,context);
const ux=context.Stage6UX;

function baseView(){
  return {
    phase:'active',
    current_player_id:'p1',
    viewing_player_id:'p1',
    players:[
      {id:'p1',hero_id:'rh',hand:[
        {id:'attack-1',type:'attack'},
        {id:'defense-1',type:'defense'},
        {id:'action-1',type:'scheme'},
        {id:'boost-1',type:'versatile'},
      ]},
      {id:'p2',hero_id:'bf'},
    ],
    legal:{
      scheme_cards:['action-1'],
      scheme_actions_by_card:{'action-1':{fighters:[]}},
      attack_cards_by_fighter:{rh:['attack-1','boost-1']},
    },
  };
}

test('pending card model separates card choices from skip choices',()=>{
  const model=ux.pendingCardModel({options:[
    {id:'defend-a',card_id:'a'},
    {id:'defend-b',card_id:'b'},
    {id:'skip',label:'Take damage'},
  ]});
  assert.deepEqual(Array.from(model.cardIDs),['a','b']);
  assert.equal(model.choiceByCard.get('b'),'defend-b');
  assert.deepEqual(Array.from(model.skipOptions,option=>option.id),['skip']);
});

test('legal cards come from server-projected attack, defense, boost, and action domains',()=>{
  const view=baseView();
  const attack=ux.legalCardIDs(view,{mode:'attack_card',values:{fighter_id:'rh'}},{boostOpen:false});
  assert.deepEqual([...attack],['attack-1','boost-1']);

  view.pending={kind:'defense',owner_id:'p1',options:[{id:'d',card_id:'defense-1'},{id:'skip'}]};
  assert.deepEqual([...ux.legalCardIDs(view,{mode:'idle',values:{}},{boostOpen:false})],['defense-1']);

  view.pending={kind:'maneuver_boost',owner_id:'p1',options:[{id:'b',card_id:'boost-1'},{id:'skip'}]};
  assert.deepEqual([...ux.legalCardIDs(view,{mode:'idle',values:{}},{boostOpen:true})],['boost-1']);
  assert.deepEqual([...ux.legalCardIDs(view,{mode:'idle',values:{}},{boostOpen:false})],[]);

  view.pending=null;
  assert.deepEqual([...ux.legalCardIDs(view,{mode:'idle',values:{}},{boostOpen:false})],['action-1']);
  assert.deepEqual([...ux.legalCardIDs(view,{mode:'attack_target',values:{fighter_id:'rh'}},{boostOpen:false})],[]);
  assert.deepEqual([...ux.legalCardIDs(view,{mode:'scheme_destination',values:{card_id:'action-1'}},{boostOpen:false})],['action-1']);
});

test('action plan uses projected fighter and destination domains',()=>{
  const view=baseView();
  assert.deepEqual(JSON.parse(JSON.stringify(ux.actionPlan(view,'action-1'))),{kind:'confirm',cardID:'action-1'});

  view.legal.scheme_actions_by_card['move-card']={fighters:[
    {fighter_id:'rh',destinations:[{destination:'s01'}]},
    {fighter_id:'outlaw',destinations:[{destination:'s02'}]},
  ]};
  assert.deepEqual(JSON.parse(JSON.stringify(ux.actionPlan(view,'move-card'))),{
    kind:'fighter',cardID:'move-card',fighterIDs:['rh','outlaw'],
  });
  assert.deepEqual(JSON.parse(JSON.stringify(ux.actionPlan(view,'move-card','rh'))),{
    kind:'destination',cardID:'move-card',fighterID:'rh',destinations:[{destination:'s01'}],
  });
});
