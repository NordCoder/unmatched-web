const STAGE4_CARD_RULES = Object.freeze({
  'a-hunters-eye':'No special effect.',
  'steal-from-the-rich':'Draw 1 card. Choose an opponent: they may discard 1 card. If they do not, draw 1 card.',
  'disarming-shot':'After combat: Draw cards equal to the combat damage dealt to the opposing fighter.',
  'piercing-shot':'After combat: Draw 2 cards.',
  'highway-robbery':'Immediately: Cancel all effects on your opponent’s card and ignore its value.',
  'defenders-of-sherwood':'After combat: Draw 1 card. If an Outlaw is defeated, return one at full health to an empty space in Robin Hood’s zone.',
  'feint':'Immediately: Cancel all effects on your opponent’s card.',
  'regroup':'After combat: If you won, draw 2 cards. Otherwise, draw 1 card.',
  'wily-fighting':'After combat: Deal 1 damage to each opposing fighter adjacent to your fighter.',
  'snark':'After combat: If the combat fighters are adjacent, draw 1 card.',
  'ambush':'During combat: Your opponent discards 1 random card. Add its BOOST value to this card.',
  'larger-than-life':'No special effect.',
  'savagery':'After combat: If you won, deal 1 damage to each opposing fighter adjacent to Bigfoot.',
  'crash-through-the-trees':'Move Bigfoot up to 5 spaces. He may move through opposing fighters.',
  'jackalope-horns':'Move the Jackalope up to 5 spaces. It may move through opposing fighters. If another fighter is adjacent after moving, choose one and deal 2 damage to it.',
  'hoax':'After combat: Move your combat fighter up to 5 spaces. It may move through opposing fighters.',
  'disengage':'After combat: Place your combat fighter in an empty space in its current zone.',
  'its-just-your-imagination':'Immediately: Cancel all effects on your opponent’s card.',
  'skirmish':'After combat: If you won, move either combat fighter up to 2 spaces.',
  'momentous-shift':'During combat: If your fighter started the turn in a different space, this card’s value is 5.',
});

function stage4UserFacingMessage(message){
  const text=String(message??'').replace(/\bscheme\b/gi,match=>match[0]===match[0].toUpperCase()?'Action':'action');
  return text.replace(/\ba Action\b/g,'an Action').replace(/\ba action\b/g,'an action');
}
function stage4AttackTypeLabel(value){return value==='ranged'?'Ranged':'Melee'}
function stage4CardTypeLabel(card){return card?.type==='scheme'?'ACTION':String(card?.type||'').toUpperCase()}
function stage4FighterDefinitionName(definitionID,fighters=allFighters()){
  const match=(fighters||[]).find(fighter=>fighter.definition_id===definitionID);
  if(match)return match.name;
  return String(definitionID||'').split('-').map(part=>part?part[0].toUpperCase()+part.slice(1):'').join(' ');
}
function stage4UsableByLabel(card,fighters=allFighters()){
  const usable=card?.usable_by||[];
  if(!usable.length||usable.includes('any'))return 'Any fighter';
  return [...new Set(usable.map(id=>stage4FighterDefinitionName(id,fighters)).filter(Boolean))].join(', ');
}
function stage4CardRulesText(card){return STAGE4_CARD_RULES[card?.definition_id]||'No special effect.'}
function stage4CardStatLine(card){return `${card?.type==='scheme'?'':`VALUE ${card?.value??0} · `}BOOST ${card?.boost??0}`}
function stage4CardDetailsContent(card){
  return `<div class="card-head"><strong>${escapeHTML(card.name)}</strong><span class="card-type">${escapeHTML(stage4CardTypeLabel(card))}</span></div><div class="card-meta">${escapeHTML(stage4CardStatLine(card))}</div><div class="card-usable">Usable by: ${escapeHTML(stage4UsableByLabel(card))}</div><p class="card-rules">${escapeHTML(stage4CardRulesText(card))}</p>`;
}
function stage4ManeuverFighterModel(fighterIDs,selectedFighterID=''){
  const fighters=[...new Set(fighterIDs||[])];
  const selected=fighters.includes(selectedFighterID)?selectedFighterID:'';
  return {fighters,selectedFighterID:selected||(fighters.length===1?fighters[0]:'')};
}
function stage4ManeuverBoostModel(pending){
  const options=pending?.kind==='maneuver_boost'?(pending.options||[]):[];
  return {noBoost:options.find(option=>!option.card_id)||null,boostOptions:options.filter(option=>option.card_id)};
}

showError=function(message){errorBox.textContent=stage4UserFacingMessage(message);errorBox.hidden=false};

const stage3BoardHighlights=boardHighlights;
boardHighlights=function(){
  const result=stage3BoardHighlights();
  if(interaction.mode==='maneuver_fighter')ownLivingFighterIDs().forEach(id=>result.fighterCandidates.add(id));
  return result;
};

const stage3InteractionInstruction=interactionInstruction;
interactionInstruction=function(){
  if(interaction.mode==='fighter_selected'){
    const fighter=fighterByID(interactionValue('fighter_id'));
    return `${fighter?.name||'Fighter'} selected. Choose Maneuver, Action, or Attack.`;
  }
  if(interaction.mode==='maneuver_fighter')return 'Choose the highlighted fighter that will move first during this Maneuver.';
  if(interaction.mode==='scheme_fighter'){
    const card=cardByID(interactionValue('card_id'));
    return `Choose a highlighted fighter to use ${card?.name||'the Action'}.`;
  }
  return stage3InteractionInstruction();
};

renderPlayers=function(){
  $('players').innerHTML=`<h2>Players</h2><div class="stats">${view.players.map(player=>{
    const fighters=view.fighters.filter(fighter=>fighter.owner_id===player.id);
    return `<div class="stat"><strong class="${player.id===view.current_player_id?'current':''}">${escapeHTML(player.name)}</strong>${fighters.map(fighter=>`<div class="fighter-row"><span>${escapeHTML(fighter.name)}: ${fighter.health}/${fighter.max_health} HP${fighter.defeated?' · defeated':''}</span><span class="range-badge">${escapeHTML(stage4AttackTypeLabel(fighter.attack_type))}</span></div>`).join('')}<div>Deck ${player.deck_count} · Hand ${player.hand_count}</div><div>Discard ${player.discard_count}</div></div>`;
  }).join('')}</div>`;
};
function stage4PendingOptionMarkup(option){
  const card=option.card_id?cardByID(option.card_id):null;
  if(!card)return `<button data-choice="${escapeAttr(option.id)}">${escapeHTML(option.label)}</button>`;
  return `<button class="pending-card-choice" data-choice="${escapeAttr(option.id)}">${stage4CardDetailsContent(card)}</button>`;
}
function stage4ManeuverBoostMarkup(pending){
  const model=stage4ManeuverBoostModel(pending);
  const fighter=fighterByID(preferredManeuverFighterID);
  return `<div class="maneuver-boost-copy"><strong>${fighter?`${escapeHTML(fighter.name)} moves first`:'Choose movement strength'}</strong><span>Continue with normal movement, or discard one card for its BOOST value.</span></div><div class="pending-actions">${model.noBoost?`<button data-choice="${escapeAttr(model.noBoost.id)}">Continue without BOOST</button>`:''}<button data-open-boost ${model.boostOptions.length?'':'disabled'}>BOOST</button></div>`;
}
renderPending=function(){
  const box=$('pending');
  const identity=pendingIdentity(view.pending);
  if(identity===renderedPendingIdentity)return;
  renderedPendingIdentity=identity;
  if(!view.pending){box.hidden=true;box.innerHTML='';return}
  box.hidden=false;
  const mine=view.pending.owner_id===view.viewing_player_id;
  let options=[];
  if(mine){
    options=isMapPendingKind(view.pending.kind)?pendingMapModel(view.pending).specialOptions:(view.pending.options||[]);
  }
  const controls=mine&&view.pending.kind==='maneuver_boost'
    ?stage4ManeuverBoostMarkup(view.pending)
    :mine
      ?`${isMapPendingKind(view.pending.kind)?'<p class="map-choice-note">Use the highlighted fighters and spaces on the map.</p>':''}<div class="prompt-options">${options.map(stage4PendingOptionMarkup).join('')}</div>`
      :`<p>Waiting for ${escapeHTML(view.pending.owner_name)}.</p>`;
  box.innerHTML=`<h2>Pending choice</h2><p>${escapeHTML(stage4UserFacingMessage(view.pending.message))}</p>${controls}`;
};
renderActions=function(){
  const box=$('actions');
  const mine=view.current_player_id===view.viewing_player_id&&!view.pending&&view.phase==='active';
  const selected=selectedOwnFighterID();
  const selectedName=selected?fighterByID(selected)?.name:'No fighter selected';
  const attackModel=attackMapModel(view.legal,selected);
  const selectedCanAttack=selected&&attackModel.selectedFighterID===selected;
  box.innerHTML=`<h2>Actions</h2><p class="selection-summary">${escapeHTML(selectedName)}</p><div class="actions">
    <button id="maneuver" ${!mine||!view.legal.can_maneuver?'disabled':''}>Maneuver${selected?' with selected fighter':''}</button>
    <button id="scheme" ${!mine||!view.legal.scheme_cards?.length?'disabled':''}>Action</button>
    <button id="attack" ${!mine||!attackMapModel(view.legal).fighters.length?'disabled':''}>Attack${selectedCanAttack?' with selected fighter':''}</button>
  </div>`;
  $('maneuver').onclick=startManeuver;
  $('scheme').onclick=openSchemeCardPicker;
  $('attack').onclick=startAttackInteraction;
};
renderHand=function(){
  const mine=view.players.find(player=>player.id===view.viewing_player_id);
  $('hand').innerHTML=(mine?.hand||[]).map(card=>`<article class="card">${stage4CardDetailsContent(card)}</article>`).join('')||'<p>No cards in hand.</p>';
};
renderEvents=function(){
  const events=[...view.events].reverse().slice(0,35);
  $('events').innerHTML=events.map(event=>`<li>${escapeHTML(stage4UserFacingMessage(event.message))}</li>`).join('');
};

startManeuver=function(){
  const model=stage4ManeuverFighterModel(ownLivingFighterIDs(),selectedOwnFighterID());
  if(model.selectedFighterID)return stage4StartManeuverWithFighter(model.selectedFighterID);
  beginInteraction('maneuver_fighter','maneuver');
  renderLocalInteraction();
  return true;
};
async function stage4StartManeuverWithFighter(fighterID){
  preferredManeuverFighterID=fighterID||null;
  return command({type:'maneuver'});
}
openSchemeCardPicker=function(){
  const me=view.players.find(player=>player.id===view.viewing_player_id);
  const cards=(me?.hand||[]).filter(card=>view.legal.scheme_cards?.includes(card.id));
  stage4OpenCardPicker('scheme-card','Choose an Action card',cards,cardID=>startSchemeMap(cardID),{submitLabel:'Use Action'});
};
openAttackCardPicker=function(attackerID,targetID){
  const cardIDs=view.legal.attack_cards_by_fighter?.[attackerID]||[];
  const cards=cardIDs.map(cardByID).filter(Boolean);
  stage4OpenCardPicker('attack-card','Choose an attack card',cards,cardID=>command({
    type:'attack',fighter_id:attackerID,target_id:targetID,card_id:cardID,
  }),{submitLabel:'Attack'});
};
function stage4OpenManeuverBoostPicker(){
  const model=stage4ManeuverBoostModel(view.pending);
  const choices=new Map(model.boostOptions.map(option=>[option.card_id,option.id]));
  const cards=model.boostOptions.map(option=>cardByID(option.card_id)).filter(Boolean);
  stage4OpenCardPicker('maneuver-boost','Choose a BOOST card',cards,cardID=>{
    const choice=choices.get(cardID);
    if(!choice)throw new Error('The selected BOOST card is no longer available.');
    return command({type:'choose',choice});
  },{submitLabel:'BOOST',note:'Any card may be discarded for BOOST. “Usable by” applies when playing the card as an Action or combat card.'});
}
function stage4OpenCardPicker(action,title,cards,onPick,options={}){
  if(!cards.length)return showError('No legal cards are available.');
  beginInteraction('dialog',action,{primary:cards[0].id});
  const dialog=$('action-dialog'),body=$('dialog-body');
  $('dialog-title').textContent=title;
  body.innerHTML=`<label>Card<select name="primary">${cards.map(card=>`<option value="${card.id}">${escapeHTML(card.name)}${card.type==='scheme'?'':` (${card.value})`}</option>`).join('')}</select></label>${options.note?`<p class="dialog-note">${escapeHTML(options.note)}</p>`:''}<article class="card card-preview" data-card-preview></article>`;
  const primary=body.querySelector('[name=primary]');
  const preview=body.querySelector('[data-card-preview]');
  const updatePreview=()=>{
    rememberInteraction('primary',primary.value);
    const card=cards.find(candidate=>candidate.id===primary.value);
    preview.innerHTML=card?stage4CardDetailsContent(card):'';
  };
  primary.onchange=updatePreview;
  updatePreview();
  $('dialog-submit').textContent=options.submitLabel||'Use card';
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

const stage3HandleBoardSelection=handleBoardSelection;
handleBoardSelection=function(spaceID,fighterID=''){
  if(interaction.mode==='maneuver_fighter'){
    if(ownLivingFighterIDs().includes(fighterID))void stage4StartManeuverWithFighter(fighterID);
    return;
  }
  return stage3HandleBoardSelection(spaceID,fighterID);
};

$('pending').onclick=event=>{
  const boost=event.target?.closest?.('[data-open-boost]');
  if(boost){stage4OpenManeuverBoostPicker();return}
  const button=event.target?.closest?.('[data-choice]');
  if(button)void command({type:'choose',choice:button.dataset.choice});
};
