import assert from 'node:assert/strict';
import {readFile,stat} from 'node:fs/promises';
import test from 'node:test';
import vm from 'node:vm';

const source=await readFile(new URL('./card-art.js',import.meta.url),'utf8');
const context={document:{addEventListener(){}}};
vm.createContext(context);
vm.runInContext(source,context);

const expected={
  'robin-hood':['a-hunters-eye','steal-from-the-rich','disarming-shot','piercing-shot','highway-robbery','defenders-of-sherwood','feint','regroup','wily-fighting','snark','ambush'],
  bigfoot:['larger-than-life','savagery','crash-through-the-trees','jackalope-horns','hoax','disengage','its-just-your-imagination','feint','regroup','skirmish','momentous-shift'],
};

test('Robin Hood and Bigfoot art paths use deck-specific identity',()=>{
  assert.equal(context.cardArtPath({deck_definition_id:'robin-hood',definition_id:'a-hunters-eye'}),'/cards/robin-hood/a-hunters-eye.png');
  assert.equal(context.cardArtPath({deck_definition_id:'bigfoot',definition_id:'larger-than-life'}),'/cards/bigfoot/larger-than-life.png');
  assert.notEqual(context.cardArtPath({deck_definition_id:'robin-hood',definition_id:'feint'}),context.cardArtPath({deck_definition_id:'bigfoot',definition_id:'feint'}));
  assert.notEqual(context.cardArtPath({deck_definition_id:'robin-hood',definition_id:'regroup'}),context.cardArtPath({deck_definition_id:'bigfoot',definition_id:'regroup'}));
});

test('missing identity selects fallback and keeps image markup out',()=>{
  assert.equal(context.cardArtPath({definition_id:'feint'}),null);
  assert.equal(context.cardArtPath({deck_definition_id:'robin-hood'}),null);
  assert.match(context.cardArtMarkup({name:'Fallback'},'<strong>Fallback</strong>'),/no-card-art/);
  assert.doesNotMatch(context.cardArtMarkup({name:'Fallback'},'<strong>Fallback</strong>'),/<img /);
});

test('image markup has accessible async loading contract',()=>{
  const markup=context.cardArtMarkup({id:'instance-7',name:'Feint',deck_definition_id:'robin-hood',definition_id:'feint'},'<strong>Feint</strong>');
  assert.match(markup,/src="\/cards\/robin-hood\/feint\.png"/);
  assert.match(markup,/alt="Feint"/);
  assert.match(markup,/loading="lazy"/);
  assert.match(markup,/decoding="async"/);
  assert.match(markup,/data-card-id="instance-7"/);
});

test('clicking card art resolves the original card instance ID',()=>{
  const picker={dataset:{cardPick:'instance-7'}};
  const image={closest(selector){return selector==='[data-card-pick]'?picker:null}};
  assert.equal(context.cardArtSelectionID(image),'instance-7');
});

test('load hides fallback and error restores it',()=>{
  const classes=new Set();
  const attrs=new Map();
  const fallback={setAttribute(name,value){attrs.set(name,value)},removeAttribute(name){attrs.delete(name)}};
  const shell={classList:{add(name){classes.add(name)},remove(name){classes.delete(name)}},querySelector(){return fallback}};
  const image={matches(selector){return selector==='[data-card-art]'},closest(){return shell},remove(){this.removed=true}};
  context.cardArtEvent({type:'load',target:image});
  assert.equal(classes.has('art-loaded'),true);
  assert.equal(attrs.get('aria-hidden'),'true');
  context.cardArtEvent({type:'error',target:image});
  assert.equal(image.removed,true);
  assert.equal(classes.has('art-loaded'),false);
  assert.equal(attrs.has('aria-hidden'),false);
});

test('visible art cards are limited to the viewing player hand',()=>{
  const view={viewing_player_id:'p1',players:[{id:'p1',hand:[{id:'own',definition_id:'feint'}]},{id:'p2',hand:[{id:'hidden',definition_id:'feint'}]}]};
  assert.deepEqual(context.visibleCardArtCards(view).map(card=>card.id),['own']);
});

test('all expected PNGs exist, are PNGs, and are 250 by 349',async()=>{
  let count=0;
  for(const [deck,ids] of Object.entries(expected))for(const id of ids){
    const path=new URL(`../static/cards/${deck}/${id}.png`,import.meta.url);
    const bytes=await readFile(path);
    const info=await stat(path);
    assert.ok(info.size>0,`${deck}/${id}`);
    assert.deepEqual([...bytes.subarray(0,8)],[137,80,78,71,13,10,26,10],`${deck}/${id} signature`);
    assert.equal(bytes.readUInt32BE(16),250,`${deck}/${id} width`);
    assert.equal(bytes.readUInt32BE(20),349,`${deck}/${id} height`);
    count++;
  }
  assert.equal(count,22);
});

test('source and generated card-art files have identical identity',async()=>{
  assert.equal(await readFile(new URL('./card-art.js',import.meta.url),'utf8'),await readFile(new URL('../static/card-art.js',import.meta.url),'utf8'));
});
