const STAGE5_TOKEN_NAMESPACE='http://www.w3.org/2000/svg';
const STAGE5_TOKEN_ROOT='/fighters/tokens';

let stage5PresentationLoad=null;
let stage5PresentationFailure=null;

function stage5FighterTokenHref(fighter){
  const definitionID=String(fighter?.definition_id||'').trim();
  if(!/^[a-z0-9-]+$/.test(definitionID))return '';
  return `${STAGE5_TOKEN_ROOT}/${definitionID}.svg`;
}

function stage5ApplyFighterTokenArt(svg,fighters=[]){
  if(!svg?.querySelectorAll)return 0;
  const byID=new Map((fighters||[]).filter(Boolean).map(fighter=>[fighter.id,fighter]));
  const documentRef=svg.ownerDocument||globalThis.document;
  if(!documentRef?.createElementNS)return 0;
  let rendered=0;
  for(const space of svg.querySelectorAll('.board-space[data-fighter]')){
    const fighterID=space.getAttribute?.('data-fighter')||space.dataset?.fighter||'';
    const fighter=byID.get(fighterID);
    const href=stage5FighterTokenHref(fighter);
    const piece=space.querySelector?.('.fighter-piece');
    const token=piece?.querySelector?.('.fighter-token');
    if(!href||!piece||!token)continue;
    const cx=Number(token.getAttribute('cx'));
    const cy=Number(token.getAttribute('cy'));
    const radius=Number(token.getAttribute('r'));
    if(!Number.isFinite(cx)||!Number.isFinite(cy)||!Number.isFinite(radius)||radius<=0)continue;
    const diameter=radius*2.18;
    const art=documentRef.createElementNS(STAGE5_TOKEN_NAMESPACE,'image');
    art.setAttribute('class','fighter-art');
    art.setAttribute('href',href);
    art.setAttribute('x',String(cx-diameter/2));
    art.setAttribute('y',String(cy-diameter/2));
    art.setAttribute('width',String(diameter));
    art.setAttribute('height',String(diameter));
    art.setAttribute('preserveAspectRatio','xMidYMid meet');
    art.setAttribute('pointer-events','none');
    art.setAttribute('aria-hidden','true');
    const health=piece.querySelector?.('.health-token');
    piece.insertBefore(art,health||null);
    rendered++;
  }
  return rendered;
}

function stage5EnsureBattlefieldPresentation(battlefieldID){
  if(!battlefieldID||BattlefieldRenderer.get(battlefieldID))return Promise.resolve(BattlefieldRenderer.get(battlefieldID));
  if(stage5PresentationLoad?.battlefieldID===battlefieldID)return stage5PresentationLoad.promise;
  if(stage5PresentationFailure===battlefieldID)return Promise.resolve(null);
  const promise=BattlefieldRenderer.load(battlefieldID)
    .then(manifest=>{
      stage5PresentationFailure=null;
      if(view?.battlefield_id===battlefieldID)renderLocalInteraction();
      return manifest;
    })
    .catch(error=>{
      stage5PresentationFailure=battlefieldID;
      console.warn(`Using fallback battlefield presentation for ${battlefieldID}:`,error);
      return null;
    });
  stage5PresentationLoad={battlefieldID,promise};
  return promise;
}

const stage5ApplyView=applyView;
applyView=function(next,options={}){
  const applied=stage5ApplyView(next,options);
  if(next?.battlefield_id)void stage5EnsureBattlefieldPresentation(next.battlefield_id);
  return applied;
};

mapDimensions=function(){
  const presentation=BattlefieldRenderer.presentationFor(view);
  return {
    width:presentation.coordinate_space.width,
    height:presentation.coordinate_space.height,
  };
};

mapY=function(percent){
  return Number(percent)/100*mapDimensions().height;
};

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
  const projectedFighters=view.fighters?.length
    ?view.fighters
    :(view.spaces||[]).flatMap(space=>space.fighter?[space.fighter]:[]);
  stage5ApplyFighterTokenArt(svg,projectedFighters);
  const image=svg.querySelector?.('.board-art');
  const fallback=presentation.art?.fallback_src;
  if(image&&fallback){
    image.addEventListener?.('error',()=>image.setAttribute('href',fallback),{once:true});
  }
};

void stage5EnsureBattlefieldPresentation('sherwood-forest');
