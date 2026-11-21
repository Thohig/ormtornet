// Sidornas texter. De står i data/spelarsidan.json och data/vardsidan.json
// och kommer på tornets språk från /texter/, som skript innan sidans eget
// skript (window.UI). Byts språket hämtar sidan nya texter med
// refreshTexts och ritar om det som syns.
"use strict";

window.UI = window.UI || { sprak: "sv", texter: {} };

// Versionsnumret i adressen när sidorna kommer från GitHub Pages ("?v=...",
// se sida.go), så att det skripten själva laddar kommer från samma version.
// Tomt från tornets program.
window.ORMTORNET_V = document.currentScript ? new URL(document.currentScript.src).search : "";

// t ger texten med nyckeln key, med {namn} och liknande ifyllda från vars.
function t(key, vars) {
  let s = UI.texter[key];
  if (s === undefined) return "[" + key + "]";
  if (vars) {
    for (const k of Object.keys(vars)) s = s.split("{" + k + "}").join(String(vars[k]));
  }
  return s;
}

// applyTexts fyller i texterna i HTML:en: data-t ger textinnehållet,
// data-t-html innehåll med enkel HTML, och data-t-title, data-t-placeholder
// och data-t-aria attributen title, placeholder och aria-label.
function applyTexts(root) {
  root = root || document;
  document.documentElement.lang = UI.sprak;
  root.querySelectorAll("[data-t]").forEach((el) => (el.textContent = t(el.dataset.t)));
  root.querySelectorAll("[data-t-html]").forEach((el) => (el.innerHTML = t(el.dataset.tHtml)));
  root.querySelectorAll("[data-t-title]").forEach((el) => (el.title = t(el.dataset.tTitle)));
  root.querySelectorAll("[data-t-placeholder]").forEach((el) => (el.placeholder = t(el.dataset.tPlaceholder)));
  root.querySelectorAll("[data-t-aria]").forEach((el) => el.setAttribute("aria-label", t(el.dataset.tAria)));
}

// refreshTexts hämtar texterna igen om språket lang inte är det sidan har.
// Svarar true om texterna byttes; då ska sidan rita om det den själv har
// skrivit med t().
let textsLoading = false;
async function refreshTexts(page, lang) {
  if (!lang || lang === UI.sprak || textsLoading) return false;
  textsLoading = true;
  try {
    const r = await fetch("/texter/" + page + ".json", { cache: "no-store" });
    if (!r.ok) return false;
    window.UI = await r.json();
    applyTexts();
    return true;
  } catch {
    return false;
  } finally {
    textsLoading = false;
  }
}

applyTexts();
