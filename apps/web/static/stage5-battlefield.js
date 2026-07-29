let stage5PresentationLoad=null;
let stage5PresentationFailure=null;

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
  BattlefieldRenderer.render({
    svg:$('board'),
    view,
    highlights:boardHighlights(),
    viewerID:view.viewing_player_id,
    initials,
    debug:BattlefieldRenderer.debugEnabled(),
  });
};

void stage5EnsureBattlefieldPresentation('sherwood-forest');
