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

test('Jackalope Horns target domain includes friendly and opposing adjacent fighters only',()=>{
  const jackalope={id:'jackalope',space_id:'s01',defeated:false};
  const fighters=[
    jackalope,
    {id:'bigfoot',space_id:'s03',defeated:false},
    {id:'robin',space_id:'s03',defeated:false},
    {id:'outlaw',space_id:'s10',defeated:false},
    {id:'defeated',space_id:'s03',defeated:true},
  ];
  const adjacency=new Map([['s02',['s01','s03']]]);
  const result=context.adjacentHornsTargets(jackalope,'s02',fighters,adjacency).map(f=>f.id);
  assert.deepEqual(Array.from(result),['bigfoot','robin']);
});
