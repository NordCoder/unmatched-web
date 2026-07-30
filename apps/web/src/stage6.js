(function(global){
  const state={cardMode:'',selectedCardID:'',confirmation:null,boostOpen:false};

  function ownHand(currentView=typeof view==='undefined'?null:view){
    return currentView?.players?.find(player=>player.id===currentView.viewing_player_id)?.hand||[];
  }
  function ownTurn(currentView=typeof view==='undefined'?null:view){
    return Boolean(currentView&&currentView.phase==='active'&&!currentView.pending&&currentView.current_player_id===currentView.viewing_player_id);
  }
  function optionCardMap(pending){
    const result=new Map();
    for(const option of pending?.options||[])if(option.card_id)result.set(option.card_id,option.id);
    return result;
  }
  function pendingCardModel(pending){
    const cards=optionCardMap(pending);
    return {
      cardIDs:[...cards.keys()],
      choiceByCard:cards,
      skipOptions:(pending?.options||[]).filter(option=>!option.card_id),
    };
  }
  function legalCardIDs(currentView=typeof view==='undefined'?null:view,currentInteraction=typeof interaction==='undefined'?null:interaction,currentState=state){
    if(!currentView)return new Set();
    const pending=currentView.pending;
    if(pending?.owner_id===currentView.viewing_player_id){
      if(pending.kind==='defense')return new Set(pendingCardModel(pending).cardIDs);
      if(pending.kind==='maneuver_boost'&&currentState.boostOpen)return new Set(pendingCardModel(pending).cardIDs);
      return new Set();
    }
    if(currentInteraction?.mode==='attack_card'){
      return new Set(currentView.legal?.attack_cards_by_fighter?.[currentInteraction.values?.fighter_id]||[]);
    }
    if(['scheme_fighter','scheme_destination','scheme_target','action_confirm'].includes(currentInteraction?.mode)){
      const cardID=currentInteraction?.values?.card_id;
      return cardID?new Set([cardID]):new Set();
    }
    if(currentInteraction&&!['idle','fighter_selected'].includes(currentInteraction.mode))return new Set();
    if(ownTurn(currentView))return new Set(currentView.legal?.scheme_cards||[]);
    return new Set();
  }
  function cardRole(card,currentView=typeof view==='undefined'?null:view,currentInteraction=typeof interaction==='undefined'?null:interaction,currentState=state){
    if(!card)return '';
    const pending=currentView?.pending;
    if(pending?.owner_id===currentView.viewing_player_id&&pending.kind==='defense')return 'defense';
    if(pending?.owner_id===currentView.viewing_player_id&&pending.kind==='maneuver_boost'&&currentState.boostOpen)return 'boost';
    if(currentInteraction?.mode==='attack_card')return 'attack';
    if(ownTurn(currentView)&&currentView.legal?.scheme_cards?.includes(card.id))return 'action';
    return '';
  }
  function actionPlan(currentView,cardID,selectedFighterID=''){
    const action=currentView?.legal?.scheme_actions_by_card?.[cardID];
    const fighters=action?.fighters||[];
    if(!fighters.length)return {kind:'confirm',cardID};
    const selected=fighters.find(candidate=>candidate.fighter_id===selectedFighterID);
    const candidate=selected||(fighters.length===1?fighters[0]:null);
    if(!candidate)return {kind:'fighter',cardID,fighterIDs:fighters.map(item=>item.fighter_id)};
    return {kind:'destination',cardID,fighterID:candidate.fighter_id,destinations:candidate.destinations||[]};
  }
  function confirmation(kind,summary,payload,extra={}){
    state.confirmation={kind,summary,payload,...extra};
    return state.confirmation;
  }
  function clearConfirmation(){state.confirmation=null;state.selectedCardID=''}

  function install(){
    if(typeof document==='undefined'||typeof $!=='function')return;

    function ensureConfirmationBar(){
      let bar=document.getElementById('card-confirmation');
      if(bar)return bar;
      bar=document.createElement('div');
      bar.id='card-confirmation';
      bar.className='card-confirmation';
      bar.hidden=true;
      const section=$('hand')?.closest?.('section');
      (section||$('game')||document.body).append(bar);
      return bar;
    }
    function renderConfirmation(){
      const bar=ensureConfirmationBar();
      const item=state.confirmation;
      if(!item){bar.hidden=true;bar.innerHTML='';return}
      bar.hidden=false;
      bar.innerHTML=`<div><strong>${escapeHTML(item.summary)}</strong><span>${escapeHTML(item.note||'Confirm the selected card and choices.')}</span></div><div class="card-confirmation-actions"><button id="stage6-confirm">${escapeHTML(item.confirmLabel||'Confirm')}</button><button id="stage6-cancel" class="ghost">Cancel</button></div>`;
      $('stage6-confirm').onclick=async()=>{
        const current=state.confirmation;
        if(!current)return;
        const accepted=await current.execute();
        if(accepted!==false){clearConfirmation();state.cardMode='';state.boostOpen=false;renderConfirmation()}
      };
      $('stage6-cancel').onclick=()=>{
        const previous=state.confirmation;
        clearConfirmation();
        if(previous?.cancel)previous.cancel();
        renderHand();renderLocalInteraction();renderConfirmation();
      };
    }

    function selectedCardSummary(card,prefix){return `${prefix}: ${card?.name||'selected card'}`}
    function selectAttackCard(card){
      if(!legalCardIDs().has(card.id))return;
      state.selectedCardID=card.id;
      state.cardMode='attack';
      confirmation('attack',selectedCardSummary(card,'Attack'),null,{
        confirmLabel:'Attack',
        execute:()=>command({type:'attack',fighter_id:interactionValue('fighter_id'),target_id:interactionValue('target_id'),card_id:card.id}),
        cancel:()=>{beginInteraction('attack_target','attack',{fighter_id:interactionValue('fighter_id')})},
      });
      renderHand();renderConfirmation();
    }
    function selectDefenseCard(card){
      const model=pendingCardModel(view.pending);
      const choice=model.choiceByCard.get(card.id);
      if(!choice)return;
      state.selectedCardID=card.id;
      state.cardMode='defense';
      confirmation('defense',selectedCardSummary(card,'Defense'),null,{
        confirmLabel:'Defend',
        execute:()=>command({type:'choose',choice}),
      });
      renderHand();renderConfirmation();
    }
    function selectBoostCard(card){
      const model=pendingCardModel(view.pending);
      const choice=model.choiceByCard.get(card.id);
      if(!choice)return;
      state.selectedCardID=card.id;
      state.cardMode='boost';
      confirmation('boost',selectedCardSummary(card,`BOOST +${card.boost??0}`),null,{
        confirmLabel:'BOOST',
        execute:()=>command({type:'choose',choice}),
      });
      renderHand();renderConfirmation();
    }
    function beginActionFromHand(card){
      if(!view.legal?.scheme_cards?.includes(card.id))return;
      state.selectedCardID=card.id;
      state.cardMode='action';
      const plan=actionPlan(view,card.id,selectedOwnFighterID());
      if(plan.kind==='confirm'){
        beginInteraction('action_confirm','scheme',{card_id:card.id});
        confirmation('action',selectedCardSummary(card,'Action'),null,{
          confirmLabel:'Play Action',
          execute:()=>command({type:'scheme',card_id:card.id}),
          cancel:resetInteraction,
        });
      }else{
        beginInteraction(plan.kind==='fighter'?'scheme_fighter':'scheme_destination','scheme',{
          card_id:card.id,
          fighter_id:plan.fighterID||'',
        });
        state.confirmation=null;
      }
      renderHand();renderLocalInteraction();renderConfirmation();
    }

    const baseApplyView=applyView;
    applyView=function(next,options={}){
      const changed=Boolean(!view||view.revision!==next?.revision||view.phase!==next?.phase||view.match_id!==next?.match_id);
      if(changed){
        clearConfirmation();
        state.cardMode=next?.pending?.kind==='defense'?'defense':'';
        state.boostOpen=false;
      }
      const applied=baseApplyView(next,options);
      if(applied&&changed){renderHand();renderConfirmation()}
      return applied;
    };

    const baseBoardHighlights=boardHighlights;
    boardHighlights=function(){
      const result=baseBoardHighlights();
      if(interaction.mode==='attack_card'){
        const fighterID=interactionValue('fighter_id');
        const targetID=interactionValue('target_id');
        result.selectedFighterID=fighterID;
        result.selectedTargetID=targetID;
        if(targetID)result.targets.add(targetID);
      }
      if(interaction.mode==='action_confirm')result.selectedFighterID=interactionValue('fighter_id');
      return result;
    };

    const baseInstruction=interactionInstruction;
    interactionInstruction=function(){
      if(interaction.mode==='attack_card')return 'Choose a highlighted attack card directly from your hand.';
      if(interaction.mode==='action_confirm')return 'Confirm the selected Action card below your hand.';
      if(view?.pending?.kind==='defense'&&view.pending.owner_id===view.viewing_player_id)return 'Choose a highlighted defense card directly from your hand.';
      if(view?.pending?.kind==='maneuver_boost'&&state.boostOpen)return 'Choose a highlighted card from your hand for BOOST.';
      return baseInstruction();
    };

    renderActions=function(){
      const box=$('actions');
      const mine=ownTurn(view);
      const selected=selectedOwnFighterID();
      if(!mine){box.innerHTML='<h2>Actions</h2><p>Waiting for the active player.</p>';return}
      if(!selected){box.innerHTML='<h2>Actions</h2><p>Select one of your fighters on the battlefield.</p>';return}
      if(!['idle','fighter_selected'].includes(interaction.mode)){
        box.innerHTML=`<h2>Current action</h2><p>${escapeHTML(interactionInstruction())}</p><button id="stage6-action-cancel" class="ghost">Cancel</button>`;
        $('stage6-action-cancel').onclick=cancelInteraction;
        return;
      }
      const fighter=fighterByID(selected);
      const attack=attackMapModel(view.legal,selected);
      const canAttack=attack.selectedFighterID===selected;
      box.innerHTML=`<h2>${escapeHTML(fighter?.name||'Selected fighter')}</h2><p class="selection-summary">Choose an action for this fighter.</p><div class="actions fighter-actions"><button id="maneuver" ${view.legal.can_maneuver?'':'disabled'}>Maneuver</button><button id="attack" ${canAttack?'':'disabled'}>Attack</button></div><p class="action-card-hint">Play Action cards directly from your hand.</p>`;
      $('maneuver').onclick=startManeuver;
      $('attack').onclick=startAttackInteraction;
    };

    openAttackCardPicker=function(attackerID,targetID){
      beginInteraction('attack_card','attack',{fighter_id:attackerID,target_id:targetID});
      state.cardMode='attack';state.selectedCardID='';clearConfirmation();
      renderLocalInteraction();renderHand();renderConfirmation();
      return true;
    };
    openSchemeCardPicker=function(){return false};
    startSchemeMap=function(cardID){const card=cardByID(cardID);if(card)beginActionFromHand(card);return true};

    submitSchemeDestination=async function(spaceID){
      const cardID=interactionValue('card_id');
      const fighterID=interactionValue('fighter_id');
      const card=cardByID(cardID);
      const action=view.legal.scheme_actions_by_card?.[cardID];
      const fighterAction=schemeFighterAction(action,fighterID);
      const destination=schemeDestinationOption(fighterAction,spaceID);
      if(!destination)return;
      rememberInteraction('destination',spaceID);
      rememberInteraction('path',[...(destination.path||[])]);
      const targets=schemeTargetIDs(fighterAction,spaceID);
      if(card?.effect==='horns'&&targets.length){
        interaction.mode='scheme_target';
        state.confirmation=null;
      }else{
        confirmation('action',`${card?.name||'Action'} → ${spaceID}`,null,{
          confirmLabel:'Play Action',
          execute:()=>command({type:'scheme',card_id:cardID,fighter_id:fighterID,path:[...(destination.path||[])]}),
          cancel:()=>{interaction.mode='scheme_destination'},
        });
      }
      renderLocalInteraction();renderHand();renderConfirmation();
    };
    submitSchemeTarget=async function(targetID){
      const cardID=interactionValue('card_id');
      const fighterID=interactionValue('fighter_id');
      const destinationID=interactionValue('destination');
      const action=view.legal.scheme_actions_by_card?.[cardID];
      const fighterAction=schemeFighterAction(action,fighterID);
      const destination=schemeDestinationOption(fighterAction,destinationID);
      if(!destination||!schemeTargetIDs(fighterAction,destinationID).includes(targetID))return;
      rememberInteraction('target_id',targetID);
      const card=cardByID(cardID),target=fighterByID(targetID);
      confirmation('action',`${card?.name||'Action'} → ${target?.name||'target'}`,null,{
        confirmLabel:'Play Action',
        execute:()=>command({type:'scheme',card_id:cardID,fighter_id:fighterID,target_id:targetID,path:[...(destination.path||[])]}),
        cancel:()=>{interaction.mode='scheme_target'},
      });
      renderLocalInteraction();renderHand();renderConfirmation();
    };

    renderHand=function(){
      const cards=ownHand(view);
      const legal=legalCardIDs(view,interaction,state);
      const activeMode=(view.pending?.kind==='defense'&&view.pending.owner_id===view.viewing_player_id)
        ?'defense'
        :(view.pending?.kind==='maneuver_boost'&&state.boostOpen)
          ?'boost'
          :interaction.mode==='attack_card'
            ?'attack'
            :['scheme_fighter','scheme_destination','scheme_target','action_confirm'].includes(interaction.mode)?'action':'';
      const html=cards.map(card=>{
        const role=cardRole(card,view,interaction,state);
        const selectable=legal.has(card.id);
        const classes=['card','hand-card','card-dock-item'];
        if(selectable)classes.push('legal-card',`legal-${role}`);
        if(activeMode&&!selectable)classes.push('card-disabled');
        if(state.selectedCardID===card.id)classes.push('selected');
        return `<button type="button" class="${classes.join(' ')}" data-hand-card="${escapeAttr(card.id)}" aria-pressed="${state.selectedCardID===card.id}">${stage4CardArtContent(card)}</button>`;
      }).join('');
      $('hand').innerHTML=html||'<p>No cards in hand.</p>';
      $('hand').classList.add('card-dock');
      renderConfirmation();
    };

    const baseRenderPending=renderPending;
    renderPending=function(){
      const pending=view.pending;
      if(!pending||pending.owner_id!==view.viewing_player_id||!['defense','maneuver_boost'].includes(pending.kind)){
        baseRenderPending();
        return;
      }
      const box=$('pending');
      box.hidden=false;
      const model=pendingCardModel(pending);
      if(pending.kind==='defense'){
        state.cardMode='defense';
        box.innerHTML=`<h2>Defense</h2><p>${escapeHTML(stage4UserFacingMessage(pending.message))}</p><p class="map-choice-note">Choose a blue highlighted card directly from your hand.</p><div class="pending-actions">${model.skipOptions.map(option=>`<button data-choice="${escapeAttr(option.id)}" class="ghost">${escapeHTML(option.label)}</button>`).join('')}</div>`;
      }else{
        box.innerHTML=`<h2>Maneuver</h2><p>${escapeHTML(stage4UserFacingMessage(pending.message))}</p><div class="pending-actions">${model.skipOptions.map(option=>`<button data-choice="${escapeAttr(option.id)}">${escapeHTML(option.label)}</button>`).join('')}<button id="stage6-open-boost" ${model.cardIDs.length?'':'disabled'}>BOOST</button></div>${state.boostOpen?'<p class="map-choice-note">Choose a yellow highlighted card directly from your hand.</p>':''}`;
        const boost=$('stage6-open-boost');
        if(boost)boost.onclick=()=>{state.boostOpen=true;state.cardMode='boost';clearConfirmation();renderPending();renderHand();renderInteractionPanel()};
      }
    };

    const hand=$('hand');
    if(hand)hand.onclick=event=>{
      const button=event.target?.closest?.('[data-hand-card]');
      if(!button)return;
      const card=cardByID(button.dataset.handCard);
      if(!card)return;
      if(view.pending?.owner_id===view.viewing_player_id&&view.pending.kind==='defense')return selectDefenseCard(card);
      if(view.pending?.owner_id===view.viewing_player_id&&view.pending.kind==='maneuver_boost'&&state.boostOpen)return selectBoostCard(card);
      if(interaction.mode==='attack_card')return selectAttackCard(card);
      if(ownTurn(view)&&['idle','fighter_selected'].includes(interaction.mode)&&view.legal?.scheme_cards?.includes(card.id))return beginActionFromHand(card);
    };

    const pendingBox=$('pending');
    if(pendingBox)pendingBox.onclick=event=>{
      const button=event.target?.closest?.('[data-choice]');
      if(button)void command({type:'choose',choice:button.dataset.choice});
    };

    if(view)render();
  }

  global.Stage6UX={
    ownHand,ownTurn,optionCardMap,pendingCardModel,legalCardIDs,cardRole,actionPlan,state,
  };
  install();
})(globalThis);
