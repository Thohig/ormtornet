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

  // Ett eko-rum som steg och andra ljud kan skickas genom.
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

  // Ormens andning i kammaren. Ju närmare att vakna, desto tätare och tyngre.
  let breathNode = null;
  let breathLevel = -1;
  function setBreath(level) {
    if (!ctx || level === breathLevel) return;
    breathLevel = level;
    const t = ctx.currentTime;
    if (breathNode) {
      const old = breathNode;
      old.out.gain.setTargetAtTime(0, t, 0.4);
      setTimeout(() => { old.src.stop(); old.lfo.stop(); }, 2500);
      breathNode = null;
    }
    if (level < 0) return;
    const rate = [0.12, 0.2, 0.33][Math.min(level, 2)];
    const depth = [0.10, 0.16, 0.24][Math.min(level, 2)];
    const src = ctx.createBufferSource();
    src.buffer = brown;
    src.loop = true;
    const lp = ctx.createBiquadFilter();
    lp.type = "lowpass";
    lp.frequency.value = 260 + level * 80;
    const amp = ctx.createGain();
    amp.gain.value = depth;
    const lfo = ctx.createOscillator();
    lfo.frequency.value = rate;
    const lfoGain = ctx.createGain();
    lfoGain.gain.value = depth;
    lfo.connect(lfoGain);
    lfoGain.connect(amp.gain);
    lfoGain.connect(lp.frequency);
    const out = ctx.createGain();
    out.gain.value = 0;
    out.gain.setTargetAtTime(1, t, 0.8);
    src.connect(lp);
    lp.connect(amp);
    amp.connect(out);
    out.connect(master);
    src.start();
    lfo.start();
    breathNode = { src, lfo, out };
  }

  // Dunket. Nivå 1: ett svagt, långsamt dunk. Nivå 2: ett hjärta.
  // Nivå 3 och 4: snabbare och starkare, med kluckande ljud emellan.
  // Nivå 0 betyder att avslöjandet är avstängt: inget dunk alls.
  let heartLevel = 0;
  let heartTimer = null;
  const HEART = {
    1: { gap: 3600, gain: 0.05, double: false, gurgle: 0 },
    2: { gap: 1700, gain: 0.10, double: true, gurgle: 0 },
    3: { gap: 1250, gain: 0.15, double: true, gurgle: 0.25 },
    4: { gap: 950, gain: 0.22, double: true, gurgle: 0.45 },
  };

  function thump(t, gain) {
    const o = ctx.createOscillator();
    o.type = "sine";
    o.frequency.setValueAtTime(58, t);
    o.frequency.exponentialRampToValueAtTime(34, t + 0.18);
    const g = ctx.createGain();
    g.gain.setValueAtTime(0.0001, t);
    g.gain.exponentialRampToValueAtTime(gain, t + 0.015);
    g.gain.exponentialRampToValueAtTime(0.0001, t + 0.3);
    const lp = ctx.createBiquadFilter();
    lp.type = "lowpass";
    lp.frequency.value = 160;
    o.connect(lp);
    lp.connect(g);
    g.connect(master);
    o.start(t);
    o.stop(t + 0.35);
  }

  function gurgle() {
    if (!ctx) return;
    const t = ctx.currentTime;
    const src = ctx.createBufferSource();
    src.buffer = white;
    const bp = ctx.createBiquadFilter();
    bp.type = "bandpass";
    bp.Q.value = 9;
    for (let i = 0; i < 6; i++) {
      bp.frequency.setValueAtTime(200 + Math.random() * 500, t + i * 0.07);
    }
    const g = ctx.createGain();
    env(g, t, 0.07, 0.05, 0.45);
    src.connect(bp);
    bp.connect(g);
    g.connect(master);
    src.start(t);
    src.stop(t + 0.55);
  }

  function heartbeat() {
    heartTimer = null;
    const h = HEART[heartLevel];
    if (!ctx || !h) return;
    const t = ctx.currentTime + 0.02;
    thump(t, h.gain);
    if (h.double) thump(t + 0.22, h.gain * 0.7);
    if (Math.random() < h.gurgle) setTimeout(gurgle, 400 + Math.random() * 300);
    heartTimer = setTimeout(heartbeat, h.gap * (0.95 + Math.random() * 0.1));
  }

  function setHeart(level) {
    if (level === heartLevel) return;
    heartLevel = level;
    if (heartTimer) clearTimeout(heartTimer);
    heartTimer = null;
    if (ctx && HEART[level]) heartTimer = setTimeout(heartbeat, 500);
  }

  // Introt: en vindstöt, ett dubbelt slag långt bort och sedan ett långt
  // andetag in, som om något drar in luft genom porten.
  function intro() {
    const t = ctx.currentTime;
    const gust = noiseHit(brown, "bandpass", 260, 0.3, 0.9, 1.8, false);
    gust.frequency.linearRampToValueAtTime(650, t + 1.1);
    thump(t + 0.75, 0.32);
    thump(t + 0.97, 0.22);
    const src = ctx.createBufferSource();
    src.buffer = brown;
    src.loop = true;
    const lp = ctx.createBiquadFilter();
    lp.type = "lowpass";
    lp.frequency.setValueAtTime(120, t + 1.3);
    lp.frequency.exponentialRampToValueAtTime(900, t + 2.9);
    const g = ctx.createGain();
    g.gain.setValueAtTime(0.0001, t + 1.3);
    g.gain.exponentialRampToValueAtTime(0.32, t + 2.7);
    g.gain.exponentialRampToValueAtTime(0.0001, t + 3.2);
    src.connect(lp);
    lp.connect(g);
    g.connect(master);
    g.connect(echo);
    src.start(t + 1.3);
    src.stop(t + 3.3);
  }

  // Stormen på startsidan: en låg drönarton i moll, vind, regn och ibland
  // ett dunk långt inne i tornet. Allt går genom en egen volym så att
  // stormen kan tona ut när man går in.
  let storm = null;

  function startStorm(reveal) {
    if (storm) return;
    const t = ctx.currentTime;
    const bus = ctx.createGain();
    bus.gain.setValueAtTime(0.0001, t);
    bus.gain.exponentialRampToValueAtTime(1, t + 4);
    bus.connect(master);
    const nodes = [];
    const loop = (buf) => {
      const src = ctx.createBufferSource();
      src.buffer = buf;
      src.loop = true;
      src.start();
      nodes.push(src);
      return src;
    };

    // Drönaren: A och C djupt nere, svävande, med en sakta pulserande volym.
    const drone = ctx.createGain();
    drone.gain.value = 0.05;
    const dlp = ctx.createBiquadFilter();
    dlp.type = "lowpass";
    dlp.frequency.value = 420;
    drone.connect(dlp);
    dlp.connect(bus);
    dlp.connect(echo);
    [[55, "sawtooth", 0.5], [55.4, "sine", 1], [65.4, "triangle", 0.7], [110, "sine", 0.25]].forEach(([f, type, level]) => {
      const o = ctx.createOscillator();
      o.type = type;
      o.frequency.value = f;
      const g = ctx.createGain();
      g.gain.value = level;
      o.connect(g);
      g.connect(drone);
      o.start();
      nodes.push(o);
    });
    lfo(0.07, 0.025, drone.gain);
    lfo(0.03, 180, dlp.frequency);

    // Tritonus: Ess som sväller in och tonar bort med ojämna mellanrum.
    const tritone = ctx.createOscillator();
    tritone.type = "triangle";
    tritone.frequency.value = 77.8;
    const tg = ctx.createGain();
    tg.gain.value = 0.0001;
    tritone.connect(tg);
    tg.connect(dlp);
    tritone.start();
    nodes.push(tritone);

    // Vinden: starkare och mer ojämn än i tornet.
    const wsrc = loop(brown);
    const wbp = ctx.createBiquadFilter();
    wbp.type = "bandpass";
    wbp.frequency.value = 500;
    wbp.Q.value = 0.9;
    const wg = ctx.createGain();
    wg.gain.value = 0.16;
    lfo(0.09, 380, wbp.frequency);
    lfo(0.17, 0.08, wg.gain);
    wsrc.connect(wbp);
    wbp.connect(wg);
    wg.connect(bus);

    // Regnet: vitt brus, ljust men inte vasst.
    const rsrc = loop(white);
    const rhp = ctx.createBiquadFilter();
    rhp.type = "highpass";
    rhp.frequency.value = 1200;
    const rlp = ctx.createBiquadFilter();
    rlp.type = "lowpass";
    rlp.frequency.value = 7000;
    const rg = ctx.createGain();
    rg.gain.value = 0.035;
    lfo(0.13, 0.012, rg.gain);
    rsrc.connect(rhp);
    rhp.connect(rlp);
    rlp.connect(rg);
    rg.connect(bus);

    storm = { bus, nodes, timers: [] };
    const later = (fn, ms) => storm && storm.timers.push(setTimeout(fn, ms));
    const swell = () => {
      if (!storm || !ctx) return;
      const now = ctx.currentTime;
      tg.gain.cancelScheduledValues(now);
      tg.gain.setValueAtTime(0.0001, now);
      tg.gain.exponentialRampToValueAtTime(0.35, now + 3.5);
      tg.gain.exponentialRampToValueAtTime(0.0001, now + 8);
      later(swell, 14000 + Math.random() * 12000);
    };
    later(swell, 6000);
    if (reveal) {
      // Ett dubbelt slag långt inne i tornet, som ett hjärta som sover.
      const beat = () => {
        if (!storm || !ctx) return;
        const now = ctx.currentTime;
        thump(now, 0.16);
        thump(now + 0.24, 0.1);
        later(beat, 15000 + Math.random() * 10000);
      };
      later(beat, 9000);
    }
  }

  function stopStorm(seconds) {
    if (!storm || !ctx) return;
    const s = storm;
    storm = null;
    s.timers.forEach(clearTimeout);
    const t = ctx.currentTime;
    s.bus.gain.cancelScheduledValues(t);
    s.bus.gain.setValueAtTime(Math.max(0.0001, s.bus.gain.value), t);
    s.bus.gain.exponentialRampToValueAtTime(0.0001, t + seconds);
    s.nodes.forEach((n) => n.stop(t + seconds + 0.1));
  }

  // Åska: ett skarpt brak när den är nära, och ett långt muller. near går
  // från 0 (långt bort) till 1 (precis här).
  function thunder(near) {
    if (!ctx) return;
    const t = ctx.currentTime;
    if (near > 0.6) {
      noiseHit(white, "highpass", 900, 0.12 * near, 0.005, 0.5, true);
    }
    const r = noiseHit(brown, "lowpass", 160 + 260 * near, 0.5 + 0.4 * near, 0.08 + 0.4 * (1 - near), 3.5 + 2 * (1 - near), true);
    r.frequency.setValueAtTime(r.frequency.value, t);
    r.frequency.linearRampToValueAtTime(90, t + 4);
    tone(38, 30, "sine", 0.18 * near + 0.05, 0.1, 3, false);
  }

  const effects = {
    intro,
    // Tornets röst: en viskning med en låg ton under.
    rost() {
      noiseHit(white, "bandpass", 1900, 0.05, 0.25, 1.1, true);
      tone(196, 147, "sine", 0.04, 0.4, 1.6, true);
    },
    // Tornet går över till en ny nivå: en djup, svällande ton.
    niva() {
      tone(41, 55, "sine", 0.3, 1.5, 2.5, false);
      tone(82, 110, "triangle", 0.05, 1.5, 2.5, true);
      const f = noiseHit(brown, "lowpass", 200, 0.3, 1.5, 2.0, false);
      f.frequency.linearRampToValueAtTime(500, ctx.currentTime + 1.5);
    },
    // Ormen vaknar: ett väsande som fyller hela tornet, och ett lågt muller.
    vasning() {
      const t = ctx.currentTime;
      const f = noiseHit(white, "highpass", 1800, 0.45, 0.25, 3.2, true);
      f.frequency.setValueAtTime(1800, t);
      f.frequency.linearRampToValueAtTime(5200, t + 1.2);
      f.frequency.linearRampToValueAtTime(2400, t + 3.4);
      const b = noiseHit(white, "bandpass", 3500, 0.25, 0.3, 2.8, false);
      b.Q.value = 2;
      tone(46, 32, "sawtooth", 0.12, 0.3, 3.0, false);
      tone(36, 28, "sine", 0.35, 0.2, 3.2, false);
    },
    // Ormen rör sig i sömnen: golvet gungar.
    ormen_ror() {
      const f = noiseHit(brown, "lowpass", 160, 0.35, 0.6, 1.8, false);
      f.frequency.linearRampToValueAtTime(90, ctx.currentTime + 2.4);
      tone(41, 36, "sine", 0.25, 0.6, 1.8, false);
    },
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
    // Ett mjukt klick när man väljer drag.
    val() {
      tone(880, 660, "triangle", 0.05, 0.003, 0.07, false);
    },
    // Ett lågt slag när rundan utförs.
    runda() {
      tone(110, 82, "sine", 0.18, 0.01, 0.5, true);
      noiseHit(white, "lowpass", 400, 0.05, 0.005, 0.12, false);
    },
    tick() {
      tone(1500, null, "square", 0.025, 0.002, 0.03, false);
    },
    start() {
      [147, 220, 294].forEach((f, i) => tone(f, null, "triangle", 0.06, 0.3, 2.2, true, i * 0.35));
    },
    // Som en klocka långt borta.
    solvarning() {
      [523, 1046, 1569].forEach((f, i) => tone(f, f * 0.995, "sine", 0.06 / (i + 1), 0.005, 2.5, true));
    },
    solnedgang() {
      tone(220, 55, "sawtooth", 0.06, 0.5, 3.5, true);
      tone(110, 40, "sine", 0.2, 0.3, 3.5, false);
      const f = noiseHit(brown, "lowpass", 600, 0.3, 1.0, 2.5, false);
      f.frequency.exponentialRampToValueAtTime(80, ctx.currentTime + 3.5);
    },
    ut() {
      const f = noiseHit(white, "bandpass", 500, 0.08, 0.6, 1.4, false);
      f.frequency.exponentialRampToValueAtTime(2500, ctx.currentTime + 2);
      tone(392, null, "triangle", 0.04, 0.4, 1.5, true, 0.3);
    },
    skatt() {
      for (let i = 0; i < 5; i++) {
        const f = 2200 + Math.random() * 1800;
        tone(f, f * 0.98, "triangle", 0.04, 0.002, 0.18, true, i * (0.05 + Math.random() * 0.06));
      }
    },
    slut() {
      [294, 247, 196, 147].forEach((f, i) => tone(f, null, "triangle", 0.05, 0.1, 1.4, true, i * 0.45));
    },
    // Fällor. Varje fälla har sitt eget ljud.
    fallgaller() {
      // Järn som slår i sten: flera ojämna metalltoner och en smäll.
      [347, 523, 811, 1213].forEach((f) => tone(f, f * 0.97, "square", 0.03, 0.002, 0.9, true));
      noiseHit(white, "highpass", 2500, 0.2, 0.002, 0.25, true);
      tone(70, 40, "sine", 0.35, 0.005, 0.4, false);
    },
    rasande_golv() {
      const f = noiseHit(brown, "lowpass", 900, 0.5, 0.05, 1.6, true);
      f.frequency.exponentialRampToValueAtTime(90, ctx.currentTime + 1.6);
      tone(80, 30, "sine", 0.35, 0.05, 1.4, false);
      for (let i = 0; i < 6; i++) tone(200 + Math.random() * 300, 80, "triangle", 0.05, 0.003, 0.15, false, 0.2 + i * 0.12);
    },
    roterande() {
      // Malande sten som vrider sig.
      const f = noiseHit(white, "bandpass", 180, 0.35, 0.3, 1.6, true);
      f.Q.value = 6;
      f.frequency.exponentialRampToValueAtTime(420, ctx.currentTime + 1.0);
      f.frequency.exponentialRampToValueAtTime(150, ctx.currentTime + 1.9);
      tone(55, 50, "sawtooth", 0.05, 0.3, 1.5, false);
    },
    pilar() {
      for (let i = 0; i < 3; i++) {
        const d = i * 0.13;
        setTimeout(() => {
          const f = noiseHit(white, "bandpass", 3000, 0.12, 0.01, 0.12, false);
          f.frequency.exponentialRampToValueAtTime(900, ctx.currentTime + 0.13);
          tone(160, 90, "triangle", 0.15, 0.002, 0.08, true, 0.12);
        }, d * 1000);
      }
    },
    dimma() {
      // Ett långt, fuktigt väsande.
      const f = noiseHit(white, "highpass", 3500, 0.07, 0.8, 2.5, true);
      f.frequency.linearRampToValueAtTime(2000, ctx.currentTime + 3);
      tone(98, 92, "sine", 0.05, 1.0, 2.0, false);
    },
    ormgrop() {
      for (let i = 0; i < 8; i++) {
        setTimeout(() => {
          const f = noiseHit(white, "highpass", 4000 + Math.random() * 2500, 0.05, 0.05, 0.35, false);
          f.Q.value = 3;
        }, i * 110 + Math.random() * 80);
      }
    },
    skinn() {
      // Torrt prassel.
      for (let i = 0; i < 14; i++) {
        setTimeout(() => noiseHit(white, "bandpass", 2500 + Math.random() * 3000, 0.06, 0.002, 0.03, false), i * 45 + Math.random() * 40);
      }
    },
    klocka() {
      // En djup klocka med långt efterljud.
      [110, 220.5, 267, 331, 440.8].forEach((f, i) => tone(f, f * 0.998, "sine", 0.12 / (i + 1), 0.004, 5, true));
    },
    morker() {
      // En pust som blåser ut lågan.
      const f = noiseHit(white, "lowpass", 1400, 0.2, 0.05, 0.5, false);
      f.frequency.exponentialRampToValueAtTime(200, ctx.currentTime + 0.5);
      tone(140, 60, "sine", 0.08, 0.05, 0.8, false);
    },
    eko() {
      for (let i = 0; i < 6; i++) tone(660, 640, "sine", 0.1 * Math.pow(0.7, i), 0.003, 0.3, false, i * 0.33);
    },
    skada() {
      tone(120, 60, "square", 0.08, 0.003, 0.18, false);
      noiseHit(white, "lowpass", 1200, 0.15, 0.002, 0.1, false);
    },
    dod() {
      tone(196, 49, "sawtooth", 0.06, 0.05, 2.5, true);
      tone(98, 40, "sine", 0.2, 0.05, 2.5, false);
    },
    hitta() {
      tone(660, null, "triangle", 0.05, 0.005, 0.3, true);
      tone(990, null, "triangle", 0.04, 0.005, 0.4, true, 0.12);
    },
    gang() {
      const f = noiseHit(brown, "lowpass", 300, 0.4, 0.2, 1.2, true);
      f.frequency.linearRampToValueAtTime(150, ctx.currentTime + 1.4);
    },
    inget() {
      noiseHit(white, "bandpass", 1200, 0.04, 0.02, 0.2, false);
    },
    drick() {
      // Klunkar.
      for (let i = 0; i < 3; i++) tone(260, 120, "sine", 0.12, 0.01, 0.12, false, i * 0.28);
    },
    // Strider.
    monster() {
      // Ett lågt, rossligt morrande.
      const t = ctx.currentTime;
      const o = ctx.createOscillator();
      o.type = "sawtooth";
      o.frequency.setValueAtTime(70, t);
      o.frequency.linearRampToValueAtTime(95, t + 0.4);
      o.frequency.linearRampToValueAtTime(60, t + 1.1);
      const trem = ctx.createOscillator();
      trem.frequency.value = 23;
      const tremGain = ctx.createGain();
      tremGain.gain.value = 0.05;
      const g = ctx.createGain();
      env(g, t, 0.08, 0.15, 1.0);
      trem.connect(tremGain);
      tremGain.connect(g.gain);
      const lp = ctx.createBiquadFilter();
      lp.type = "lowpass";
      lp.frequency.value = 600;
      o.connect(lp);
      lp.connect(g);
      g.connect(master);
      g.connect(echo);
      o.start(t);
      trem.start(t);
      o.stop(t + 1.3);
      trem.stop(t + 1.3);
    },
    hugg() {
      const f = noiseHit(white, "bandpass", 1800, 0.25, 0.002, 0.12, false);
      f.frequency.exponentialRampToValueAtTime(400, ctx.currentTime + 0.12);
      tone(90, 50, "sine", 0.25, 0.002, 0.15, false);
    },
    parera() {
      [1320, 1980, 2710].forEach((f) => tone(f, f * 0.99, "triangle", 0.04, 0.001, 0.5, true));
    },
    undan() {
      const f = noiseHit(white, "bandpass", 600, 0.1, 0.05, 0.25, false);
      f.frequency.exponentialRampToValueAtTime(2200, ctx.currentTime + 0.3);
    },
    monster_dor() {
      tone(180, 40, "sawtooth", 0.07, 0.01, 1.2, true);
      noiseHit(brown, "lowpass", 400, 0.3, 0.01, 0.8, false);
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
    },
    // stop tonar ut allt ljud och stänger av ljudet för gott, till exempel
    // när tornet har stängts.
    stop() {
      if (!ctx) return;
      const t = ctx.currentTime;
      master.gain.cancelScheduledValues(t);
      master.gain.setValueAtTime(master.gain.value, t);
      master.gain.linearRampToValueAtTime(0, t + 2.5);
      const closing = ctx;
      ctx = null; // så att inga nya ljud startas
      if (heartTimer) clearTimeout(heartTimer);
      heartTimer = null;
      setTimeout(() => closing.close(), 2800);
    },
    // storm startar stormen på startsidan. reveal lägger till tornets dunk.
    storm(reveal) {
      if (!ctx) return;
      if (ctx.state === "suspended") ctx.resume();
      startStorm(reveal);
    },
    calm(seconds) {
      stopStorm(seconds || 3);
    },
    thunder(near) {
      thunder(near);
    },
    running() {
      return !!ctx && ctx.state === "running";
    },
    breath(level) {
      setBreath(level);
    },
    level(level) {
      setHeart(level);
    },
    play(name) {
      if (!ctx || !effects[name]) return;
      if (ctx.state === "suspended") ctx.resume();
      effects[name]();
    },
  };
})();
