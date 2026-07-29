const $ = (id) => document.getElementById(id);
const lobby = $('lobby'), game = $('game'), errorBox = $('error');
const zoneColors = {gray:'#a9aaa5','light-gray':'#d7d6c9',green:'#598b5b',brown:'#9a7655','light-green':'#8ebf78',orange:'#d39a50',yellow:'#d8c45e'};
let session = JSON.parse(sessionStorage.getItem('unmatched-session') || 'null');
let view = null;
let pollHandle = null;
let commandInFlight = false;
let refreshInFlight = false;
let renderedPendingIdentity = null;
const interaction = {mode:'idle',action:null,revision:null,values:{}};

function saveSession(next){session=next;if(next)sessionStorage.setItem('unmatched-session',JSON.stringify(next));else sessionStorage.removeItem('unmatched-session')}
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
function shouldApplyView(current,next,force=false){return Boolean(force||!current||!next||current.match_id!==next.match_id||current.revision!==next.revision||current.phase!==next.phase||current.viewing_player_id!==next.viewing_player_id)}
function interactionCanContinue(state,next){return Boolean(state?.mode==='dialog'&&next&&state.revision===next.revision&&next.phase==='active'&&!next.pending&&next.current_player_id===next.viewing_player_id)}
function resetInteraction(){interaction.mode='idle';interaction.action=null;interaction.revision=null;interaction.values={}}
function beginInteraction(action){if(interaction.action!==action||interaction.revision!==view?.revision)resetInteraction();interaction.mode='dialog';interaction.action=action;interaction.revision=view?.revision??null}
function rememberInteraction(name,value){interaction.values[name]=value}
function interactionValue(name,fallback=''){return interaction.values[name]??fallback}
function reconcileInteraction(next){const dialog=$('action-dialog');if(!dialog?.open)return;if(interactionCanContinue(interaction,next))return;if(typeof dialog.close==='function')dialog.close();resetInteraction()}
function applyView(next,options={}){if(!next||!shouldApplyView(view,next,Boolean(options.force)))return false;view=next;reconcileInteraction(next);render();return true}
async function createMatch(){try{const body=await request('/api/matches',{method:'POST',body:'{}'});saveSession({matchID:body.match_id,token:body.token,playerID:body.player_id});applyView(body.view,{force:true});enterGame()}catch(error){showError(error.message)}}
async function joinMatch(){const code=$('join-code').value.trim();if(!code)return showError('Enter a match code.');try{const body=await request(`/api/matches/${encodeURIComponent(code)}/join`,{method:'POST',body:'{}'});saveSession({matchID:body.match_id,token:body.token,playerID:body.player_id});applyView(body.view,{force:true});enterGame()}catch(error){showError(error.message)}}
async function refresh(){if(!session||commandInFlight||refreshInFlight)return false;refreshInFlight=true;try{const body=await request(`/api/matches/${encodeURIComponent(session.matchID)}`);return applyView(body.view)}catch(error){showError(error.message);return false}finally{refreshInFlight=false}}
function beginCommand(){if(commandInFlight)return false;commandInFlight=true;return true}
function endCommand(){commandInFlight=false}
async function command(payload){if(!beginCommand())return false;let failed=false;try{const body=await request(`/api/matches/${encodeURIComponent(session.matchID)}/commands`,{method:'POST',body:JSON.stringify({...payload,expected_revision:view.revision})});applyView(body.view);return true}catch(error){showError(error.message);failed=true;return false}finally{endCommand();if(failed)void refresh()}}
function enterGame(){lobby.hidden=true;game.hidden=false;$('reset').hidden=false;if(pollHandle)clearInterval(pollHandle);pollHandle=setInterval(refresh,1200)}
function leave(){if(pollHandle)clearInterval(pollHandle);pollHandle=null;saveSession(null);view=null;refreshInFlight=false;renderedPendingIdentity=null;resetInteraction();game.hidden=true;lobby.hidden=false;$('reset').hidden=true;clearError()}

function render(){if(!view)return;renderStatus();renderBoard();renderPlayers();renderPending();renderActions();renderHand();renderEvents()}
function renderStatus(){const current=view.players.find(p=>p.id===view.current_player_id);const mine=view.viewing_player_id===view.current_player_id;const status=view.phase==='ended'?`<span class="winner">${escapeHTML(view.winner_name)} wins</span>`:view.phase==='waiting_for_player'?`Waiting for Bigfoot. Match code: <strong>${escapeHTML(view.match_id)}</strong>`:`${mine?'<span class="current">Your turn</span>':`${escapeHTML(current?.name||'Opponent')}'s turn`} · ${current?.actions_remaining??0} actions remaining`;$('status').innerHTML=`<div>${status}</div><small>Match <strong>${escapeHTML(view.match_id)}</strong> · revision ${view.revision}</small>`}
function renderBoard(){const svg=$('board');const spaces=new Map(view.spaces.map(space=>[space.id,space]));const edges=view.edges.map(edge=>{const a=spaces.get(edge.from),b=spaces.get(edge.to);return `<line class="edge" x1="${a.x}" y1="${a.y}" x2="${b.x}" y2="${b.y}"/>`}).join('');const nodes=view.spaces.map(space=>{const colors=space.zones.map(z=>zoneColors[z]||'#999');const fill=colors[0];const fighter=space.fighter;return `<g data-space="${space.id}"><circle class="space" cx="${space.x}" cy="${space.y}" r="3.5" fill="${fill}"><title>${space.id} · ${space.zones.join(', ')}</title></circle><text class="space-label" x="${space.x}" y="${space.y+.8}">${space.id.slice(1)}</text>${fighter?`<circle class="fighter" cx="${space.x}" cy="${space.y-5}" r="3" fill="${fighter.owner_id===view.viewing_player_id?'#3f87d8':'#b94b43'}"><title>${escapeHTML(fighter.name)} ${fighter.health}/${fighter.max_health}</title></circle><text class="fighter-label" x="${space.x}" y="${space.y-4.3}">${initials(fighter.name)}</text>`:''}</g>`}).join('');svg.innerHTML=edges+nodes}
function renderPlayers(){$('players').innerHTML=`<h2>Players</h2><div class="stats">${view.players.map(player=>{const fighters=view.fighters.filter(f=>f.owner_id===player.id);return `<div class="stat"><strong class="${player.id===view.current_player_id?'current':''}">${escapeHTML(player.name)}</strong>${fighters.map(f=>`<div>${escapeHTML(f.name)}: ${f.health}/${f.max_health} HP${f.defeated?' · defeated':''}</div>`).join('')}<div>Deck ${player.deck_count} · Hand ${player.hand_count}</div><div>Discard ${player.discard_count}</div></div>`}).join('')}</div>`}
function pendingIdentity(pending){if(!pending)return 'none';return JSON.stringify([pending.kind,pending.owner_id,pending.message,(pending.options||[]).map(option=>[option.id,option.label])])}
function renderPending(){const box=$('pending');const identity=pendingIdentity(view.pending);if(identity===renderedPendingIdentity)return;renderedPendingIdentity=identity;if(!view.pending){box.hidden=true;box.innerHTML='';return}box.hidden=false;const mine=view.pending.owner_id===view.viewing_player_id;box.innerHTML=`<h2>Pending choice</h2><p>${escapeHTML(view.pending.message)}</p>${mine?`<div class="prompt-options">${view.pending.options.map(option=>`<button data-choice="${escapeAttr(option.id)}">${escapeHTML(option.label)}</button>`).join('')}</div>`:`<p>Waiting for ${escapeHTML(view.pending.owner_name)}.</p>`}`}
function renderActions(){const box=$('actions');const mine=view.current_player_id===view.viewing_player_id&&!view.pending&&view.phase==='active';box.innerHTML=`<h2>Actions</h2><div class="actions"><button id="maneuver" ${!mine||!view.legal.can_maneuver?'disabled':''}>Maneuver</button><button id="scheme" ${!mine||!view.legal.scheme_cards?.length?'disabled':''}>Scheme</button><button id="attack" ${!mine||!Object.keys(view.legal.attack_cards_by_fighter||{}).length?'disabled':''}>Attack</button></div>`;$('maneuver').onclick=()=>command({type:'maneuver'});$('scheme').onclick=openScheme;$('attack').onclick=openAttack}
function renderHand(){const mine=view.players.find(p=>p.id===view.viewing_player_id);$('hand').innerHTML=(mine?.hand||[]).map(card=>`<article class="card"><strong>${escapeHTML(card.name)}</strong><small>${card.type.toUpperCase()} ${card.type==='scheme'?'':card.value} · BOOST ${card.boost}</small><small>${escapeHTML(card.effect.replaceAll('_',' '))}</small></article>`).join('')||'<p>No cards in hand.</p>'}
function renderEvents(){const events=[...view.events].reverse().slice(0,35);$('events').innerHTML=events.map(event=>`<li>${escapeHTML(event.message)}</li>`).join('')}

function schemeFighterAction(action,fighterID){return (action?.fighters||[]).find(candidate=>candidate.fighter_id===fighterID)}
function schemeDestinationOption(fighterAction,destinationID){return (fighterAction?.destinations||[]).find(option=>option.destination===destinationID)}
function schemeTargetIDs(fighterAction,destinationID){return fighterAction?.targets_by_destination?.[destinationID]||[]}
function renderSchemeExtras(card,action,extras){
  const fighters=action?.fighters||[];
  if(!fighters.length){extras.innerHTML='';return}
  const rememberedFighter=interactionValue('fighter_id');
  const fighterID=fighters.some(candidate=>candidate.fighter_id===rememberedFighter)?rememberedFighter:fighters[0].fighter_id;
  const fighterAction=schemeFighterAction(action,fighterID);
  const destinations=fighterAction?.destinations||[];
  const rememberedDestination=interactionValue('destination');
  const destinationID=destinations.some(option=>option.destination===rememberedDestination)?rememberedDestination:destinations[0]?.destination||'';
  rememberInteraction('fighter_id',fighterID);rememberInteraction('destination',destinationID);
  const fighterField=fighters.length>1?`<label>Fighter<select name="fighter">${fighters.map(candidate=>{const fighter=fighterByID(candidate.fighter_id);return `<option value="${candidate.fighter_id}">${escapeHTML(fighter?.name||candidate.fighter_id)} at ${escapeHTML(fighter?.space_id||'')}</option>`}).join('')}</select></label>`:'';
  extras.innerHTML=`${fighterField}<label>Destination<select name="destination">${destinations.map(option=>`<option value="${option.destination}">${escapeHTML(option.label||option.destination)}</option>`).join('')}</select></label>${card.effect==='horns'?'<div data-horns-target></div>':''}`;
  const fighterSelect=extras.querySelector('[name=fighter]');if(fighterSelect){fighterSelect.value=fighterID;fighterSelect.onchange=()=>{rememberInteraction('fighter_id',fighterSelect.value);rememberInteraction('destination','');rememberInteraction('target_id','');renderSchemeExtras(card,action,extras)}}
  const destinationSelect=extras.querySelector('[name=destination]');if(destinationSelect){destinationSelect.value=destinationID;destinationSelect.onchange=()=>{rememberInteraction('destination',destinationSelect.value);rememberInteraction('target_id','');renderSchemeExtras(card,action,extras)}}
  if(card.effect!=='horns')return;
  const targetBox=extras.querySelector('[data-horns-target]');
  const targets=schemeTargetIDs(fighterAction,destinationID);
  if(!targets.length){targetBox.innerHTML='<p>No adjacent fighter; the damage step is skipped.</p>';return}
  const rememberedTarget=interactionValue('target_id');
  const targetID=targets.includes(rememberedTarget)?rememberedTarget:targets[0];
  rememberInteraction('target_id',targetID);
  targetBox.innerHTML=`<label>Adjacent target<select name="target">${targets.map(id=>{const candidate=fighterByID(id);return `<option value="${id}">${escapeHTML(candidate?.name||id)} at ${escapeHTML(candidate?.space_id||'')}</option>`}).join('')}</select></label>`;
  const targetSelect=targetBox.querySelector('[name=target]');targetSelect.value=targetID;targetSelect.onchange=()=>rememberInteraction('target_id',targetSelect.value);
}
function openScheme(){
  const me=view.players.find(p=>p.id===view.viewing_player_id);
  const cards=me.hand.filter(card=>view.legal.scheme_cards.includes(card.id));
  const actions=view.legal.scheme_actions_by_card||{};
  openDialog('scheme','Play a scheme',cards.map(card=>`<option value="${card.id}">${escapeHTML(card.name)}</option>`).join(''),async(cardID,body)=>{
    const card=cards.find(candidate=>candidate.id===cardID);
    const payload={type:'scheme',card_id:cardID};
    const action=actions[cardID];
    if(action?.fighters?.length){
      const fighterID=body.querySelector('[name=fighter]')?.value||action.fighters[0].fighter_id;
      const fighterAction=schemeFighterAction(action,fighterID);
      const destinationID=body.querySelector('[name=destination]')?.value;
      const destination=schemeDestinationOption(fighterAction,destinationID);
      if(!destination)throw new Error('Choose a legal destination supplied by the server.');
      payload.fighter_id=fighterID;payload.path=[...(destination.path||[])];
      if(card.effect==='horns')addOptionalHornsTarget(payload,body.querySelector('[name=target]'));
    }
    return command(payload);
  },(select,body)=>{
    const update=()=>{rememberInteraction('primary',select.value);const card=cards.find(candidate=>candidate.id===select.value);renderSchemeExtras(card,actions[card.id],body.querySelector('.extras'))};
    select.onchange=()=>{rememberInteraction('fighter_id','');rememberInteraction('destination','');rememberInteraction('target_id','');update()};update();
  });
}
function openAttack(){const attackers=Object.keys(view.legal.attack_cards_by_fighter||{});const options=attackers.map(id=>{const fighter=fighterByID(id);return `<option value="${id}">${escapeHTML(fighter.name)} at ${fighter.space_id}</option>`}).join('');openDialog('attack','Attack',options,async(attackerID,body)=>command({type:'attack',fighter_id:attackerID,target_id:body.querySelector('[name=target]').value,card_id:body.querySelector('[name=card]').value}),(select,body)=>{const update=()=>{rememberInteraction('primary',select.value);const attackerID=select.value;const targetIDs=view.legal.attack_targets_by_fighter[attackerID]||[];const cardIDs=view.legal.attack_cards_by_fighter[attackerID]||[];const rememberedTarget=interactionValue('target_id');const rememberedCard=interactionValue('card_id');const targetID=targetIDs.includes(rememberedTarget)?rememberedTarget:targetIDs[0];const cardID=cardIDs.includes(rememberedCard)?rememberedCard:cardIDs[0];body.querySelector('.extras').innerHTML=`<label>Target<select name="target">${targetIDs.map(id=>{const f=fighterByID(id);return `<option value="${id}">${escapeHTML(f.name)} · ${f.health} HP · ${f.space_id}</option>`}).join('')}</select></label><label>Attack card<select name="card">${cardIDs.map(id=>{const c=cardByID(id);return `<option value="${id}">${escapeHTML(c.name)} (${c.value})</option>`}).join('')}</select></label>`;const target=body.querySelector('[name=target]'),card=body.querySelector('[name=card]');target.value=targetID;card.value=cardID;rememberInteraction('target_id',targetID);rememberInteraction('card_id',cardID);target.onchange=()=>rememberInteraction('target_id',target.value);card.onchange=()=>rememberInteraction('card_id',card.value)};select.onchange=()=>{rememberInteraction('target_id','');rememberInteraction('card_id','');update()};update()})}
function openDialog(action,title,primaryOptions,submit,configure){beginInteraction(action);const dialog=$('action-dialog'),body=$('dialog-body');$('dialog-title').textContent=title;body.innerHTML=`<label>Selection<select name="primary">${primaryOptions}</select></label><div class="extras"></div>`;const primary=body.querySelector('[name=primary]');const remembered=interactionValue('primary');if(Array.from(primary.options||[]).some(option=>option.value===remembered))primary.value=remembered;rememberInteraction('primary',primary.value);configure(primary,body);dialog.showModal();$('dialog-submit').onclick=async(event)=>{event.preventDefault();try{const accepted=await submit(primary.value,body);if(accepted!==false){if(dialog.open)dialog.close();resetInteraction()}}catch(error){showError(error.message)}};dialog.oncancel=resetInteraction}

function allFighters(){return view.fighters||view.spaces.flatMap(space=>space.fighter?[space.fighter]:[])}
function fighterByID(id){return allFighters().find(f=>f.id===id)}
function cardByID(id){return view.players.find(p=>p.id===view.viewing_player_id).hand.find(c=>c.id===id)}
function addOptionalHornsTarget(payload,targetSelect){if(!targetSelect)return payload;if(!targetSelect.value)throw new Error('Choose a living fighter adjacent to Jackalope’s destination.');payload.target_id=targetSelect.value;return payload}
function initials(name){return name.split(/\s+/).map(part=>part[0]).join('').slice(0,2).toUpperCase()}
function escapeHTML(value){return String(value??'').replace(/[&<>\"]/g,ch=>({'&':'&amp;','<':'&lt;','>':'&gt;','\"':'&quot;'}[ch]))}
function escapeAttr(value){return escapeHTML(value).replace(/'/g,'&#39;')}

$('pending').onclick=(event)=>{const button=event.target?.closest?.('[data-choice]');if(button)void command({type:'choose',choice:button.dataset.choice})};
$('create').onclick=createMatch;$('join').onclick=joinMatch;$('reset').onclick=leave;
if(session){enterGame();void refresh()}
