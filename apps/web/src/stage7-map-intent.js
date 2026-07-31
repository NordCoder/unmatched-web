(function(global){
  let maneuverContinuation=null;

  function ownTurn(currentView=typeof view==='undefined'?null:view){
    return Boolean(
      currentView&&
      currentView.phase==='active'&&
      !currentView.pending&&
      currentView.current_player_id===currentView.viewing_player_id,
    );
  }

  function viewFighters(currentView){
    if(currentView?.fighters?.length)return currentView.fighters;
    return currentView?.spaces?.flatMap(space=>space.fighter?[space.fighter]:[])||[];
  }

  function viewFighter(currentView,fighterID){
    return viewFighters(currentView).find(fighter=>fighter.id===fighterID);
  }

  function maneuverMapModel(legal,selectedFighterID=''){
    const byFighter=legal?.maneuver_action?.destinations_by_fighter||{};
    const fighters=Object.keys(byFighter).filter(fighterID=>(byFighter[fighterID]||[]).length>0);
    const fighterID=fighters.includes(selectedFighterID)?selectedFighterID:'';
    return {
      baseMovement:Number(legal?.maneuver_action?.base_movement)||0,
      fighters,
      selectedFighterID:fighterID,
      destinations:fighterID?(byFighter[fighterID]||[]):[],
    };
  }

  function attackIntentModel(legal,selectedFighterID=''){
    const cards=legal?.attack_cards_by_fighter||{};
    const targets=legal?.attack_targets_by_fighter||{};
    const fighters=Object.keys(cards).filter(fighterID=>(cards[fighterID]||[]).length&&(targets[fighterID]||[]).length);
    const fighterID=fighters.includes(selectedFighterID)?selectedFighterID:'';
    return {
      fighters,
      selectedFighterID:fighterID,
      targetIDs:fighterID?(targets[fighterID]||[]):[],
    };
  }

  function fighterIntentModel(currentView,selectedFighterID=''){
    const fighter=viewFighter(currentView,selectedFighterID);
    const selectable=Boolean(
      ownTurn(currentView)&&
      fighter&&
      fighter.owner_id===currentView.viewing_player_id&&
      !fighter.defeated,
    );
    if(!selectable){
      return {fighterID:'',baseMovement:0,destinations:[],targetIDs:[]};
    }
    const maneuver=maneuverMapModel(currentView.legal,selectedFighterID);
    const attack=attackIntentModel(currentView.legal,selectedFighterID);
    return {
      fighterID:selectedFighterID,
      baseMovement:maneuver.baseMovement,
      destinations:maneuver.destinations,
      targetIDs:attack.targetIDs,
    };
  }

  function activationFor(currentView,currentInteraction,spaceID,fighterID=''){
    if(!ownTurn(currentView)||!['idle','fighter_selected'].includes(currentInteraction?.mode))return {kind:'none'};
    const fighter=viewFighter(currentView,fighterID);
    const selected=currentInteraction?.values?.fighter_id||'';
    if(fighter&&fighter.owner_id===currentView.viewing_player_id&&!fighter.defeated){
      return fighter.id===selected
        ?{kind:'clear-fighter',fighterID:fighter.id}
        :{kind:'select-fighter',fighterID:fighter.id};
    }
    if(currentInteraction.mode!=='fighter_selected'||!selected)return {kind:'none'};
    const intent=fighterIntentModel(currentView,selected);
    if(fighterID&&intent.targetIDs.includes(fighterID)){
      return {kind:'attack',fighterID:selected,targetID:fighterID};
    }
    const destination=intent.destinations.find(option=>option.destination===spaceID);
    if(destination){
      return {
        kind:'maneuver',
        fighterID:selected,
        destination:destination.destination,
        path:[...(destination.path||[])],
      };
    }
    return {kind:'none'};
  }

  function createManeuverContinuation(currentView,fighterID,destination,path=[]){
    const intent=fighterIntentModel(currentView,fighterID);
    const option=intent.destinations.find(candidate=>candidate.destination===destination);
    if(!option)return null;
    return {
      kind:'maneuver',
      match_id:currentView.match_id,
      player_id:currentView.viewing_player_id,
      source_revision:currentView.revision,
      fighter_id:fighterID,
      intended_destination:destination,
      intended_path:[...(option.path||path||[])],
      selected_destination:destination,
      selected_path:[...(option.path||path||[])],
      choice_id:'',
    };
  }

  function maneuverContinuationOption(pending,continuation,destination=continuation?.selected_destination||continuation?.intended_destination){
    if(!pending||pending.kind!=='maneuver_move'||!continuation)return null;
    return (pending.options||[]).find(option=>
      option.fighter_id===continuation.fighter_id&&
      option.destination===destination,
    )||null;
  }

  function clearContinuationResult(reason){
    return {continuation:null,state:'cleared',option:null,reason};
  }

  function reconcileManeuverContinuation(continuation,currentView){
    if(!continuation)return {continuation:null,state:'none',option:null,reason:''};
    if(!currentView)return clearContinuationResult('missing_view');
    if(currentView.match_id!==continuation.match_id)return clearContinuationResult('match_changed');
    if(currentView.viewing_player_id!==continuation.player_id)return clearContinuationResult('viewer_changed');
    if(currentView.phase!=='active')return clearContinuationResult('phase_changed');
    if(currentView.current_player_id!==continuation.player_id)return clearContinuationResult('turn_changed');
    const fighter=viewFighter(currentView,continuation.fighter_id);
    if(!fighter||fighter.owner_id!==continuation.player_id||fighter.defeated){
      return clearContinuationResult('fighter_unavailable');
    }

    const pending=currentView.pending;
    if(!pending){
      if(currentView.revision===continuation.source_revision){
        return {continuation,state:'starting',option:null,reason:''};
      }
      return clearContinuationResult('expected_maneuver_pending');
    }
    if(pending.owner_id!==continuation.player_id)return clearContinuationResult('pending_owner_changed');
    if(pending.kind==='maneuver_boost'){
      return {continuation,state:'boost',option:null,reason:''};
    }
    if(pending.kind!=='maneuver_move')return clearContinuationResult('unexpected_pending_kind');

    const option=maneuverContinuationOption(pending,continuation);
    if(!option)return clearContinuationResult('destination_unavailable');
    return {
      continuation:{
        ...continuation,
        selected_destination:option.destination,
        selected_path:[...(option.path||[])],
        choice_id:option.id,
      },
      state:'destination',
      option,
      reason:'',
    };
  }

  function selectManeuverContinuationDestination(continuation,pending,destination){
    const option=maneuverContinuationOption(pending,continuation,destination);
    if(!option)return null;
    return {
      continuation:{
        ...continuation,
        selected_destination:option.destination,
        selected_path:[...(option.path||[])],
        choice_id:option.id,
      },
      option,
    };
  }

  function activeContinuation(currentView=typeof view==='undefined'?null:view){
    const hadContinuation=Boolean(maneuverContinuation);
    const result=reconcileManeuverContinuation(maneuverContinuation,currentView);
    maneuverContinuation=result.continuation;
    if(hadContinuation&&!maneuverContinuation&&typeof preferredManeuverFighterID!=='undefined'){
      preferredManeuverFighterID=null;
    }
    return result;
  }

  function clearManeuverContinuation(){
    maneuverContinuation=null;
    if(typeof preferredManeuverFighterID!=='undefined')preferredManeuverFighterID=null;
  }

  function install(){
    if(typeof document==='undefined'||typeof $!=='function'||typeof boardHighlights!=='function')return;

    const baseCommand=command;
    command=async function(payload){
      const continuationCommand=Boolean(
        maneuverContinuation&&
        payload?.type==='choose'&&
        ['maneuver_boost','maneuver_move'].includes(view?.pending?.kind),
      );
      if(continuationCommand&&view.pending.kind==='maneuver_move')clearManeuverContinuation();
      const accepted=await baseCommand(payload);
      if(accepted===false&&(
        payload?.type==='maneuver'||
        continuationCommand
      ))clearManeuverContinuation();
      return accepted;
    };

    const baseLeave=leave;
    leave=function(){
      clearManeuverContinuation();
      baseLeave();
    };

    const baseSyncPendingMapInteraction=syncPendingMapInteraction;
    syncPendingMapInteraction=function(){
      const current=activeContinuation(view);
      baseSyncPendingMapInteraction();
      if(current.state!=='destination'||!current.option)return;
      beginInteraction('pending_destination','pending',{
        pending_identity:pendingIdentity(view.pending),
        fighter_id:current.continuation.fighter_id,
        destination:current.option.destination,
        path:[...(current.option.path||[])],
        continuation_choice_id:current.option.id,
      });
    };

    const baseBoardHighlights=boardHighlights;
    boardHighlights=function(){
      const result=baseBoardHighlights();
      if(!ownTurn(view)&&['idle','fighter_selected'].includes(interaction.mode)){
        result.fighterCandidates.clear();
        result.destinations.clear();
        result.targets.clear();
        result.selectedFighterID='';
        return result;
      }
      if(interaction.mode==='fighter_selected'){
        const intent=fighterIntentModel(view,interactionValue('fighter_id'));
        for(const option of intent.destinations)result.destinations.set(option.destination,option);
        for(const targetID of intent.targetIDs)result.targets.add(targetID);
        result.selectedFighterID=intent.fighterID;
      }
      const current=activeContinuation(view);
      if(current.state==='destination'&&current.option){
        result.selectedFighterID=current.continuation.fighter_id;
        result.selectedDestinationID=current.option.destination;
      }
      return result;
    };

    const basePathOptionForSpace=pathOptionForSpace;
    pathOptionForSpace=function(spaceID){
      if(interaction.mode==='fighter_selected'){
        return fighterIntentModel(view,interactionValue('fighter_id')).destinations.find(option=>option.destination===spaceID);
      }
      return basePathOptionForSpace(spaceID);
    };

    const baseInstruction=interactionInstruction;
    interactionInstruction=function(){
      const current=activeContinuation(view);
      if(current.state==='boost')return 'Choose whether to BOOST. The originally selected fighter and destination will be restored afterwards.';
      if(current.state==='destination'&&current.option){
        const fighter=fighterByID(current.continuation.fighter_id);
        return `${fighter?.name||'Fighter'} is ready to move to ${current.option.destination}. Confirm below or choose another highlighted destination.`;
      }
      if(interaction.mode==='fighter_selected'){
        const fighter=fighterByID(interactionValue('fighter_id'));
        const intent=fighterIntentModel(view,interactionValue('fighter_id'));
        const choices=[];
        if(intent.destinations.length)choices.push('a highlighted empty space to start Maneuver');
        if(intent.targetIDs.length)choices.push('a red opponent to attack');
        if(!choices.length)return `${fighter?.name||'Fighter'} selected. No direct move or attack is currently available.`;
        return `${fighter?.name||'Fighter'} selected. Choose ${choices.join(' or ')}.`;
      }
      return baseInstruction();
    };

    const baseRenderInteractionPanel=renderInteractionPanel;
    renderInteractionPanel=function(){
      baseRenderInteractionPanel();
      const current=activeContinuation(view);
      if(current.state!=='destination')return;
      const cancel=$('interaction-cancel');
      if(cancel)cancel.hidden=true;
    };

    const baseRenderBoard=renderBoard;
    renderBoard=function(){
      baseRenderBoard();
      const board=$('board');
      if(!board)return;
      const current=activeContinuation(view);
      board.dataset.interactionMode=interaction.mode==='fighter_selected'
        ?'fighter-intent'
        :current.state==='destination'?'maneuver-confirmation':'';
      for(const node of board.querySelectorAll?.('[data-space]')||[]){
        const label=node.getAttribute?.('aria-label')||'';
        if(node.classList?.contains('legal-target'))node.setAttribute?.('aria-label',`${label}. Legal attack target`);
        else if(node.classList?.contains('legal-destination'))node.setAttribute?.('aria-label',`${label}. Legal Maneuver destination`);
      }
      if(current.state!=='destination'||!current.option)return;
      const node=board.querySelector?.(`[data-space="${current.option.destination}"]`);
      if(!node)return;
      node.classList?.add('preselected-destination');
      const label=node.getAttribute?.('aria-label')||`Space ${current.option.destination}`;
      node.setAttribute?.('aria-label',`${label}. Preselected Maneuver destination. Confirm below or choose another destination.`);
      const point=global.BattlefieldRenderer?.point?.(
        global.BattlefieldRenderer.presentationFor(view),
        current.option.destination,
      );
      if(point&&typeof document.createElementNS==='function'){
        const marker=document.createElementNS('http://www.w3.org/2000/svg','text');
        marker.setAttribute('class','destination-check');
        marker.setAttribute('x',String(point.x));
        marker.setAttribute('y',String(point.y));
        marker.setAttribute('aria-hidden','true');
        marker.textContent='✓';
        node.append(marker);
      }
      setPathPreview(current.option.destination);
    };

    async function confirmManeuverDestination(){
      const current=activeContinuation(view);
      if(current.state!=='destination'||!current.option)return false;
      const choice=current.option.id;
      clearManeuverContinuation();
      return command({type:'choose',choice});
    }

    renderActions=function(){
      const box=$('actions');
      if(!box)return;
      const current=activeContinuation(view);
      if(current.state==='destination'&&current.option){
        const fighter=fighterByID(current.continuation.fighter_id);
        const name=fighter?.name||'Fighter';
        box.innerHTML=`<h2>Maneuver</h2><p class="selection-summary">${escapeHTML(name)} → ${escapeHTML(current.option.destination)}</p><p>Confirm the preselected destination or choose another highlighted space. The final path is the server-projected path after BOOST.</p><div class="actions"><button id="stage7-confirm-maneuver">Move ${escapeHTML(name)} to ${escapeHTML(current.option.destination)}</button></div><p class="action-card-hint">Finish maneuver remains available in the pending-choice panel.</p>`;
        $('stage7-confirm-maneuver').onclick=()=>void confirmManeuverDestination();
        return;
      }
      if(!ownTurn(view)){
        box.innerHTML='<h2>Actions</h2><p>Waiting for the active player.</p>';
        return;
      }
      if(!['idle','fighter_selected'].includes(interaction.mode)){
        box.innerHTML=`<h2>Current action</h2><p>${escapeHTML(interactionInstruction())}</p><button id="stage7-action-cancel" class="ghost">Cancel</button>`;
        $('stage7-action-cancel').onclick=cancelInteraction;
        return;
      }
      const selected=selectedOwnFighterID();
      if(!selected){
        box.innerHTML='<h2>Actions</h2><p>Select one of your fighters on the battlefield. Action cards remain playable directly from your hand.</p>';
        return;
      }
      const fighter=fighterByID(selected);
      const intent=fighterIntentModel(view,selected);
      const moveText=intent.destinations.length
        ?`${intent.destinations.length} base-movement destination${intent.destinations.length===1?'':'s'}`
        :'No base-movement destination';
      const attackText=intent.targetIDs.length
        ?`${intent.targetIDs.length} attack target${intent.targetIDs.length===1?'':'s'}`
        :'No attack target';
      box.innerHTML=`<h2>${escapeHTML(fighter?.name||'Selected fighter')}</h2><p class="selection-summary">Use the battlefield directly.</p><p><strong>Move:</strong> ${escapeHTML(moveText)} · <strong>Attack:</strong> ${escapeHTML(attackText)}</p><p class="action-card-hint">Highlighted empty spaces start Maneuver. Red opponents start Attack. Play Action cards from your hand.</p><div class="actions"><button id="stage7-clear-selection" class="ghost">Clear selection</button></div><details><summary>Additional action</summary><button id="stage7-start-maneuver" class="ghost" ${view.legal?.can_maneuver?'':'disabled'}>Start Maneuver without a destination</button></details>`;
      $('stage7-clear-selection').onclick=()=>{resetInteraction();renderLocalInteraction()};
      $('stage7-start-maneuver').onclick=()=>{clearManeuverContinuation();void startManeuver()};
    };

    const baseHandleBoardSelection=handleBoardSelection;
    handleBoardSelection=function(spaceID,fighterID=''){
      const current=activeContinuation(view);
      if(current.state==='destination'&&interaction.mode==='pending_destination'){
        const selection=selectManeuverContinuationDestination(maneuverContinuation,view.pending,spaceID);
        if(!selection)return;
        maneuverContinuation=selection.continuation;
        rememberInteraction('destination',selection.option.destination);
        rememberInteraction('path',[...(selection.option.path||[])]);
        rememberInteraction('continuation_choice_id',selection.option.id);
        renderLocalInteraction();
        return;
      }
      if(!['idle','fighter_selected'].includes(interaction.mode)){
        return baseHandleBoardSelection(spaceID,fighterID);
      }
      const activation=activationFor(view,interaction,spaceID,fighterID);
      switch(activation.kind){
      case 'select-fighter':
        clearManeuverContinuation();
        beginInteraction('fighter_selected','select',{fighter_id:activation.fighterID});
        renderLocalInteraction();
        return;
      case 'clear-fighter':
        clearManeuverContinuation();
        resetInteraction();
        renderLocalInteraction();
        return;
      case 'attack':
        clearManeuverContinuation();
        rememberInteraction('target_id',activation.targetID);
        openAttackCardPicker(activation.fighterID,activation.targetID);
        return;
      case 'maneuver':
        maneuverContinuation=createManeuverContinuation(
          view,
          activation.fighterID,
          activation.destination,
          activation.path,
        );
        if(!maneuverContinuation)return;
        preferredManeuverFighterID=activation.fighterID;
        void startManeuver();
        return;
      default:
        return;
      }
    };

    document.addEventListener('keydown',event=>{
      if(event.key!=='Escape'||interaction.mode!=='fighter_selected')return;
      event.preventDefault();
      clearManeuverContinuation();
      resetInteraction();
      renderLocalInteraction();
    });

    if(view)renderLocalInteraction();
  }

  global.Stage7MapIntent={
    ownTurn,
    maneuverMapModel,
    attackIntentModel,
    fighterIntentModel,
    activationFor,
    createManeuverContinuation,
    maneuverContinuationOption,
    reconcileManeuverContinuation,
    selectManeuverContinuationDestination,
  };
  install();
})(globalThis);
