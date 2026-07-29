(function(global){
  const cache=new Map();
  const inflight=new Map();
  const DEFAULT_WIDTH=1337;
  const DEFAULT_HEIGHT=742;

  function finite(value){return Number.isFinite(Number(value))}
  function number(value){return Number(value)}
  function escapeHTML(value){
    return String(value??'').replace(/[&<>"]/g,ch=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;'}[ch]));
  }
  function escapeAttr(value){return escapeHTML(value).replace(/'/g,'&#39;')}

  function manifestPath(battlefieldID){
    return `/battlefields/${encodeURIComponent(battlefieldID)}/manifest.json`;
  }

  function shapeCenter(shape){
    if(!shape)return null;
    if(shape.type==='circle'||shape.type==='ellipse')return {x:number(shape.cx),y:number(shape.cy)};
    if(shape.type==='polygon'){
      const points=shape.points||[];
      if(!points.length)return null;
      const total=points.reduce((sum,point)=>({x:sum.x+number(point[0]),y:sum.y+number(point[1])}),{x:0,y:0});
      return {x:total.x/points.length,y:total.y/points.length};
    }
    return null;
  }

  function shapeRadius(shape){
    if(!shape)return 0;
    if(shape.type==='circle')return number(shape.r);
    if(shape.type==='ellipse')return Math.min(number(shape.rx),number(shape.ry));
    if(shape.type==='polygon'){
      const center=shapeCenter(shape);
      return Math.min(...shape.points.map(point=>Math.hypot(number(point[0])-center.x,number(point[1])-center.y)));
    }
    return 0;
  }

  function adjustShape(shape,delta){
    if(shape.type==='circle')return {...shape,r:Math.max(1,number(shape.r)+delta)};
    if(shape.type==='ellipse')return {...shape,rx:Math.max(1,number(shape.rx)+delta),ry:Math.max(1,number(shape.ry)+delta)};
    if(shape.type==='polygon'){
      const center=shapeCenter(shape);
      const scale=Math.max(.05,(shapeRadius(shape)+delta)/Math.max(1,shapeRadius(shape)));
      return {...shape,points:shape.points.map(point=>[
        center.x+(number(point[0])-center.x)*scale,
        center.y+(number(point[1])-center.y)*scale,
      ])};
    }
    return shape;
  }

  function shapeMarkup(shape,attributes=''){
    if(shape.type==='circle')return `<circle ${attributes} cx="${number(shape.cx)}" cy="${number(shape.cy)}" r="${number(shape.r)}"></circle>`;
    if(shape.type==='ellipse')return `<ellipse ${attributes} cx="${number(shape.cx)}" cy="${number(shape.cy)}" rx="${number(shape.rx)}" ry="${number(shape.ry)}"></ellipse>`;
    if(shape.type==='polygon')return `<polygon ${attributes} points="${shape.points.map(point=>`${number(point[0])},${number(point[1])}`).join(' ')}"></polygon>`;
    return '';
  }

  function validateShape(shape,width,height,spaceID){
    const errors=[];
    if(!shape||!['circle','ellipse','polygon'].includes(shape.type)){
      return [`${spaceID}: unsupported or missing shape`];
    }
    if(shape.type==='circle'){
      if(!finite(shape.cx)||!finite(shape.cy)||!finite(shape.r)||number(shape.r)<=0)errors.push(`${spaceID}: invalid circle`);
      if(number(shape.cx)-number(shape.r)<0||number(shape.cy)-number(shape.r)<0||number(shape.cx)+number(shape.r)>width||number(shape.cy)+number(shape.r)>height)errors.push(`${spaceID}: circle exceeds coordinate space`);
    }else if(shape.type==='ellipse'){
      if(!finite(shape.cx)||!finite(shape.cy)||!finite(shape.rx)||!finite(shape.ry)||number(shape.rx)<=0||number(shape.ry)<=0)errors.push(`${spaceID}: invalid ellipse`);
      if(number(shape.cx)-number(shape.rx)<0||number(shape.cy)-number(shape.ry)<0||number(shape.cx)+number(shape.rx)>width||number(shape.cy)+number(shape.ry)>height)errors.push(`${spaceID}: ellipse exceeds coordinate space`);
    }else{
      if(!Array.isArray(shape.points)||shape.points.length<3)errors.push(`${spaceID}: polygon needs at least three points`);
      for(const point of shape.points||[]){
        if(!Array.isArray(point)||point.length!==2||!finite(point[0])||!finite(point[1]))errors.push(`${spaceID}: invalid polygon point`);
        else if(number(point[0])<0||number(point[1])<0||number(point[0])>width||number(point[1])>height)errors.push(`${spaceID}: polygon point exceeds coordinate space`);
      }
    }
    return errors;
  }

  function validateManifest(manifest,runtimeSpaceIDs=[]){
    const errors=[];
    if(!manifest||manifest.schema_version!==1)errors.push('schema_version must be 1');
    if(!manifest?.battlefield_id)errors.push('battlefield_id is required');
    const width=number(manifest?.coordinate_space?.width);
    const height=number(manifest?.coordinate_space?.height);
    if(!finite(width)||!finite(height)||width<=0||height<=0)errors.push('coordinate_space width and height must be positive');
    if(!manifest?.art?.src)errors.push('art.src is required');
    const spaces=manifest?.spaces||{};
    const ids=Object.keys(spaces);
    if(!ids.length)errors.push('spaces must not be empty');
    if(width>0&&height>0){
      for(const id of ids)errors.push(...validateShape(spaces[id]?.shape,width,height,id));
    }
    const runtime=new Set(runtimeSpaceIDs);
    if(runtime.size){
      for(const id of runtime)if(!spaces[id])errors.push(`${id}: presentation geometry is missing`);
      for(const id of ids)if(!runtime.has(id))errors.push(`${id}: presentation geometry has no runtime space`);
    }
    return {ok:errors.length===0,errors};
  }

  function fallbackPresentation(view){
    const width=DEFAULT_WIDTH;
    const height=DEFAULT_HEIGHT;
    const spaces={};
    for(const space of view?.spaces||[]){
      spaces[space.id]={
        shape:{type:'circle',cx:number(space.x)/100*width,cy:number(space.y)/100*height,r:54},
        token_anchor:{x:number(space.x)/100*width,y:number(space.y)/100*height},
      };
    }
    return {
      schema_version:1,
      battlefield_id:view?.battlefield_id||'fallback',
      coordinate_space:{width,height},
      art:{src:'/sherwood-forest.svg',width,height},
      defaults:{highlight_inset:6,hit_padding:8,token_scale:.36},
      spaces,
      fallback:true,
    };
  }

  function register(manifest){
    const validation=validateManifest(manifest);
    if(!validation.ok)throw new Error(`Invalid battlefield presentation: ${validation.errors.join('; ')}`);
    cache.set(manifest.battlefield_id,manifest);
    return manifest;
  }

  function get(battlefieldID){return cache.get(battlefieldID)||null}

  async function load(battlefieldID,fetcher=global.fetch){
    if(cache.has(battlefieldID))return cache.get(battlefieldID);
    if(inflight.has(battlefieldID))return inflight.get(battlefieldID);
    if(typeof fetcher!=='function')throw new Error('fetch is unavailable');
    const promise=(async()=>{
      const response=await fetcher(manifestPath(battlefieldID),{cache:'force-cache'});
      if(!response?.ok)throw new Error(`Battlefield presentation ${battlefieldID} returned HTTP ${response?.status??'unknown'}`);
      const manifest=await response.json();
      if(manifest.battlefield_id!==battlefieldID)throw new Error(`Battlefield presentation identity mismatch: ${manifest.battlefield_id}`);
      return register(manifest);
    })().finally(()=>inflight.delete(battlefieldID));
    inflight.set(battlefieldID,promise);
    return promise;
  }

  function presentationFor(view){
    return get(view?.battlefield_id)||fallbackPresentation(view);
  }

  function geometry(manifest,spaceID){return manifest?.spaces?.[spaceID]||null}

  function point(manifest,spaceID){
    const entry=geometry(manifest,spaceID);
    if(!entry)return null;
    const anchor=entry.token_anchor;
    if(anchor&&finite(anchor.x)&&finite(anchor.y))return {x:number(anchor.x),y:number(anchor.y)};
    return shapeCenter(entry.shape);
  }

  function debugEnabled(search=global.location?.search||''){
    return new URLSearchParams(search).get('battlefieldDebug')==='1';
  }

  function renderFighter(fighter,entry,viewerID,initials){
    if(!fighter)return '';
    const center=point({spaces:{space:entry}},'space');
    const base=shapeRadius(entry.shape);
    const scale=number(entry.token_scale)||.36;
    const tokenRadius=Math.max(12,base*scale);
    const healthRadius=Math.max(7,tokenRadius*.34);
    const healthX=center.x+tokenRadius*.62;
    const healthY=center.y+tokenRadius*.62;
    const labelY=center.y+tokenRadius*.18;
    const ownerClass=fighter.owner_id===viewerID?'friendly-token':'opposing-token';
    return `<g class="fighter-piece" aria-hidden="true">
      <circle class="fighter-token ${ownerClass}" cx="${center.x}" cy="${center.y}" r="${tokenRadius}"></circle>
      <text class="fighter-label" x="${center.x}" y="${labelY}" style="font-size:${Math.max(10,tokenRadius*.62)}px">${escapeHTML(initials(fighter.name))}</text>
      <circle class="health-token" cx="${healthX}" cy="${healthY}" r="${healthRadius}"></circle>
      <text class="health-label" x="${healthX}" y="${healthY+healthRadius*.35}" style="font-size:${Math.max(9,healthRadius*1.05)}px">${fighter.health}</text>
    </g>`;
  }

  function render(options){
    const {svg,view,highlights,viewerID,initials=(name)=>name,debug=debugEnabled()}=options;
    if(!svg||!view)return null;
    const manifest=presentationFor(view);
    const validation=validateManifest(manifest,(view.spaces||[]).map(space=>space.id));
    const effective=validation.ok?manifest:fallbackPresentation(view);
    const width=number(effective.coordinate_space.width);
    const height=number(effective.coordinate_space.height);
    const defaults={highlight_inset:6,hit_padding:8,token_scale:.36,...effective.defaults};
    const clips=[];
    const nodes=[];
    for(const space of view.spaces||[]){
      const entry=effective.spaces[space.id];
      if(!entry)continue;
      const fighter=space.fighter;
      const fighterID=fighter?.id||'';
      const classes=['board-space'];
      if(highlights.destinations.has(space.id))classes.push('legal-destination');
      if(fighterID&&highlights.fighterCandidates.has(fighterID))classes.push('legal-fighter');
      if(fighterID&&highlights.targets.has(fighterID))classes.push('legal-target');
      if(fighterID&&fighterID===highlights.selectedFighterID)classes.push('selected-fighter');
      if(fighterID&&fighterID===highlights.selectedTargetID)classes.push('selected-target');
      const interactive=classes.length>1||(fighter&&fighter.owner_id===viewerID&&!fighter.defeated);
      const inset=finite(entry.highlight_inset)?number(entry.highlight_inset):number(defaults.highlight_inset);
      const padding=finite(entry.hit_padding)?number(entry.hit_padding):number(defaults.hit_padding);
      const highlightShape=adjustShape(entry.shape,-inset);
      const hitShape=adjustShape(entry.shape,padding);
      const clipID=`space-clip-${space.id}`;
      clips.push(`<clipPath id="${clipID}">${shapeMarkup(adjustShape(entry.shape,-1))}</clipPath>`);
      const aria=fighter?`${fighter.name}, ${fighter.health} health, space ${space.id}`:`Space ${space.id}`;
      const entryWithScale={...entry,token_scale:entry.token_scale??defaults.token_scale};
      const debugMarkup=debug?`${shapeMarkup(entry.shape,'class="debug-space-boundary"')}<text class="debug-space-id" x="${shapeCenter(entry.shape).x}" y="${shapeCenter(entry.shape).y}">${escapeHTML(space.id)}</text>`:'';
      nodes.push(`<g class="${classes.join(' ')}" data-space="${escapeAttr(space.id)}"${fighterID?` data-fighter="${escapeAttr(fighterID)}"`:''}${interactive?' role="button" tabindex="0"':''} aria-label="${escapeAttr(aria)}">
        ${shapeMarkup(hitShape,'class="space-hitbox"')}
        <g clip-path="url(#${clipID})">${shapeMarkup(highlightShape,'class="space-highlight"')}</g>
        ${renderFighter(fighter,entryWithScale,viewerID,initials)}
        ${debugMarkup}
        <title>${escapeHTML(aria)} · ${escapeHTML((space.zones||[]).join(', '))}</title>
      </g>`);
    }
    const edgeMarkup=debug?(view.edges||[]).map(edge=>{
      const a=point(effective,edge.from),b=point(effective,edge.to);
      return a&&b?`<line class="debug-edge" x1="${a.x}" y1="${a.y}" x2="${b.x}" y2="${b.y}"></line>`:'';
    }).join(''):'';
    svg.setAttribute('viewBox',`0 0 ${width} ${height}`);
    svg.setAttribute('preserveAspectRatio','xMidYMid meet');
    svg.style.aspectRatio=`${width}/${height}`;
    svg.dataset.presentation=effective.fallback?'fallback':'calibrated';
    svg.innerHTML=`<defs>${clips.join('')}</defs>
      <image class="board-art" href="${escapeAttr(effective.art.src)}" x="0" y="0" width="${width}" height="${height}" preserveAspectRatio="xMidYMid meet"></image>
      <g class="debug-edges">${edgeMarkup}</g>
      <polyline id="path-preview" class="path-preview" points=""></polyline>
      <g class="board-overlay">${nodes.join('')}</g>`;
    return {manifest:effective,validation};
  }

  global.BattlefieldRenderer={
    manifestPath,
    validateManifest,
    fallbackPresentation,
    register,
    get,
    load,
    presentationFor,
    geometry,
    point,
    shapeCenter,
    shapeRadius,
    adjustShape,
    shapeMarkup,
    debugEnabled,
    render,
  };
})(globalThis);
