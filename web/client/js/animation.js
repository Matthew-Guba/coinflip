const Animation = {
  // Крутим монетку и возвращаем Promise, который резолвится когда анимация завершена.
  // finalSide: 'heads' | 'tails'
  flip(coinEl, finalSide) {
    return new Promise((resolve) => {
      // Итоговый угол: heads = кратно 360, tails = 360*N + 180
      const rot = finalSide === 'heads' ? 1800 : 1980;
      coinEl.style.setProperty('--final-rot', rot + 'deg');
      coinEl.classList.remove('flipping');
      // рефлоу
      void coinEl.offsetWidth;
      coinEl.classList.add('flipping');

      // тики звука во время крутки
      let ticks = 0;
      const interval = setInterval(() => {
        Sound.flipTick();
        if (++ticks > 15) clearInterval(interval);
      }, 100);

      const onEnd = () => {
        coinEl.removeEventListener('animationend', onEnd);
        clearInterval(interval);
        resolve();
      };
      coinEl.addEventListener('animationend', onEnd);
    });
  },

  pulse(el) {
    el.animate([
      { transform: 'scale(1)' },
      { transform: 'scale(1.2)' },
      { transform: 'scale(1)' }
    ], { duration: 300, easing: 'ease-out' });
  }
};