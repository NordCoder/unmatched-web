let stage5PresentationLoad=null;
let stage5PresentationFailure=null;
let stage5AssetLoad=null;
let stage5AssetsReadyKey='';
const stage5FighterTokenHref=BattlefieldRenderer.fighterTokenHref;
const stage5BundledTokenDefinitions=['robin-hood','outlaw','bigfoot','jackalope'];

function stage5LoadingElement(){return $('battlefield-loading')}
function stage5ShowLoading(message='Preparing battlefield…'){
  const overlay=stage5LoadingElement();
  const label=$('battlefield-loading-label');
  if(label)label.textContent=message;
  if(overlay)overlay.hidden=false;
  if(game)game.setAttribute('aria-busy','true');
}
function stage5HideLoading(){
  const overlay=stage5LoadingElement();
  if(overlay)overlay.hidden=true;
  if(game)game.removeAttribute('aria-busy');
}
function stage5BoardDisplayWidth(){
  const frame=document.querySelector?.('.board-frame');
  return Math.max(320,Number(frame?.clientWidth)||Number(globalThis.innerWidth)||1337);
}
function stage5DevicePixelRatio(){return Math.max(1,Number(globalThis.devicePixelRatio)||1)}
function stage5AssetKey(nextView){
  const density=stage5DevicePixelRatio()>=1.45?'high':'standard';
  const width=stage5BoardDisplayWidth()>=900?'wide':'compact';
  return `${nextView?.battlefield_id||'none'}:${density}:${width}`;
}
function stage5ImageFactory(){return typeof Image==='function'?new Image():null}
function stage5PreloadImage(src){
  if(!src||typeof Image!=='function')return Promise.resolve(src);
  return new Promise((resolve,reject)=>{
    const image=new Image();
    let settled=false;
    const done=async()=>{
      if(settled)return;
      try{if(typeof image.decode==='function')await image.decode()}catch(_error){}
      settled=true;resolve(src);
    };
    image.onload=done;
    image.onerror=()=>{if(!settled){settled=true;reject(new Error(`failed to load ${src}`))}};
    image.src=src;
    if(image.complete&&Number(image.naturalWidth)>0)void done();
  });
}
function stage5WithTimeout(promise,milliseconds,label){
  let timeoutID;
  const timeout=new Promise((_,reject)=>{
    timeoutID=setTimeout(()=>reject(new Error(`${label} timed out`)),milliseconds);
  });
  return Promise.race([promise,timeout]).finally(()=>clearTimeout(timeoutID));
}
function stage5ConfigurePresentation(manifest){
  if(!manifest||manifest.battlefield_id!=='sherwood-forest')return manifest;
  const variants=[
    {id:'1x',src:'/battlefields/sherwood-forest/board-1x.webp',pixel_width:1337,pixel_height:742},
    {id:'2x',src:'/battlefields/sherwood-forest/board-2x.webp',pixel_width:2674,pixel_height:1484},
  ];
  manifest.art={
    ...manifest.art,
    src:variants[0].src,
    variants,
  };
  manifest.defaults={
    ...manifest.defaults,
    hero_token_diameter_ratio:.82,
    sidekick_token_diameter_ratio:.77,
    health_badge_ratio:.15,
  };
  return manifest;
}
function stage5VisibleFighters(nextView){
  return nextView?.fighters?.length
    ?nextView.fighters
    :(nextView?.spaces||[]).flatMap(space=>space.fighter?[space.fighter]:[]);
}

function stage5EnsureBattlefieldPresentation(battlefieldID){
  if(!battlefieldID)return Promise.resolve(null);
  const cached=BattlefieldRenderer.get(battlefieldID);
  if(cached)return Promise.resolve(stage5ConfigurePresentation(cached));
  if(stage5PresentationLoad?.battlefieldID===battlefieldID)return stage5PresentationLoad.promise;
  if(stage5PresentationFailure===battlefieldID)return Promise.resolve(null);
  const promise=BattlefieldRenderer.load(battlefieldID)
    .then(manifest=>{
      stage5PresentationFailure=null;
      return stage5ConfigurePresentation(manifest);
    })
    .catch(error=>{
      stage5PresentationFailure=battlefieldID;
      console.warn(`Using fallback battlefield presentation for ${battlefieldID}:`,error);
      return null;
    });
  stage5PresentationLoad={battlefieldID,promise};
  return promise;
}

async function stage5PrepareAssets(nextView){
  if(!nextView?.battlefield_id)return null;
  const key=stage5AssetKey(nextView);
  if(stage5AssetsReadyKey===key)return BattlefieldRenderer.get(nextView.battlefield_id);
  if(stage5AssetLoad?.key===key)return stage5AssetLoad.promise;
  stage5ShowLoading('Preparing Sherwood Forest…');
  const promise=(async()=>{
    const manifest=await stage5EnsureBattlefieldPresentation(nextView.battlefield_id);
    if(manifest){
      await stage5WithTimeout(BattlefieldRenderer.preloadArt(manifest,{
        displayWidth:stage5BoardDisplayWidth(),
        devicePixelRatio:stage5DevicePixelRatio(),
        imageFactory:stage5ImageFactory,
      }),9000,'battlefield art');
    }
    const tokenFighters=[...stage5VisibleFighters(nextView),...stage5BundledTokenDefinitions.map(definition_id=>({definition_id}))];
    const tokens=[...new Set(tokenFighters.map(BattlefieldRenderer.fighterTokenHref).filter(Boolean))];
    await Promise.allSettled(tokens.map(src=>stage5WithTimeout(stage5PreloadImage(src),5000,`token ${src}`)));
    stage5AssetsReadyKey=key;
    if(view?.match_id===nextView.match_id&&view?.battlefield_id===nextView.battlefield_id)renderLocalInteraction();
    return manifest;
  })().catch(error=>{
    console.warn('Battlefield preload failed; continuing with available fallbacks:',error);
    stage5AssetsReadyKey=key;
    if(view?.match_id===nextView.match_id)renderLocalInteraction();
    return BattlefieldRenderer.get(nextView.battlefield_id);
  }).finally(()=>{
    if(stage5AssetLoad?.key===key)stage5AssetLoad=null;
    stage5HideLoading();
  });
  stage5AssetLoad={key,promise};
  return promise;
}

const stage5ApplyView=applyView;
applyView=function(next,options={}){
  const applied=stage5ApplyView(next,options);
  if(next?.battlefield_id&&stage5AssetsReadyKey!==stage5AssetKey(next))void stage5PrepareAssets(next);
  return applied;
};

const stage5EnterGame=enterGame;
enterGame=function(){
  stage5ShowLoading('Loading match…');
  stage5EnterGame();
  if(view?.battlefield_id)void stage5PrepareAssets(view);
};

const stage5Leave=leave;
leave=function(){
  stage5AssetsReadyKey='';
  stage5AssetLoad=null;
  stage5HideLoading();
  stage5Leave();
};

mapDimensions=function(){
  const presentation=BattlefieldRenderer.presentationFor(view);
  return {width:presentation.coordinate_space.width,height:presentation.coordinate_space.height};
};
mapY=function(percent){return Number(percent)/100*mapDimensions().height};
boardPoint=function(spaceID){
  const presentation=BattlefieldRenderer.presentationFor(view);
  const point=BattlefieldRenderer.point(presentation,spaceID);
  return point?`${point.x},${point.y}`:'';
};

renderBoard=function(){
  if(!view)return;
  void stage5EnsureBattlefieldPresentation(view.battlefield_id);
  const svg=$('board');
  const presentation=BattlefieldRenderer.presentationFor(view);
  BattlefieldRenderer.render({
    svg,
    view,
    highlights:boardHighlights(),
    viewerID:view.viewing_player_id,
    initials,
    debug:BattlefieldRenderer.debugEnabled(),
  });
  for(const image of svg.querySelectorAll?.('.board-art')||[]){
    const fallback=presentation.art?.fallback_src;
    if(image&&fallback)image.addEventListener?.('error',()=>image.setAttribute('href',fallback),{once:true});
  }
};

void stage5EnsureBattlefieldPresentation('sherwood-forest');
if(session&&!view)stage5ShowLoading('Loading match…');
