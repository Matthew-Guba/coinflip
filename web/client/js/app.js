const state = {
  userId: 1,                         // берётся из вашей системы аутентификации
  balance: 0,
  side: 'heads',
  series: 0,
  nonce: 1,
  clientSeed: 'seed_' + Math.random().toString(36).slice(2, 10),
  history: [],
  busy: false
};

const $ = (id) => document.getElementById(id);

const els = {
  balance: $('balance'),
  bet: $('bet'),
  playBtn: $('playBtn'),
  coin: $('coin'),
  coef: $('coef'),
  series: $('series'),
  result: $('resultMsg'),
  histList: $('historyList'),
  soundBtn: $('soundToggle'),
  sideHeads: $('sideHeads'),
  sideTails: $('sideTails'),
  half: $('half'),
  double: $('double'),
  max: $('max')
};

// ---- init ----
function init() {
  bindEvents();
  loadBalance();
  loadHistory();
}

function bindEvents() {
  els.playBtn.addEventListener('click', onPlay);

  els.sideHeads.addEventListener('click', () => selectSide('heads'));
  els.sideTails.addEventListener('click', () => selectSide('tails'));

  els.half.addEventListener('click', () => { Sound.click(); changeBet(v => v/2); });
  els.double.addEventListener('click', () => { Sound.click(); changeBet(v => v*2); });
  els.max.addEventListener('click', () => { Sound.click(); els.bet.value = state.balance.toFixed(2); });

  els.soundBtn.addEventListener('click', () => {
    const on = Sound.toggle();
    els.soundBtn.textContent = on ? '🔊' : '🔇';
  });
}

function changeBet(fn) {
  const v = parseFloat(els.bet.value) || 0;
  const next = Math.max(0.01, fn(v));
  els.bet.value = next.toFixed(2);
}

function selectSide(side) {
  Sound.click();
  state.side = side;
  els.sideHeads.classList.toggle('active', side === 'heads');
  els.sideTails.classList.toggle('active', side === 'tails');
}

// ---- balance ----
async function loadBalance() {
  try {
    const r = await API.balance(state.userId);
    state.balance = parseFloat(r.balance);
    updateBalance();
  } catch (e) {
    console.warn('balance fail', e);
  }
}

function updateBalance() {
  els.balance.textContent = '$' + state.balance.toFixed(2);
}

// ---- play ----
async function onPlay() {
  if (state.busy) return;
  const bet = parseFloat(els.bet.value);
  if (!bet || bet <= 0) return flashMsg('Введите ставку', 'lose');
  if (bet > state.balance) return flashMsg('Недостаточно средств', 'lose');

  state.busy = true;
  els.playBtn.disabled = true;
  els.result.textContent = '';
  els.result.className = 'result-msg';

  Sound.flipStart();

  // отправляем запрос параллельно с началом анимации
  const playPromise = API.play({
    userId: state.userId,
    bet: bet,
    side: state.side,
    clientSeed: state.clientSeed,
    nonce: state.nonce++
  });

  try {
    const res = await playPromise;
    // ждём завершение анимации на конкретную сторону
    await Animation.flip(els.coin, res.Result);

    state.balance = parseFloat(res.NewBalance);
    updateBalance();

    els.coef.textContent = 'x' + res.Coefficient.toFixed(2);

    if (res.Won) {
      state.series++;
      flashMsg(`+ $${res.WinAmount.toFixed(2)}`, 'win');
      Sound.win();
      Animation.pulse(els.balance);
    } else {
      state.series = 0;
      flashMsg(`Проигрыш $${bet.toFixed(2)}`, 'lose');
      Sound.lose();
    }
    els.series.textContent = state.series;

    pushHistory(res.Won);
  } catch (e) {
    flashMsg('Ошибка: ' + e.message, 'lose');
  } finally {
    state.busy = false;
    els.playBtn.disabled = false;
  }
}

function flashMsg(text, cls) {
  els.result.textContent = text;
  els.result.className = 'result-msg ' + (cls || '');
}

// ---- history ----
async function loadHistory() {
  try {
    const rows = await API.history(state.userId, 20);
    state.history = rows.map(r => r.win_amount > 0);
    renderHistory();
  } catch (e) {
    console.warn('history fail', e);
  }
}

function pushHistory(won) {
  state.history.unshift(won);
  if (state.history.length > 20) state.history.pop();
  renderHistory();
}

function renderHistory() {
  els.histList.innerHTML = state.history
    .map(w => `<span class="hist-item ${w ? 'hist-win' : 'hist-lose'}">${w ? 'W' : 'L'}</span>`)
    .join('');
}

init();