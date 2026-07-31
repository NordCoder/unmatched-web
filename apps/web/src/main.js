const $ = (id) => document.getElementById(id);
const lobby = $('lobby'), game = $('game'), errorBox = $('error');
const MAP_WIDTH = 100;
const MAP_HEIGHT = 742 / 1337 * 100;
const MAP_ASSET = '/sherwood-forest.svg';

let session = JSON.parse(sessionStorage.getItem('unmatched-session') || 'null');
let view = null;
let pollHandle = null;
let commandInFlight = false;
let refreshInFlight = false;
let renderedPendingIdentity = null;
let preferredManeuverFighterID = null;
const interaction = {mode:'idle',action:null,revision:null,values:{}};

function saveSession(next){
  session=next;
  if(next)sessionStorage.setItem('unmatched-session',JSON.stringify(next));
  else sessionStorage.removeItem('unmatched-session');
}
function showError(message){errorBox.textContent=message;errorBox.hidden=false}
function clearError(){errorBox.hidden=true;errorBox.textContent=''}
async function request(path,options={}){
  clearError();
  const headers={'Content-Type':'application/json',...(options.headers||{})};
  if(session?.token)headers.Authorization=`Bearer ${session.token}`;
  const response=await fetch(path,{...options,headers});
  const body=await response.json().catch(()=>({error:{message:`HTTP ${response.status}`}}));
  if(!response.ok)throw new Error(body.error?.message||`HTTP ${response.status}`);
  return body;
}

function shouldApplyView(current,next,force=false){
  return Boolean(force||!current||!next||current.match_id!==next.match_id||current.revision!==next.revision||current.phase!==next.phase||current.viewing_player_id!==next.viewing_player_id);
}
function interactionCanContinue(state,next){
  if(!state||state.mode==='idle'||!next||state.revision!==next.revision)return false;
  if(state.mode==='dialog')return next.phase==='active'&&!next.pending&&next.current_player_id===next.viewing_player_id;
  return true;
}
function resetInteraction(){
  interaction.mode='idle';
  interaction.action=null;
  interaction.revision=null;
  interaction.values={};
}
function beginInteraction(mode,action,values={}){
  interaction.mode=mode;
  interaction.action=action;
  interaction.revision=view?.revision??null;
  interaction.values={...values};
}
function rememberInteraction(name,value){interaction.values[name]=value}
function interactionValue(name,fallback=''){return interaction.values[name]??fallback}
function closeActionDialog(){
  const dialog=$('action-dialog');
  if(dialog?.open&&typeof dialog.close==='function')dialog.close();
}
function reconcileInteraction(next){
  if(!view||view.revision===next?.revision)return;
  closeActionDialog();
  resetInteraction();
}
function applyView(next,options={}){
  if(!next||!shouldApplyView(view,next,Boolean(options.force)))return false;
  reconcileInteraction(next);
  view=next;
  render();
  return true;
}

async function createMatch(){
  try{
    const body=await request('/api/matches',{method:'POST',body:'{}'});
    saveSession({matchID:body.match_id,token:body.token,playerID:body.player_id});
    applyView(body.view,{force:true});
    enterGame();
  }catch(error){showError(error.message)}
}
async function joinMatch(){
  const code=$('join-code').value.trim();
  if(!code)return showError('Enter a match code.');
  try{
    const body=await request(`/api/matches/${encodeURIComponent(code)}/join`,{method:'POST',body:'{}'});
    saveSession({matchID:body.match_id,token:body.token,playerID:body.player_id});
    applyView(body.view,{force:true});
    enterGame();
  }catch(error){showError(error.message)}
}
async function refresh(){
  if(!session||commandInFlight||refreshInFlight)return false;
  refreshInFlight=true;
  try{
    const body=await request(`/api/matches/${encodeURIComponent(session.matchID)}`);
    return applyView(body.view);
  }catch(error){
    showError(error.message);
    return false;
  }finally{refreshInFlight=false}
}
function beginCommand(){if(commandInFlight)return false;commandInFlight=true;return true}
function endCommand(){commandInFlight=false}
async function command(payload){
  if(!beginCommand())return false;
  let failed=false;
  try{
    const body=await request(`/api/matches/${encodeURIComponent(session.matchID)}/commands`,{
      method:'POST',
      body:JSON.stringify({...payload,expected_revision:view.revision}),
    });
    applyView(body.view);
    return true;
  }catch(error){
    showError(error.message);
    failed=true;
    return false;
  }finally{
    endCommand();
    if(failed)void refresh();
  }
}
function enterGame(){
  lobby.hidden=true;
  game.hidden=false;
  $('reset').hidden=false;
  if(pollHandle)clearInterval(pollHandle);
  pollHandle=setInterval(refresh,1200);
}
function leave(){
  if(pollHandle)clearInterval(pollHandle);
  pollHandle=null;
  saveSession(null);
  view=null;
  refreshInFlight=false;
  renderedPendingIdentity=null;
  preferredManeuverFighterID=null;
  resetInteraction();
  game.hidden=true;
  lobby.hidden=false;
  $('reset').hidden=true;
  clearError();
}

function render(){
  if(!view)return;
  syncPendingMapInteraction();
  renderStatus();
  renderInteractionPanel();
  renderBoard();
  renderPlayers();
  renderPending();
  renderActions();
  renderCombat();
  renderHand();
  renderEvents();
}
function renderLocalInteraction(){
  if(!view)return;
  renderInteractionPanel();
  renderBoard();
  renderActions();
}
function renderStatus(){
  const current=view.players.find(player=>player.id===view.current_player_id);
  const mine=view.viewing_player_id===view.current_player_id;
  const status=view.phase==='ended'
    ?`<span class="winner">${escapeHTML(view.winner_name)} wins</span>`
    :view.phase==='waiting_for_player'
      ?`Waiting for Bigfoot. Match code: <strong>${escapeHTML(view.match_id)}</strong>`
      :`${mine?'<span class="current">Your turn</span>':`${escapeHTML(current?.name||'Opponent')}'s turn`} · ${current?.actions_remaining??0} actions remaining`;
  $('status').innerHTML=`<div>${status}</div><small>Match <strong>${escapeHTML(view.match_id)}</strong> · revision ${view.revision}</small>`;
}

function mapDimensions(){return {width:MAP_WIDTH,height:MAP_HEIGHT}}
function mapY(percent){return percent/100*MAP_HEIGHT}
function spaceByID(id){return view?.spaces?.find(space=>space.id===id)}
function allFighters(){return view?.fighters||view?.spaces?.flatMap(space=>space.fighter?[space.fighter]:[])||[]}
function fighterByID(id){return allFighters().find(fighter=>fighter.id===id)}
function cardByID(id){return view?.players?.find(player=>player.id===view.viewing_player_id)?.hand?.find(card=>card.id===id)}
function ownLivingFighterIDs(){
  return allFighters().filter(fighter=>fighter.owner_id===view.viewing_player_id&&!fighter.defeated).map(fighter=>fighter.id);
}

function pendingIdentity(pending){
  if(!pending)return 'none';
  return JSON.stringify([
    pending.kind,
    pending.owner_id,
    pending.message,
    (pending.options||[]).map(option=>[option.id,option.label,option.fighter_id,option.destination,option.path]),
  ]);
}
function isMapPendingKind(kind){
  return ['maneuver_move','move','place','skirmish','return_sidekick'].includes(kind);
}
function pendingMapModel(pending,selectedFighterID=''){
  const options=(pending?.options||[]).filter(option=>option.fighter_id&&option.destination);
  const specialOptions=(pending?.options||[]).filter(option=>!option.fighter_id||!option.destination);
  const fighters=[...new Set(options.map(option=>option.fighter_id))];
  const preferred=fighters.includes(selectedFighterID)?selectedFighterID:'';
  const destinationOptions=preferred?options.filter(option=>option.fighter_id===preferred):[];
  return {fighters,selectedFighterID:preferred,destinationOptions,specialOptions,options};
}
function syncPendingMapInteraction(){
  const pending=view.pending;
  if(!pending||pending.owner_id!==view.viewing_player_id||!isMapPendingKind(pending.kind)){
    if(interaction.action==='pending')resetInteraction();
    return;
  }
  const identity=pendingIdentity(pending);
  if(interaction.action==='pending'&&interaction.values.pending_identity===identity)return;
  const preferred=pending.kind==='maneuver_move'?preferredManeuverFighterID:'';
  const model=pendingMapModel(pending,preferred);
  const fighterID=model.selectedFighterID||(model.fighters.length===1?model.fighters[0]:'');
  beginInteraction(fighterID?'pending_destination':'pending_fighter','pending',{
    pending_identity:identity,
    fighter_id:fighterID,
  });
}

function attackMapModel(legal,selectedFighterID=''){
  const cards=legal?.attack_cards_by_fighter||{};
  const targets=legal?.attack_targets_by_fighter||{};
  const fighters=Object.keys(cards).filter(fighterID=>(cards[fighterID]||[]).length&&(targets[fighterID]||[]).length);
  const fighterID=fighters.includes(selectedFighterID)?selectedFighterID:'';
  return {fighters,selectedFighterID:fighterID,targetIDs:fighterID?(targets[fighterID]||[]):[]};
}
function schemeFighterAction(action,fighterID){
  return (action?.fighters||[]).find(candidate=>candidate.fighter_id===fighterID);
}
function schemeDestinationOption(fighterAction,destinationID){
  return (fighterAction?.destinations||[]).find(option=>option.destination===destinationID);
}
function schemeTargetIDs(fighterAction,destinationID){
  return fighterAction?.targets_by_destination?.[destinationID]||[];
}
function schemeMapModel(action,selectedFighterID='',selectedDestinationID=''){
  const fighters=(action?.fighters||[]).map(candidate=>candidate.fighter_id);
  const fighterID=fighters.includes(selectedFighterID)?selectedFighterID:'';
  const fighterAction=schemeFighterAction(action,fighterID);
  const destinations=fighterAction?.destinations||[];
  const destination=schemeDestinationOption(fighterAction,selectedDestinationID);
  return {
    fighters,
    selectedFighterID:fighterID,
    destinations,
    selectedDestination:destination,
    targetIDs:destination?schemeTargetIDs(fighterAction,destination.destination):[],
  };
}

function boardHighlights(){
  const result={
    fighterCandidates:new Set(),
    targets:new Set(),
    destinations:new Map(),
    selectedFighterID:interactionValue('fighter_id'),
    selectedTargetID:interactionValue('target_id'),
  };
  switch(interaction.mode){
  case 'idle':
  case 'fighter_selected':
    ownLivingFighterIDs().forEach(id=>result.fighterCandidates.add(id));
    break;
  case 'attack_fighter':{
    const model=attackMapModel(view.legal);
    model.fighters.forEach(id=>result.fighterCandidates.add(id));
    break;
  }
  case 'attack_target':{
    const model=attackMapModel(view.legal,interactionValue('fighter_id'));
    model.targetIDs.forEach(id=>result.targets.add(id));
    break;
  }
  case 'scheme_fighter':{
    const action=view.legal.scheme_actions_by_card?.[interactionValue('card_id')];
    schemeMapModel(action).fighters.forEach(id=>result.fighterCandidates.add(id));
    break;
  }
  case 'scheme_destination':{
    const action=view.legal.scheme_actions_by_card?.[interactionValue('card_id')];
    const model=schemeMapModel(action,interactionValue('fighter_id'));
    model.destinations.forEach(option=>result.destinations.set(option.destination,option));
    break;
  }
  case 'scheme_target':{
    const action=view.legal.scheme_actions_by_card?.[interactionValue('card_id')];
    const model=schemeMapModel(action,interactionValue('fighter_id'),interactionValue('destination'));
    model.targetIDs.forEach(id=>result.targets.add(id));
    break;
  }
  case 'pending_fighter':{
    const model=pendingMapModel(view.pending);
    model.fighters.forEach(id=>result.fighterCandidates.add(id));
    break;
  }
  case 'pending_destination':{
    const model=pendingMapModel(view.pending,interactionValue('fighter_id'));
    model.destinationOptions.forEach(option=>result.destinations.set(option.destination,option));
    model.fighters.forEach(id=>result.fighterCandidates.add(id));
    break;
  }
  }
  return result;
}
function boardPoint(spaceID){
  const space=spaceByID(spaceID);
  return space?`${space.x},${mapY(space.y)}`:'';
}
function renderBoard(){
  const svg=$('board');
  const highlights=boardHighlights();
  const nodes=view.spaces.map(space=>{
    const fighter=space.fighter;
    const fighterID=fighter?.id||'';
    const classes=['board-space'];
    if(highlights.destinations.has(space.id))classes.push('legal-destination');
    if(fighterID&&highlights.fighterCandidates.has(fighterID))classes.push('legal-fighter');
    if(fighterID&&highlights.targets.has(fighterID))classes.push('legal-target');
    if(fighterID&&fighterID===highlights.selectedFighterID)classes.push('selected-fighter');
    if(fighterID&&fighterID===highlights.selectedTargetID)classes.push('selected-target');
    const interactive=classes.length>1||(fighter&&fighter.owner_id===view.viewing_player_id&&!fighter.defeated);
    const x=space.x,y=mapY(space.y);
    const fighterMarkup=fighter?`
      <circle class="fighter-token ${fighter.owner_id===view.viewing_player_id?'friendly-token':'opposing-token'}" cx="${x}" cy="${y}" r="2.15"></circle>
      <text class="fighter-label" x="${x}" y="${y+.68}">${initials(fighter.name)}</text>
      <text class="health-label" x="${x}" y="${y+4.15}">${fighter.health}</text>`:'';
    const aria=fighter?`${fighter.name}, ${fighter.health} health, space ${space.id}`:`Space ${space.id}`;
    return `<g class="${classes.join(' ')}" data-space="${space.id}"${fighterID?` data-fighter="${fighterID}"`:''}${interactive?' role="button" tabindex="0"':''} aria-label="${escapeAttr(aria)}">
      <circle class="space-hitbox" cx="${x}" cy="${y}" r="4.25"></circle>
      <circle class="space-highlight" cx="${x}" cy="${y}" r="4.05"></circle>
      ${fighterMarkup}
      <title>${escapeHTML(aria)} · ${escapeHTML(space.zones.join(', '))}</title>
    </g>`;
  }).join('');
  svg.setAttribute('viewBox',`0 0 ${MAP_WIDTH} ${MAP_HEIGHT}`);
  svg.innerHTML=`
    <image class="board-art" href="${MAP_ASSET}" x="0" y="0" width="${MAP_WIDTH}" height="${MAP_HEIGHT}" preserveAspectRatio="none"></image>
    <polyline id="path-preview" class="path-preview" points=""></polyline>
    <g class="board-overlay">${nodes}</g>`;
}
function pathOptionForSpace(spaceID){
  if(interaction.mode==='pending_destination'){
    return pendingMapModel(view.pending,interactionValue('fighter_id')).destinationOptions.find(option=>option.destination===spaceID);
  }
  if(interaction.mode==='scheme_destination'){
    const action=view.legal.scheme_actions_by_card?.[interactionValue('card_id')];
    return schemeDestinationOption(schemeFighterAction(action,interactionValue('fighter_id')),spaceID);
  }
  return undefined;
}
function setPathPreview(spaceID=''){
  const line=$('path-preview');
  if(!line)return;
  const option=spaceID?pathOptionForSpace(spaceID):undefined;
  const fighter=fighterByID(interactionValue('fighter_id'));
  if(!option||!fighter){line.setAttribute('points','');return}
  const points=[fighter.space_id,...(option.path||[])].map(boardPoint).filter(Boolean);
  line.setAttribute('points',points.join(' '));
}

function interactionInstruction(){
  const fighter=fighterByID(interactionValue('fighter_id'));
  const card=cardByID(interactionValue('card_id'));
  switch(interaction.mode){
  case 'fighter_selected': return `${fighter?.name||'Fighter'} selected. Choose Maneuver, Scheme, or Attack.`;
  case 'attack_fighter': return 'Choose a highlighted friendly fighter to attack.';
  case 'attack_target': return `Choose a highlighted target for ${fighter?.name||'the attacker'}.`;
  case 'scheme_fighter': return `Choose a highlighted fighter to use ${card?.name||'the Scheme'}.`;
  case 'scheme_destination': return `Choose a highlighted destination for ${fighter?.name||'the fighter'}. Hover a space to preview the path.`;
  case 'scheme_target': return 'Choose a highlighted adjacent fighter for the damage step.';
  case 'pending_fighter': return 'Choose which highlighted fighter to move.';
  case 'pending_destination': return `Choose a highlighted destination for ${fighter?.name||'the fighter'}. Hover a space to preview the path.`;
  default: return 'Select one of your fighters on the map, then choose an action.';
  }
}
function renderInteractionPanel(){
  const panel=$('board-interaction');
  if(!panel)return;
  const pendingMode=interaction.action==='pending';
  const cancellable=interaction.mode!=='idle'&&(!pendingMode||pendingMapModel(view.pending).fighters.length>1);
  panel.innerHTML=`<div><strong>Map control</strong><span>${escapeHTML(interactionInstruction())}</span></div><button id="interaction-cancel" class="ghost" ${cancellable?'':'hidden'}>${pendingMode?'Change fighter':'Cancel'}</button>`;
  const cancel=$('interaction-cancel');
  if(cancel)cancel.onclick=cancelInteraction;
}
function cancelInteraction(){
  if(interaction.action==='pending'){
    const model=pendingMapModel(view.pending);
    beginInteraction(model.fighters.length===1?'pending_destination':'pending_fighter','pending',{
      pending_identity:pendingIdentity(view.pending),
      fighter_id:model.fighters.length===1?model.fighters[0]:'',
    });
  }else{
    resetInteraction();
  }
  renderLocalInteraction();
}

function renderPlayers(){
  $('players').innerHTML=`<h2>Players</h2><div class="stats">${view.players.map(player=>{
    const fighters=view.fighters.filter(fighter=>fighter.owner_id===player.id);
    return `<div class="stat"><strong class="${player.id===view.current_player_id?'current':''}">${escapeHTML(player.name)}</strong>${fighters.map(fighter=>`<div>${escapeHTML(fighter.name)}: ${fighter.health}/${fighter.max_health} HP${fighter.defeated?' · defeated':''}</div>`).join('')}<div>Deck ${player.deck_count} · Hand ${player.hand_count}</div><div>Discard ${player.discard_count}</div></div>`;
  }).join('')}</div>`;
}
function renderPending(){
  const box=$('pending');
  const identity=pendingIdentity(view.pending);
  if(identity===renderedPendingIdentity)return;
  renderedPendingIdentity=identity;
  if(!view.pending){box.hidden=true;box.innerHTML='';return}
  box.hidden=false;
  const mine=view.pending.owner_id===view.viewing_player_id;
  let options=[];
  if(mine){
    options=isMapPendingKind(view.pending.kind)
      ?pendingMapModel(view.pending).specialOptions
      :(view.pending.options||[]);
  }
  box.innerHTML=`<h2>Pending choice</h2><p>${escapeHTML(view.pending.message)}</p>${mine
    ?`${isMapPendingKind(view.pending.kind)?'<p class="map-choice-note">Use the highlighted fighters and spaces on the map.</p>':''}<div class="prompt-options">${options.map(option=>`<button data-choice="${escapeAttr(option.id)}">${escapeHTML(option.label)}</button>`).join('')}</div>`
    :`<p>Waiting for ${escapeHTML(view.pending.owner_name)}.</p>`}`;
}
function renderCombat(){
  const box=$('combat');
  if(!box)return;
  if(!view.combat){box.hidden=true;box.innerHTML='';return}
  box.hidden=false;
  box.innerHTML=`<h2>Combat</h2><p>${view.combat.waiting_for_defense?'Defense choice is hidden until reveal.':'Combat cards revealed.'}</p>`;
}
function selectedOwnFighterID(){
  const id=interactionValue('fighter_id');
  const fighter=fighterByID(id);
  return fighter&&fighter.owner_id===view.viewing_player_id&&!fighter.defeated?id:'';
}
function renderActions(){
  const box=$('actions');
  const mine=view.current_player_id===view.viewing_player_id&&!view.pending&&view.phase==='active';
  const selected=selectedOwnFighterID();
  const selectedName=selected?fighterByID(selected)?.name:'No fighter selected';
  const attackModel=attackMapModel(view.legal,selected);
  const selectedCanAttack=selected&&attackModel.selectedFighterID===selected;
  box.innerHTML=`<h2>Actions</h2><p class="selection-summary">${escapeHTML(selectedName)}</p><div class="actions">
    <button id="maneuver" ${!mine||!view.legal.can_maneuver?'disabled':''}>Maneuver</button>
    <button id="scheme" ${!mine||!view.legal.scheme_cards?.length?'disabled':''}>Scheme</button>
    <button id="attack" ${!mine||!attackMapModel(view.legal).fighters.length?'disabled':''}>Attack${selectedCanAttack?' with selected fighter':''}</button>
  </div>`;
  $('maneuver').onclick=startManeuver;
  $('scheme').onclick=openSchemeCardPicker;
  $('attack').onclick=startAttackInteraction;
}
function renderHand(){
  const mine=view.players.find(player=>player.id===view.viewing_player_id);
  $('hand').innerHTML=visibleCardArtCards(view).map(card=>`<article class="card hand-card">${cardArtMarkup(card,`<strong>${escapeHTML(card.name)}</strong><small>${card.type.toUpperCase()} ${card.type==='scheme'?'':card.value} · BOOST ${card.boost}</small><small>${escapeHTML(card.effect.replaceAll('_',' '))}</small>`)}</article>`).join('')||'<p>No cards in hand.</p>';
}
function renderEvents(){
  const events=[...view.events].reverse().slice(0,35);
  $('events').innerHTML=events.map(event=>`<li>${escapeHTML(event.message)}</li>`).join('');
}

async function startManeuver(){
  preferredManeuverFighterID=selectedOwnFighterID()||null;
  await command({type:'maneuver'});
}
function startAttackInteraction(){
  const selected=selectedOwnFighterID();
  const model=attackMapModel(view.legal,selected);
  if(selected&&model.selectedFighterID){
    beginInteraction('attack_target','attack',{fighter_id:selected});
  }else{
    beginInteraction('attack_fighter','attack');
  }
  renderLocalInteraction();
}
function startSchemeMap(cardID){
  const action=view.legal.scheme_actions_by_card?.[cardID];
  const fighters=action?.fighters||[];
  if(!fighters.length)return command({type:'scheme',card_id:cardID});
  const selected=selectedOwnFighterID();
  const selectedIsLegal=fighters.some(candidate=>candidate.fighter_id===selected);
  const fighterID=selectedIsLegal?selected:(fighters.length===1?fighters[0].fighter_id:'');
  beginInteraction(fighterID?'scheme_destination':'scheme_fighter','scheme',{
    card_id:cardID,
    fighter_id:fighterID,
  });
  renderLocalInteraction();
  return true;
}
function openSchemeCardPicker(){
  const me=view.players.find(player=>player.id===view.viewing_player_id);
  const cards=(me?.hand||[]).filter(card=>view.legal.scheme_cards?.includes(card.id));
  openCardPicker('scheme-card','Choose a Scheme card',cards,cardID=>startSchemeMap(cardID));
}
function openAttackCardPicker(attackerID,targetID){
  const cardIDs=view.legal.attack_cards_by_fighter?.[attackerID]||[];
  const cards=cardIDs.map(cardByID).filter(Boolean);
  openCardPicker('attack-card','Choose an attack card',cards,cardID=>command({
    type:'attack',
    fighter_id:attackerID,
    target_id:targetID,
    card_id:cardID,
  }));
}
function openCardPicker(action,title,cards,onPick){
  if(!cards.length)return showError('No legal cards are available.');
  beginInteraction('dialog',action,{primary:cards[0].id});
  const dialog=$('action-dialog'),body=$('dialog-body');
  $('dialog-title').textContent=title;
  body.innerHTML=`<label>Card<select name="primary">${cards.map(card=>`<option value="${card.id}">${escapeHTML(card.name)}${card.type==='scheme'?'':` (${card.value})`}</option>`).join('')}</select></label>`;
  const primary=body.querySelector('[name=primary]');
  primary.onchange=()=>rememberInteraction('primary',primary.value);
  dialog.showModal();
  $('dialog-submit').onclick=async event=>{
    event.preventDefault();
    try{
      const accepted=await onPick(primary.value);
      if(accepted!==false){
        if(dialog.open)dialog.close();
        if(interaction.mode==='dialog')resetInteraction();
      }
    }catch(error){showError(error.message)}
  };
  dialog.oncancel=resetInteraction;
}

async function choosePendingDestination(spaceID){
  const model=pendingMapModel(view.pending,interactionValue('fighter_id'));
  const option=model.destinationOptions.find(candidate=>candidate.destination===spaceID);
  if(!option)return;
  if(view.pending.kind==='maneuver_move')preferredManeuverFighterID=null;
  await command({type:'choose',choice:option.id});
}
async function submitSchemeDestination(spaceID){
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
    renderLocalInteraction();
    return;
  }
  await command({type:'scheme',card_id:cardID,fighter_id:fighterID,path:[...(destination.path||[])]});
}
async function submitSchemeTarget(targetID){
  const cardID=interactionValue('card_id');
  const fighterID=interactionValue('fighter_id');
  const action=view.legal.scheme_actions_by_card?.[cardID];
  const fighterAction=schemeFighterAction(action,fighterID);
  const destinationID=interactionValue('destination');
  const destination=schemeDestinationOption(fighterAction,destinationID);
  const targets=schemeTargetIDs(fighterAction,destinationID);
  if(!destination||!targets.includes(targetID))return;
  rememberInteraction('target_id',targetID);
  await command({
    type:'scheme',
    card_id:cardID,
    fighter_id:fighterID,
    target_id:targetID,
    path:[...(destination.path||[])],
  });
}

function handleBoardSelection(spaceID,fighterID=''){
  const fighter=fighterByID(fighterID);
  switch(interaction.mode){
  case 'idle':
  case 'fighter_selected':
    if(fighter&&fighter.owner_id===view.viewing_player_id&&!fighter.defeated){
      beginInteraction('fighter_selected','select',{fighter_id:fighter.id});
      renderLocalInteraction();
    }
    break;
  case 'attack_fighter':{
    const model=attackMapModel(view.legal);
    if(model.fighters.includes(fighterID)){
      beginInteraction('attack_target','attack',{fighter_id:fighterID});
      renderLocalInteraction();
    }
    break;
  }
  case 'attack_target':{
    const model=attackMapModel(view.legal,interactionValue('fighter_id'));
    if(model.fighters.includes(fighterID)){
      rememberInteraction('fighter_id',fighterID);
      renderLocalInteraction();
    }else if(model.targetIDs.includes(fighterID)){
      rememberInteraction('target_id',fighterID);
      openAttackCardPicker(interactionValue('fighter_id'),fighterID);
    }
    break;
  }
  case 'scheme_fighter':{
    const action=view.legal.scheme_actions_by_card?.[interactionValue('card_id')];
    if(schemeMapModel(action).fighters.includes(fighterID)){
      interaction.mode='scheme_destination';
      rememberInteraction('fighter_id',fighterID);
      renderLocalInteraction();
    }
    break;
  }
  case 'scheme_destination':
    void submitSchemeDestination(spaceID);
    break;
  case 'scheme_target':{
    const action=view.legal.scheme_actions_by_card?.[interactionValue('card_id')];
    const model=schemeMapModel(action,interactionValue('fighter_id'),interactionValue('destination'));
    if(model.targetIDs.includes(fighterID))void submitSchemeTarget(fighterID);
    break;
  }
  case 'pending_fighter':{
    const model=pendingMapModel(view.pending);
    if(model.fighters.includes(fighterID)){
      interaction.mode='pending_destination';
      rememberInteraction('fighter_id',fighterID);
      renderLocalInteraction();
    }
    break;
  }
  case 'pending_destination':{
    const model=pendingMapModel(view.pending,interactionValue('fighter_id'));
    if(model.fighters.includes(fighterID)&&fighterID!==interactionValue('fighter_id')){
      rememberInteraction('fighter_id',fighterID);
      renderLocalInteraction();
    }else{
      void choosePendingDestination(spaceID);
    }
    break;
  }
  }
}
function boardEventTarget(event){
  const group=event.target?.closest?.('[data-space]');
  if(!group)return null;
  return {spaceID:group.dataset.space||'',fighterID:group.dataset.fighter||''};
}
function handleBoardClick(event){
  const target=boardEventTarget(event);
  if(target)handleBoardSelection(target.spaceID,target.fighterID);
}
function handleBoardKeydown(event){
  if(event.key!=='Enter'&&event.key!==' ')return;
  const target=boardEventTarget(event);
  if(!target)return;
  event.preventDefault();
  handleBoardSelection(target.spaceID,target.fighterID);
}
function handleBoardPointerOver(event){
  const target=boardEventTarget(event);
  if(target)setPathPreview(target.spaceID);
}
function handleBoardPointerOut(event){
  const from=event.target?.closest?.('[data-space]');
  const to=event.relatedTarget?.closest?.('[data-space]');
  if(from&&from!==to)setPathPreview('');
}

function addOptionalHornsTarget(payload,targetSelect){
  if(!targetSelect)return payload;
  if(!targetSelect.value)throw new Error('Choose a living fighter adjacent to Jackalope’s destination.');
  payload.target_id=targetSelect.value;
  return payload;
}
function initials(name){return name.split(/\s+/).map(part=>part[0]).join('').slice(0,2).toUpperCase()}
function escapeHTML(value){return String(value??'').replace(/[&<>\"]/g,ch=>({'&':'&amp;','<':'&lt;','>':'&gt;','\"':'&quot;'}[ch]))}
function escapeAttr(value){return escapeHTML(value).replace(/'/g,'&#39;')}

$('pending').onclick=event=>{
  const button=event.target?.closest?.('[data-choice]');
  if(button)void command({type:'choose',choice:button.dataset.choice});
};
$('board').onclick=handleBoardClick;
$('board').onkeydown=handleBoardKeydown;
$('board').onpointerover=handleBoardPointerOver;
$('board').onpointerout=handleBoardPointerOut;
$('create').onclick=createMatch;
$('join').onclick=joinMatch;
$('reset').onclick=leave;
if(session){enterGame();void refresh()}
