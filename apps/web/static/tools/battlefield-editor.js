const editor={
  manifest:null,
  selectedID:'',
  dragging:false,
  dragOffset:{x:0,y:0},
};

const editor$=id=>document.getElementById(id);
const editorSVG=editor$('editor-board');
const status=message=>{editor$('status').textContent=message};

function editorSpaceIDs(){
  return Object.keys(editor.manifest?.spaces||{}).sort((a,b)=>a.localeCompare(b,undefined,{numeric:true}));
}

function editorPoint(event){
  const point=editorSVG.createSVGPoint();
  point.x=event.clientX;
  point.y=event.clientY;
  return point.matrixTransform(editorSVG.getScreenCTM().inverse());
}

function editorSelectedEntry(){
  return editor.manifest?.spaces?.[editor.selectedID]||null;
}

function editorSetCenter(shape,x,y){
  const center=BattlefieldRenderer.shapeCenter(shape);
  if(!center)return;
  const dx=x-center.x,dy=y-center.y;
  if(shape.type==='circle'||shape.type==='ellipse'){
    shape.cx=Number(x.toFixed(1));
    shape.cy=Number(y.toFixed(1));
  }else if(shape.type==='polygon'){
    shape.points=shape.points.map(point=>[
      Number((Number(point[0])+dx).toFixed(1)),
      Number((Number(point[1])+dy).toFixed(1)),
    ]);
  }
}

function editorSetRadius(shape,radius){
  const current=BattlefieldRenderer.shapeRadius(shape);
  if(!current||!Number.isFinite(radius)||radius<=0)return;
  const adjusted=BattlefieldRenderer.adjustShape(shape,radius-current);
  Object.assign(shape,adjusted);
}

function editorRefreshControls(){
  const ids=editorSpaceIDs();
  const select=editor$('space-select');
  select.innerHTML=ids.map(id=>`<option value="${id}">${id}</option>`).join('');
  if(!ids.includes(editor.selectedID))editor.selectedID=ids[0]||'';
  select.value=editor.selectedID;
  const entry=editorSelectedEntry();
  const center=BattlefieldRenderer.shapeCenter(entry?.shape);
  const radius=BattlefieldRenderer.shapeRadius(entry?.shape);
  editor$('space-x').value=center?.x??'';
  editor$('space-y').value=center?.y??'';
  editor$('space-r').value=radius?radius.toFixed(1):'';
  const circleLike=['circle','ellipse'].includes(entry?.shape?.type);
  editor$('space-x').disabled=!entry;
  editor$('space-y').disabled=!entry;
  editor$('space-r').disabled=!entry||(!circleLike&&entry?.shape?.type!=='polygon');
}

function editorRender(){
  if(!editor.manifest)return;
  const width=Number(editor.manifest.coordinate_space.width);
  const height=Number(editor.manifest.coordinate_space.height);
  const zoom=Number(editor$('zoom').value)/100;
  editorSVG.setAttribute('viewBox',`0 0 ${width} ${height}`);
  editorSVG.setAttribute('preserveAspectRatio','xMidYMid meet');
  editorSVG.style.width=`${width*zoom}px`;
  editorSVG.style.height=`${height*zoom}px`;

  const showIDs=editor$('show-ids').checked;
  const showHitboxes=editor$('show-hitboxes').checked;
  const hitPadding=Number(editor.manifest.defaults?.hit_padding??8);
  const nodes=editorSpaceIDs().map(id=>{
    const entry=editor.manifest.spaces[id];
    const shape=entry.shape;
    const center=BattlefieldRenderer.shapeCenter(shape);
    const selected=id===editor.selectedID?' selected':'';
    const hitbox=showHitboxes
      ?BattlefieldRenderer.shapeMarkup(BattlefieldRenderer.adjustShape(shape,hitPadding),'class="editor-hitbox"')
      :'';
    const label=showIDs?`<text class="editor-label" x="${center.x}" y="${center.y}">${id}</text>`:'';
    return `<g data-space="${id}">
      ${hitbox}
      ${BattlefieldRenderer.shapeMarkup(shape,`class="editor-space${selected}" data-space="${id}"`)}
      <circle class="editor-center" cx="${center.x}" cy="${center.y}" r="3"></circle>
      ${label}
    </g>`;
  }).join('');

  editorSVG.innerHTML=`<image href="${editor.manifest.art.src}" x="0" y="0" width="${width}" height="${height}" preserveAspectRatio="xMidYMid meet"></image>${nodes}`;
  editorRefreshControls();
}

async function editorLoad(battlefieldID){
  status(`Loading ${battlefieldID}…`);
  try{
    const response=await fetch(BattlefieldRenderer.manifestPath(battlefieldID),{cache:'no-store'});
    if(!response.ok)throw new Error(`HTTP ${response.status}`);
    const manifest=await response.json();
    const validation=BattlefieldRenderer.validateManifest(manifest);
    if(!validation.ok)throw new Error(validation.errors.join('\n'));
    editor.manifest=manifest;
    editor.selectedID=editorSpaceIDs()[0]||'';
    editorRender();
    status(`${Object.keys(manifest.spaces).length} spaces loaded.`);
  }catch(error){
    status(`Load failed: ${error.message}`);
  }
}

function editorSelect(id){
  if(!editor.manifest?.spaces?.[id])return;
  editor.selectedID=id;
  editorRender();
}

function editorMoveSelection(delta){
  const ids=editorSpaceIDs();
  const index=ids.indexOf(editor.selectedID);
  if(index<0)return;
  editorSelect(ids[(index+delta+ids.length)%ids.length]);
}

function editorUpdateFromInputs(){
  const entry=editorSelectedEntry();
  if(!entry)return;
  const x=Number(editor$('space-x').value);
  const y=Number(editor$('space-y').value);
  const radius=Number(editor$('space-r').value);
  if(Number.isFinite(x)&&Number.isFinite(y))editorSetCenter(entry.shape,x,y);
  if(Number.isFinite(radius)&&radius>0)editorSetRadius(entry.shape,radius);
  editorRender();
}

function editorExport(){
  if(!editor.manifest)return;
  const validation=BattlefieldRenderer.validateManifest(editor.manifest);
  if(!validation.ok){
    status(`Cannot export:\n${validation.errors.join('\n')}`);
    return;
  }
  const blob=new Blob([`${JSON.stringify(editor.manifest,null,2)}\n`],{type:'application/json'});
  const url=URL.createObjectURL(blob);
  const link=document.createElement('a');
  link.href=url;
  link.download='manifest.json';
  link.click();
  URL.revokeObjectURL(url);
  status('manifest.json exported.');
}

editorSVG.addEventListener('pointerdown',event=>{
  const target=event.target.closest?.('[data-space]');
  const id=target?.dataset?.space;
  if(!id)return;
  editorSelect(id);
  const center=BattlefieldRenderer.shapeCenter(editorSelectedEntry().shape);
  const point=editorPoint(event);
  editor.dragOffset={x:center.x-point.x,y:center.y-point.y};
  editor.dragging=true;
  editorSVG.setPointerCapture?.(event.pointerId);
  event.preventDefault();
});

editorSVG.addEventListener('pointermove',event=>{
  if(!editor.dragging)return;
  const point=editorPoint(event);
  editorSetCenter(editorSelectedEntry().shape,point.x+editor.dragOffset.x,point.y+editor.dragOffset.y);
  editorRender();
});

function editorEndDrag(event){
  if(!editor.dragging)return;
  editor.dragging=false;
  editorSVG.releasePointerCapture?.(event.pointerId);
}
editorSVG.addEventListener('pointerup',editorEndDrag);
editorSVG.addEventListener('pointercancel',editorEndDrag);

editorSVG.addEventListener('wheel',event=>{
  const entry=editorSelectedEntry();
  if(!entry)return;
  event.preventDefault();
  const radius=BattlefieldRenderer.shapeRadius(entry.shape);
  editorSetRadius(entry.shape,Math.max(1,radius+(event.deltaY<0?1:-1)));
  editorRender();
},{passive:false});

editor$('space-select').onchange=event=>editorSelect(event.target.value);
for(const id of ['space-x','space-y','space-r'])editor$(id).onchange=editorUpdateFromInputs;
editor$('zoom').oninput=editorRender;
editor$('show-ids').onchange=editorRender;
editor$('show-hitboxes').onchange=editorRender;
editor$('previous').onclick=()=>editorMoveSelection(-1);
editor$('next').onclick=()=>editorMoveSelection(1);
editor$('load').onclick=()=>editorLoad(editor$('battlefield-id').value.trim());
editor$('export').onclick=editorExport;
editor$('import').onchange=async event=>{
  const file=event.target.files?.[0];
  if(!file)return;
  try{
    const manifest=JSON.parse(await file.text());
    const validation=BattlefieldRenderer.validateManifest(manifest);
    if(!validation.ok)throw new Error(validation.errors.join('\n'));
    editor.manifest=manifest;
    editor$('battlefield-id').value=manifest.battlefield_id;
    editor.selectedID=editorSpaceIDs()[0]||'';
    editorRender();
    status(`${Object.keys(manifest.spaces).length} spaces imported.`);
  }catch(error){status(`Import failed: ${error.message}`)}
};

void editorLoad(editor$('battlefield-id').value);
