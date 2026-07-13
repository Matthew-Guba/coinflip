const API = {
  baseURL: '', // тот же домен

  async play({userId, bet, side, clientSeed, nonce}) {
    const r = await fetch(this.baseURL + '/api/play', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({
        user_id: userId,
        bet: bet,
        side: side,
        client_seed: clientSeed,
        nonce: nonce
      })
    });
    if (!r.ok) throw new Error(await r.text());
    return r.json();
  },

  async history(userId, limit = 20) {
    const r = await fetch(`${this.baseURL}/api/history?user_id=${userId}&limit=${limit}`);
    if (!r.ok) throw new Error(await r.text());
    return r.json();
  },

  async balance(userId) {
    const r = await fetch(`${this.baseURL}/api/balance?user_id=${userId}`);
    if (!r.ok) throw new Error(await r.text());
    return r.json();
  }
};