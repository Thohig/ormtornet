// Vandrarens bok (docs/vandrarens-bok.md): en hjälpsida i ett lager över
// spelarsidan, med flikar. Grundernas texter är sidans egna (bok_* i
// data/spelarsidan.json). Resten kommer från tornet (GET /bok, bok.go) på
// tornets språk: klasserna, och föremålen, kombinationerna, varelserna,
// fällorna, händelserna och ormens kort, alla öppna från början. Det man
// har mött får märket Mött: tornet säger det i tillståndet (bok), och
// webbläsaren sparar det i localStorage, från vandring till vandring. Boken
// öppnas med B, knappen i styrningen och knappen i introt, och stängs med
// B, Esc, ✕ eller ett klick utanför.
"use strict";

window.Bok = (() => {
  const $ = (id) => document.getElementById(id);
  const FLIKAR = ["grunder", "klasser", "foremal", "kombinationer", "varelser", "fallor", "handelser", "ormen"];
  // Flikarna med kort, med slaget i tornets nycklar ("varelse:krypare").
  const SLAG = {
    foremal: "foremal", kombinationer: "kombination", varelser: "varelse",
    fallor: "falla", handelser: "handelse", ormen: "ormkort",
  };
  const NYCKEL = "ormtornet-bok";
  let flik = "grunder";
  let innehall = null; // svaret från /bok
  let hamtar = null;

  // mott är allt man har mött, i den här webbläsaren.
  const mott = new Set();
  try {
    for (const k of JSON.parse(localStorage.getItem(NYCKEL) || "[]")) mott.add(k);
  } catch {}

  // merge lägger till det tornet säger att man har mött.
  function merge(lista) {
    if (!Array.isArray(lista)) return;
    let nytt = false;
    for (const k of lista) {
      if (!mott.has(k)) {
        mott.add(k);
        nytt = true;
      }
    }
    if (!nytt) return;
    try {
      localStorage.setItem(NYCKEL, JSON.stringify([...mott]));
    } catch {}
    if (oppen()) rita();
  }

  function el(tag, cls, text) {
    const e = document.createElement(tag);
    if (cls) e.className = cls;
    if (text !== undefined) e.textContent = text;
    return e;
  }

  function oppen() {
    return !$("bok").hidden;
  }

  // hamta frågar tornet efter bokens innehåll, varje gång boken öppnas, så
  // att den följer tornets språk och innehåll.
  function hamta() {
    if (hamtar) return hamtar;
    hamtar = fetch("/bok", { cache: "no-store" })
      .then((r) => (r.ok ? r.json() : Promise.reject(r.status)))
      .then((v) => {
        innehall = v;
      })
      .catch(() => {
        innehall = innehall || { fel: true };
      })
      .finally(() => {
        hamtar = null;
        if (oppen()) rita();
      });
    return hamtar;
  }

  function ritaFlikar() {
    const rad = $("bokFlikar");
    rad.innerHTML = "";
    for (const id of FLIKAR) {
      const b = el("button", "bok-flik", t("bok_flik_" + id));
      b.type = "button";
      b.setAttribute("role", "tab");
      b.setAttribute("aria-selected", id === flik ? "true" : "false");
      b.addEventListener("click", () => {
        flik = id;
        rita();
      });
      rad.appendChild(b);
    }
  }

  function ritaGrunder(sida) {
    for (const id of t("bok_grunder").split(/\s+/).filter(Boolean)) {
      const avsnitt = el("section", "bok-avsnitt");
      avsnitt.append(el("h3", "", t("bok_" + id + "_rubrik")), el("p", "", t("bok_" + id)));
      sida.appendChild(avsnitt);
    }
  }

  // klasskort bygger ett klasskort som i klassvalet, men med framsidan upp
  // och utan att gå att välja.
  function klasskort(c) {
    const kort = el("div", "card vand bok-klass");
    const inner = el("div", "kk-inner");
    const front = el("div", "kk-front");
    front.append(el("div", "cname", (c.tecken ? c.tecken + " " : "") + c.namn));
    const bild = el("div", "cbild");
    bild.innerHTML = Kort.bild(c.bild);
    front.append(bild, el("div", "desc", c.beskrivning));
    const rader = [];
    if (c.liv) rader.push(t(c.liv > 0 ? "klasskort_liv_mer" : "klasskort_liv_mindre", { n: Math.abs(c.liv) }));
    if (c.drycker) rader.push(t(c.drycker === 1 ? "klasskort_dryck_en" : "klasskort_dryck_flera", { n: c.drycker }));
    rader.push(t("klasskort_strid", { kraft: c.kraft || 0, forsvar: c.forsvar || 0 }));
    front.append(el("div", "stats", rader.join(" · ")));
    front.append(el("div", "good", t("klasskort_styrka", { styrka: c.styrka })));
    front.append(el("div", "bad", t("klasskort_svaghet", { svaghet: c.svaghet })));
    inner.appendChild(front);
    kort.appendChild(inner);
    return kort;
  }

  function ritaKlasser(sida) {
    if (!innehall) {
      sida.appendChild(el("p", "bok-inledning", t("bok_laddar")));
      return;
    }
    if (innehall.fel) {
      sida.appendChild(el("p", "bok-inledning", t("bok_fel")));
      return;
    }
    sida.appendChild(el("p", "bok-inledning", t("bok_klasser_inledning")));
    const lek = el("div", "bok-lek");
    for (const c of innehall.klasser || []) lek.appendChild(klasskort(c));
    sida.appendChild(lek);
  }

  // kort bygger ett spelkort i boken, med framsidan upp: rubriken, bilden,
  // texten och raderna under, och märket Mött om man har mött det.
  function kort(c, rader, harMott) {
    const k = el("div", "okort vand bok-kort");
    const inner = el("div", "ok-inner");
    const front = el("div", "ok-front");
    front.appendChild(el("div", "ok-namn", c.rubrik));
    if (harMott) front.appendChild(el("div", "bok-mott", t("bok_mott")));
    const bild = el("div", "ok-bild");
    bild.innerHTML = Kort.bild(c.bild);
    front.appendChild(bild);
    if (c.text) front.appendChild(el("div", "ok-text", c.text));
    for (const [cls, text] of rader) if (text) front.appendChild(el("div", "bok-rad " + cls, text));
    inner.appendChild(front);
    k.appendChild(inner);
    return k;
  }

  // ritaLek ritar en flik med kort: föremålen, kombinationerna, varelserna,
  // fällorna, händelserna eller ormens kort.
  function ritaLek(sida) {
    if (!innehall || innehall.fel) {
      sida.appendChild(el("p", "bok-inledning", t(innehall ? "bok_fel" : "bok_laddar")));
      return;
    }
    const slag = SLAG[flik];
    sida.appendChild(el("p", "bok-inledning", t("bok_" + flik + "_inledning")));
    const lek = el("div", "bok-lek");
    for (const c of innehall[flik] || []) {
      let rader = [];
      if (flik === "foremal") {
        if (c.forbannelse) rader.push(["bad", t("bok_forbannelse", { text: c.forbannelse })]);
        if (c.vander) rader.push(["bad", t("bok_vander", { text: c.vander })]);
      } else if (flik === "kombinationer") {
        rader = [["om", c.delar && t("bok_delar", { delar: c.delar.join(" + ") })]];
      } else if (flik !== "fallor") {
        rader = [["om", c.om]];
      } else {
        rader = [["om", c.om], ["spar", c.spar && t("bok_spar", { text: c.spar })]];
      }
      lek.appendChild(kort(c, rader, mott.has(slag + ":" + c.id)));
    }
    sida.appendChild(lek);
  }

  function rita() {
    ritaFlikar();
    const sida = $("bokSida");
    sida.innerHTML = "";
    sida.dataset.flik = flik;
    if (flik === "klasser") ritaKlasser(sida);
    else if (SLAG[flik]) ritaLek(sida);
    else ritaGrunder(sida);
  }

  function open() {
    const bok = $("bok");
    if (!oppen()) {
      bok.hidden = false;
      $("bokSida").scrollTop = 0;
    }
    rita();
    hamta();
    $("bokSida").focus({ preventScroll: true });
  }

  function close() {
    if (!oppen()) return false;
    $("bok").hidden = true;
    if (document.activeElement && $("bok").contains(document.activeElement)) document.activeElement.blur();
    return true;
  }

  function toggle() {
    if (oppen()) close();
    else open();
  }

  $("bokStang").addEventListener("click", close);
  // Ett klick utanför sidan stänger boken. Flikarna ritas om när man
  // klickar på dem, så det är bara klick på själva lagret som räknas.
  $("bok").addEventListener("click", (e) => {
    if (e.target === $("bok")) close();
  });
  // Boken tar tangenterna medan den är öppen, så att de inte blir drag; B
  // öppnar den var man än är, utom när man skriver. Skriptet läses före
  // app.js, så att den här lyssnaren kommer först.
  document.addEventListener("keydown", (e) => {
    if (e.metaKey || e.ctrlKey || e.altKey) return;
    const skriver = e.target && (e.target.tagName === "INPUT" || e.target.tagName === "TEXTAREA");
    if (oppen()) {
      if (e.key === "Escape" || e.key === "b" || e.key === "B") {
        e.preventDefault();
        close();
      } else if (e.key === "ArrowLeft" || e.key === "ArrowRight") {
        e.preventDefault();
        const i = FLIKAR.indexOf(flik) + (e.key === "ArrowRight" ? 1 : -1);
        flik = FLIKAR[(i + FLIKAR.length) % FLIKAR.length];
        rita();
      }
      e.stopImmediatePropagation();
      return;
    }
    if (!skriver && (e.key === "b" || e.key === "B")) {
      e.preventDefault();
      e.stopImmediatePropagation();
      open();
    }
  }, true);

  return { open, close, toggle, oppen, merge };
})();
