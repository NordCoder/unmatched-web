function cardArtEscapeHTML(value){
  return String(value??'').replace(/[&<>\"]/g,ch=>({'&':'&amp;','<':'&lt;','>':'&gt;','\"':'&quot;'}[ch]));
}

function cardArtPath(card){
  const deck=card?.deck_definition_id;
  const definition=card?.definition_id;
  if(!deck||!definition)return null;
  return `/cards/${encodeURIComponent(deck)}/${encodeURIComponent(definition)}.png`;
}

function visibleCardArtCards(view){
  const viewer=view?.viewing_player_id;
  return view?.players?.find(player=>player.id===viewer)?.hand||[];
}

function cardArtMarkup(card,fallbackMarkup,options={}){
  const path=cardArtPath(card);
  const selected=options.selected?' selected':'';
  const cardID=card?.id?` data-card-id="${cardArtEscapeHTML(card.id)}"`:'';
  const image=path
    ?`<img class="card-art-image" data-card-art src="${cardArtEscapeHTML(path)}" alt="${cardArtEscapeHTML(card.name||'Card')}" loading="lazy" decoding="async">`
    :'';
  return `<div class="card-art-shell${selected}${path?'':' no-card-art'}" data-card-art-shell${cardID}>${image}<div class="card-fallback" data-card-fallback>${fallbackMarkup}</div></div>`;
}

function cardArtSelectionID(target){
  const picker=target?.closest?.('[data-card-pick]');
  if(picker?.dataset?.cardPick)return picker.dataset.cardPick;
  const shell=target?.closest?.('[data-card-art-shell]');
  return shell?.dataset?.cardId||null;
}

function cardArtEvent(event){
  const image=event.target;
  if(!image?.matches?.('[data-card-art]'))return;
  const shell=image.closest?.('[data-card-art-shell]');
  if(!shell)return;
  const fallback=shell.querySelector?.('[data-card-fallback]');
  if(event.type==='load'){
    shell.classList.add('art-loaded');
    fallback?.setAttribute('aria-hidden','true');
    return;
  }
  image.remove();
  shell.classList.remove('art-loaded');
  fallback?.removeAttribute('aria-hidden');
}

function installCardArtHandlers(root=document){
  root?.addEventListener?.('load',cardArtEvent,true);
  root?.addEventListener?.('error',cardArtEvent,true);
}

globalThis.cardArtPath=cardArtPath;
globalThis.cardArtMarkup=cardArtMarkup;
globalThis.cardArtSelectionID=cardArtSelectionID;
globalThis.visibleCardArtCards=visibleCardArtCards;
globalThis.cardArtEvent=cardArtEvent;
globalThis.installCardArtHandlers=installCardArtHandlers;
installCardArtHandlers();
