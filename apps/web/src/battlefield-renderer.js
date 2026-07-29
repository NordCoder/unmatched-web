(function(global){
  const cache=new Map();
  const inflight=new Map();
  const DEFAULT_WIDTH=1337;
  const DEFAULT_HEIGHT=742;
  const TOKEN_ROOT='/fighters/tokens';

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
    if(!shape||!['circle','ellipse','polygon'].includes(shape.type))return [`${spaceID}: unsupported or missing shape`];
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

  function artVariants(art){
    const variants=Array.isArray(art?.variants)?art.variants.filter(variant=>variant?.src||variant?.base64_parts?.length):[];
    if(variants.length)return variants;
    return art?.src||art?.base64_parts?.length?[{
      id:'default',src:art.src,base64_parts:art.base64_parts,mime:art.mime,
      pixel_width:art.pixel_width||art.width,pixel_height:art.pixel_height||art.height,
    }]:[];
  }

  function validateManifest(manifest,runtimeSpaceIDs=[]){
    const errors=[];
    if(!manifest||manifest.schema_version!==1)errors.push('schema_version must be 1');
    if(!manifest?.battlefield_id)errors.push('battlefield_id is required');
    const width=number(manifest?.coordinate_space?.width);
    const height=number(manifest?.coordinate_space?.height);
    if(!finite(width)||!finite(height)||width<=0||height<=0)errors.push('coordinate_space width and height must be positive');
    if(!artVariants(manifest?.art).length)errors.push('art source or variants are required');
    const spaces=manifest?.spaces||{};
    const ids=Object.keys(spaces);
    if(!ids.length)errors.push('spaces must not be empty');
    if(width>0&&height>0)for(const id of ids)errors.push(...validateShape(spaces[id]?.shape,width,height,id));
    const runtime=new Set(runtimeSpaceIDs);
    if(runtime.size){
      for(const id of runtime)if(!spaces[id])errors.push(`${id}: presentation geometry is missing`);
      for(const id of ids)if(!runtime.has(id))errors.push(`${id}: presentation geometry has no runtime space`);
    }
    return {ok:errors.length===0,errors};
  }

  function fallbackPresentation(view){
    const width=DEFAULT_WIDTH,height=DEFAULT_HEIGHT,spaces={};
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
      art:{src:'/sherwood-forest.svg',width,height,pixel_width:width,pixel_height:height},
      defaults:{highlight_inset:6,hit_padding:8,hero_token_diameter_ratio:.82,sidekick_token_diameter_ratio:.77},
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

  function selectArtVariant(manifest,options={}){
    const variants=artVariants(manifest?.art);
    if(!variants.length)return null;
    const displayWidth=Math.max(1,number(options.displayWidth)||DEFAULT_WIDTH);
    const dpr=Math.max(1,number(options.devicePixelRatio)||number(global.devicePixelRatio)||1);
    const required=displayWidth*dpr;
    const sorted=[...variants].sort((left,right)=>(number(left.pixel_width)||0)-(number(right.pixel_width)||0));
    return sorted.find(variant=>(number(variant.pixel_width)||0)>=required)||sorted[sorted.length-1];
  }

  async function resolveArtVariant(variant,fetcher=global.fetch){
    if(!variant)return null;
    if(variant.resolved_src)return variant.resolved_src;
    if(!Array.isArray(variant.base64_parts)||!variant.base64_parts.length)return variant.src;
    if(typeof fetcher!=='function')throw new Error('fetch is unavailable');
    const parts=await Promise.all(variant.base64_parts.map(async part=>{
      const response=await fetcher(part,{cache:'force-cache'});
      if(!response?.ok)throw new Error(`Battlefield art part ${part} returned HTTP ${response?.status??'unknown'}`);
      return (await response.text()).trim();
    }));
    const encoded=parts.join('');
    if(encoded.length<1000)throw new Error('Battlefield art parts are unexpectedly empty');
    variant.resolved_src=`data:${variant.mime||'image/webp'};base64,${encoded}`;
    return variant.resolved_src;
  }

  async function prepareArt(manifest,fetcher=global.fetch,options={}){
    const art=manifest?.art;
    if(!art)return null;
    const selected=selectArtVariant(manifest,options);
    if(!selected)return null;
    art.active_variant=selected;
    art.resolved_src=await resolveArtVariant(selected,fetcher);
    return art.resolved_src;
  }

  function imageReady(src,imageFactory){
    if(!src)return Promise.reject(new Error('image source is empty'));
    if(typeof imageFactory!=='function')return Promise.resolve(src);
    return new Promise((resolve,reject)=>{
      const image=imageFactory();
      if(!image)return reject(new Error('image factory returned no image'));
      let settled=false;
      const done=async()=>{
        if(settled)return;
        try{if(typeof image.decode==='function')await image.decode()}catch(_error){}
        settled=true;resolve(src);
      };
      image.onload=done;
      image.onerror=()=>{if(!settled){settled=true;reject(new Error(`failed to load ${src}`))}};
      image.src=src;
      if(image.complete&&number(image.naturalWidth)>0)void done();
    });
  }

  async function preloadArt(manifest,options={}){
    const fetcher=options.fetcher||global.fetch;
    const imageFactory=options.imageFactory||(()=>typeof global.Image==='function'?new global.Image():null);
    const selected=selectArtVariant(manifest,options);
    const variants=artVariants(manifest?.art);
    const ordered=[];
    if(selected)ordered.push(selected);
    for(const variant of [...variants].sort((a,b)=>(number(b.pixel_width)||0)-(number(a.pixel_width)||0))){
      if(!ordered.includes(variant)&&(number(variant.pixel_width)||0)<=(number(selected?.pixel_width)||Infinity))ordered.push(variant);
    }
    let lastError=null;
    for(const variant of ordered){
      try{
        const src=await resolveArtVariant(variant,fetcher);
        await imageReady(src,imageFactory);
        manifest.art.active_variant=variant;
        manifest.art.resolved_src=src;
        manifest.art_error='';
        return src;
      }catch(error){lastError=error}
    }
    const fallback=manifest?.art?.fallback_src;
    if(fallback){
      await imageReady(fallback,imageFactory);
      manifest.art.active_variant={id:'fallback',src:fallback,pixel_width:manifest.coordinate_space.width,pixel_height:manifest.coordinate_space.height};
      manifest.art.resolved_src=fallback;
      manifest.art_error=lastError?.message||'';
      return fallback;
    }
    throw lastError||new Error('battlefield art could not be loaded');
  }

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

  function presentationFor(view){return get(view?.battlefield_id)||fallbackPresentation(view)}
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
  function fighterTokenHref(fighter){
    const definitionID=String(fighter?.definition_id||'').trim();
    if(!/^[a-z0-9-]+$/.test(definitionID))return '';
    return `${TOKEN_ROOT}/${definitionID}.svg`;
  }
  function fighterIsHero(view,fighter){
    return Boolean((view?.players||[]).some(player=>player.hero_id===fighter?.id));
  }

  function renderFighter(fighter,entry,viewerID,initials,view,defaults){
    if(!fighter)return '';
    const center=point({spaces:{space:entry}},'space');
    const base=shapeRadius(entry.shape);
    const hero=fighterIsHero(view,fighter);
    const configured=number(entry.token_diameter_ratio);
    const ratio=configured>0?configured:number(hero?defaults.hero_token_diameter_ratio:defaults.sidekick_token_diameter_ratio)||.78;
    const diameter=Math.max(28,base*2*Math.min(.9,Math.max(.5,ratio)));
    const tokenRadius=diameter/2;
    const healthRadius=Math.max(8,base*(number(defaults.health_badge_ratio)||.15));
    const healthOffset=Math.min(base-healthRadius-2,tokenRadius*.7);
    const healthX=center.x+healthOffset;
    const healthY=center.y+healthOffset;
    const labelY=center.y+tokenRadius*.16;
    const ownerClass=fighter.owner_id===viewerID?'friendly-token':'opposing-token';
    const href=fighterTokenHref(fighter);
    const image=href?`<image class="fighter-art" href="${escapeAttr(href)}" x="${center.x-tokenRadius}" y="${center.y-tokenRadius}" width="${diameter}" height="${diameter}" preserveAspectRatio="xMidYMid meet" pointer-events="none" aria-hidden="true"></image>`:'';
    return `<g class="fighter-piece ${hero?'hero-piece':'sidekick-piece'}" aria-hidden="true">
      <circle class="fighter-token ${ownerClass}" cx="${center.x}" cy="${center.y}" r="${tokenRadius}"></circle>
      <text class="fighter-label" x="${center.x}" y="${labelY}" style="font-size:${Math.max(12,tokenRadius*.55)}px">${escapeHTML(initials(fighter.name))}</text>
      ${image}
      <circle class="health-token" cx="${healthX}" cy="${healthY}" r="${healthRadius}"></circle>
      <text class="health-label" x="${healthX}" y="${healthY+healthRadius*.35}" style="font-size:${Math.max(10,healthRadius*1.05)}px">${fighter.health}</text>
    </g>`;
  }

  function render(options){
    const {svg,view,highlights,viewerID,initials=(name)=>name,debug=debugEnabled()}=options;
    if(!svg||!view)return null;
    const manifest=presentationFor(view);
    const validation=validateManifest(manifest,(view.spaces||[]).map(space=>space.id));
    const effective=validation.ok?manifest:fallbackPresentation(view);
    const width=number(effective.coordinate_space.width),height=number(effective.coordinate_space.height);
    const defaults={highlight_inset:6,hit_padding:8,hero_token_diameter_ratio:.82,sidekick_token_diameter_ratio:.77,health_badge_ratio:.15,...effective.defaults};
    const clips=[],nodes=[];
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
      const highlightShape=adjustShape(entry.shape,-inset),hitShape=adjustShape(entry.shape,padding);
      const clipID=`space-clip-${space.id}`;
      clips.push(`<clipPath id="${clipID}">${shapeMarkup(adjustShape(entry.shape,-1))}</clipPath>`);
      const aria=fighter?`${fighter.name}, ${fighter.health} health, space ${space.id}`:`Space ${space.id}`;
      const debugMarkup=debug?`${shapeMarkup(entry.shape,'class="debug-space-boundary"')}<text class="debug-space-id" x="${shapeCenter(entry.shape).x}" y="${shapeCenter(entry.shape).y}">${escapeHTML(space.id)}</text>`:'';
      nodes.push(`<g class="${classes.join(' ')}" data-space="${escapeAttr(space.id)}"${fighterID?` data-fighter="${escapeAttr(fighterID)}"`:''}${interactive?' role="button" tabindex="0"':''} aria-label="${escapeAttr(aria)}">
        ${shapeMarkup(hitShape,'class="space-hitbox"')}
        <g clip-path="url(#${clipID})">${shapeMarkup(highlightShape,'class="space-highlight"')}</g>
        ${renderFighter(fighter,entry,viewerID,initials,view,defaults)}
        ${debugMarkup}
        <title>${escapeHTML(aria)} · ${escapeHTML((space.zones||[]).join(', '))}</title>
      </g>`);
    }
    const edgeMarkup=debug?(view.edges||[]).map(edge=>{
      const a=point(effective,edge.from),b=point(effective,edge.to);
      return a&&b?`<line class="debug-edge" x1="${a.x}" y1="${a.y}" x2="${b.x}" y2="${b.y}"></line>`:'';
    }).join(''):'';
    const artSource=effective.art.resolved_src||effective.art.active_variant?.resolved_src||effective.art.active_variant?.src||effective.art.src||effective.art.fallback_src;
    svg.setAttribute('viewBox',`0 0 ${width} ${height}`);
    svg.setAttribute('preserveAspectRatio','xMidYMid meet');
    svg.style.aspectRatio=`${width}/${height}`;
    svg.dataset.presentation=effective.fallback?'fallback':'calibrated';
    svg.innerHTML=`<defs>${clips.join('')}</defs>
      <image class="board-art" href="${escapeAttr(artSource)}" x="0" y="0" width="${width}" height="${height}" preserveAspectRatio="xMidYMid meet"></image>
      <g class="debug-edges">${edgeMarkup}</g>
      <polyline id="path-preview" class="path-preview" points=""></polyline>
      <g class="board-overlay">${nodes.join('')}</g>`;
    return {manifest:effective,validation};
  }

  global.BattlefieldRenderer={
    manifestPath,validateManifest,fallbackPresentation,register,get,load,
    artVariants,selectArtVariant,resolveArtVariant,prepareArt,preloadArt,
    presentationFor,geometry,point,shapeCenter,shapeRadius,adjustShape,shapeMarkup,
    debugEnabled,fighterTokenHref,fighterIsHero,render,
  };
})(globalThis);
