class SoundEngine {
  constructor() {
    this.enabled = true;
    this.ctx = null;
  }

  _ensureCtx() {
    if (!this.ctx) {
      this.ctx = new (window.AudioContext || window.webkitAudioContext)();
    }
  }

  toggle() { this.enabled = !this.enabled; return this.enabled; }

  _beep({freq = 440, duration = 0.15, type = 'sine', volume = 0.15, slideTo = null}) {
    if (!this.enabled) return;
    this._ensureCtx();
    const t0 = this.ctx.currentTime;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    osc.type = type;
    osc.frequency.setValueAtTime(freq, t0);
    if (slideTo) osc.frequency.exponentialRampToValueAtTime(slideTo, t0 + duration);
    gain.gain.setValueAtTime(volume, t0);
    gain.gain.exponentialRampToValueAtTime(0.001, t0 + duration);
    osc.connect(gain).connect(this.ctx.destination);
    osc.start(t0);
    osc.stop(t0 + duration);
  }

  flipStart() {
    this._beep({freq: 300, slideTo: 900, duration: 0.4, type: 'triangle', volume: 0.1});
  }

  flipTick() {
    this._beep({freq: 800, duration: 0.05, type: 'square', volume: 0.05});
  }

  win() {
    this._beep({freq: 523, duration: 0.12, type: 'sine', volume: 0.2});
    setTimeout(() => this._beep({freq: 659, duration: 0.12, type: 'sine', volume: 0.2}), 120);
    setTimeout(() => this._beep({freq: 784, duration: 0.2,  type: 'sine', volume: 0.2}), 240);
  }

  lose() {
    this._beep({freq: 300, slideTo: 120, duration: 0.5, type: 'sawtooth', volume: 0.15});
  }

  click() {
    this._beep({freq: 600, duration: 0.05, type: 'square', volume: 0.08});
  }
}

const Sound = new SoundEngine();