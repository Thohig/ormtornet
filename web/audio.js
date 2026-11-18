// Allt ljud skapas här med Web Audio API. Inga ljudfiler.
"use strict";

const Sound = (() => {
  let ctx = null;
  let master = null;
  let echo = null;
  let brown = null;
  let white = null;

  function noiseBuffer(seconds, kind) {
    const n = Math.floor(ctx.sampleRate * seconds);
    const buf = ctx.createBuffer(1, n, ctx.sampleRate);
    const d = buf.getChannelData(0);
    let last = 0;
    for (let i = 0; i < n; i++) {
      const w = Math.random() * 2 - 1;
      if (kind === "brown") {
        last = (last + 0.02 * w) / 1.02;
        d[i] = last * 3.5;
      } else {
        d[i] = w;
      }
    }
    return buf;
  }

  // Ett eko-rum som droppar och steg kan skickas genom.
  function makeEcho() {
    const delay = ctx.createDelay(1.5);
    delay.delayTime.value = 0.29;
    const fb = ctx.createGain();
    fb.gain.value = 0.38;
    const lp = ctx.createBiquadFilter();
    lp.type = "lowpass";
    lp.frequency.value = 1600;
    delay.connect(lp);
    lp.connect(fb);
    fb.connect(delay);
    lp.connect(master);
    return delay;
  }

  function lfo(freq, depth, target) {
    const o = ctx.createOscillator();
    o.frequency.value = freq;
    const g = ctx.createGain();
    g.gain.value = depth;
    o.connect(g);
    g.connect(target);
    o.start();
  }

  // Vind: brunt brus genom ett bandpassfilter som sakta vandrar.
  function wind() {
    const src = ctx.createBufferSource();
    src.buffer = brown;
    src.loop = true;
    const bp = ctx.createBiquadFilter();
    bp.type = "bandpass";
    bp.frequency.value = 420;
    bp.Q.value = 0.8;
    const g = ctx.createGain();
    g.gain.value = 0.10;
    lfo(0.05, 260, bp.frequency);
    lfo(0.11, 0.05, g.gain);
    src.connect(bp);
    bp.connect(g);
    g.connect(master);
    src.start();
  }

  function drip() {
    const t = ctx.currentTime;
    const f0 = 1300 + Math.random() * 1000;
    const o = ctx.createOscillator();
    o.type = "sine";
    o.frequency.setValueAtTime(f0, t);
    o.frequency.exponentialRampToValueAtTime(f0 * 0.45, t + 0.09);
    const g = ctx.createGain();
    g.gain.setValueAtTime(0, t);
    g.gain.linearRampToValueAtTime(0.05 + Math.random() * 0.04, t + 0.004);
    g.gain.exponentialRampToValueAtTime(0.0001, t + 0.18);
    o.connect(g);
    g.connect(master);
    g.connect(echo);
    o.start(t);
    o.stop(t + 0.2);
    setTimeout(drip, 1500 + Math.random() * 5500);
  }

  function env(g, t, peak, attack, release) {
    g.gain.setValueAtTime(0.0001, t);
    g.gain.exponentialRampToValueAtTime(peak, t + attack);
    g.gain.exponentialRampToValueAtTime(0.0001, t + attack + release);
  }

  function noiseHit(buf, filterType, freq, peak, attack, release, toEcho) {
    const t = ctx.currentTime;
    const src = ctx.createBufferSource();
    src.buffer = buf;
    src.playbackRate.value = 0.8 + Math.random() * 0.4;
    const f = ctx.createBiquadFilter();
    f.type = filterType;
    f.frequency.value = freq;
    const g = ctx.createGain();
    env(g, t, peak, attack, release);
    src.connect(f);
    f.connect(g);
    g.connect(master);
    if (toEcho) g.connect(echo);
    src.start(t);
    src.stop(t + attack + release + 0.05);
    return f;
  }

  function tone(freq, endFreq, type, peak, attack, release, toEcho, delay) {
    const t = ctx.currentTime + (delay || 0);
    const o = ctx.createOscillator();
    o.type = type;
    o.frequency.setValueAtTime(freq, t);
    if (endFreq) o.frequency.exponentialRampToValueAtTime(endFreq, t + attack + release);
    const g = ctx.createGain();
    env(g, t, peak, attack, release);
    o.connect(g);
    g.connect(master);
    if (toEcho) g.connect(echo);
    o.start(t);
    o.stop(t + attack + release + 0.05);
  }

  const effects = {
    // Sten som skrapar mot sten, och en dov duns.
    nyttrum() {
      const f = noiseHit(white, "lowpass", 900, 0.22, 0.12, 0.9, true);
      f.frequency.exponentialRampToValueAtTime(180, ctx.currentTime + 0.9);
      tone(62, 38, "sine", 0.35, 0.02, 0.6, false);
    },
    steg() {
      noiseHit(white, "bandpass", 700, 0.06, 0.005, 0.06, true);
    },
    vagg() {
      tone(95, 50, "sine", 0.25, 0.005, 0.15, false);
    },
    ansluter() {
      tone(196, null, "triangle", 0.05, 0.4, 1.6, true);
      tone(294, null, "triangle", 0.035, 0.4, 1.6, true, 0.25);
    },
  };

  return {
    start() {
      if (ctx) {
        ctx.resume();
        return;
      }
      const AC = window.AudioContext || window.webkitAudioContext;
      if (!AC) return;
      ctx = new AC();
      master = ctx.createGain();
      master.gain.value = 0.8;
      master.connect(ctx.destination);
      echo = makeEcho();
      brown = noiseBuffer(4, "brown");
      white = noiseBuffer(2, "white");
      wind();
      setTimeout(drip, 800);
    },
    play(name) {
      if (!ctx || !effects[name]) return;
      if (ctx.state === "suspended") ctx.resume();
      effects[name]();
    },
  };
})();
