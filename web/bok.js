// Vandrarens bok (docs/vandrarens-bok.md): en hjälpsida i ett lager över
// spelarsidan, med flikar. Grundernas texter är sidans egna (bok_* i
// data/spelarsidan.json); klasserna kommer från tornet (GET /bok, bok.go)
// på tornets språk. Boken öppnas med B, knappen i styrningen och knappen i
// introt, och stängs med B, Esc, ✕ eller ett klick utanför.
"use strict";

window.Bok = (() => {
  const $ = (id) => document.getElementById(id);
  const FLIKAR = ["grunder", "klasser"];
  let flik = "grunder";
  let innehall = null; // svaret från /bok
  let hamtar = null;

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

  function rita() {
    ritaFlikar();
    const sida = $("bokSida");
    sida.innerHTML = "";
    sida.dataset.flik = flik;
    if (flik === "klasser") ritaKlasser(sida);
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

  return { open, close, toggle, oppen };
})();
