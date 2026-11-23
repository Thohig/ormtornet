// Webbläsarprovet: startar tornet och spelar en kort vandring i två
// webbläsare, en dator och en telefon, och letar efter det som robotarna
// inte ser: fel i skripten, texter på sidan som är trasiga, tangenterna vid
// dörrarna som ligger över varandra, en sida som är bredare än skärmen,
// minikort som inte går att öppna, mellanslag som inte spelar efter ett
// minikort och en omladdning som inte berättar rummet igen. Det tar en
// bild av varje fas i båda webbläsarna. Båda står på tornets dator och är
// alltså spelledaren som spelar själv (med märket i översta raden). Se
// README, "Webbläsarprovet".
//
//   node verktyg/webbprov.mjs [mapp] [sv|en]
//
// Mappen (förval webbprov/) får bilderna och rapport.txt. Provet slutar
// med 1 om något hittades. Kräver Node 22, Go och Playwright med Chromium
// (npm install playwright i klient/ eller globalt).
import { spawn, execFileSync } from "node:child_process";
import { createRequire } from "node:module";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { kontrollera } from "./textkontroll.mjs";

const rot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const mapp = path.resolve(process.argv[2] || "webbprov");
const sprak = process.argv[3] || "sv";
const port = 8190 + Math.floor(Math.random() * 100);
const bas = `http://127.0.0.1:${port}`;

// Playwright hämtas där det finns: i klient/, globalt eller bredvid.
function hamtaPlaywright() {
  for (const fran of [path.join(rot, "klient", "x.js"), path.join(process.cwd(), "x.js")]) {
    try { return createRequire(fran)("playwright"); } catch {}
  }
  try {
    const glob = execFileSync("npm", ["root", "-g"], { encoding: "utf8" }).trim();
    return createRequire(path.join(glob, "x.js"))("playwright");
  } catch {}
  console.error("Hittar inte Playwright. Kör npm install playwright i klient/ eller globalt.");
  process.exit(2);
}
const { chromium, devices } = hamtaPlaywright();

fs.mkdirSync(mapp, { recursive: true });
const rapport = [];
const fynd = [];
const logga = (s) => { rapport.push(s); console.log(s); };
const fynna = (vem, s) => {
  const rad = `FYND ${vem}: ${s}`;
  if (fynd.includes(rad)) return;
  fynd.push(rad);
  logga(rad);
};
const vanta = (ms) => new Promise((r) => setTimeout(r, ms));

// Tornet: byggs till en tillfällig fil och startas med en kort vandring,
// två människor och en AI-spelare.
const bin = path.join(os.tmpdir(), `ormtornet-webbprov-${process.pid}`);
logga(`bygger tornet …`);
execFileSync("go", ["build", "-o", bin, "."], { cwd: rot, stdio: "inherit" });
const server = spawn(bin, [
  "-addr", `127.0.0.1:${port}`, "-installningar", "", "-min", "3", "-ai", "1",
  "-sol", "8", "-turtid", "8s", "-sprak", sprak, "-autostopp", "0",
], { cwd: rot, stdio: ["ignore", "pipe", "pipe"] });
let serverlogg = "";
server.stdout.on("data", (d) => (serverlogg += d));
server.stderr.on("data", (d) => (serverlogg += d));
function stang(kod) {
  server.kill();
  try { fs.unlinkSync(bin); } catch {}
  fs.writeFileSync(path.join(mapp, "rapport.txt"), rapport.join("\n") + "\n");
  process.exit(kod);
}
for (let i = 0; ; i++) {
  try { await fetch(bas + "/"); break; } catch {}
  if (i > 100) { console.error(serverlogg); stang(2); }
  await vanta(100);
}

const browser = await chromium.launch({ args: ["--use-gl=angle", "--use-angle=swiftshader", "--enable-unsafe-swiftshader"] });

// En spelare: en sida i en egen webbläsare, med det den har sett.
async function spelare(namn, opt) {
  const ctx = await browser.newContext(opt);
  const page = await ctx.newPage();
  const s = { namn, page, sett: new Set(), bilder: new Set(), mobil: !!opt.isMobile };
  page.on("pageerror", (e) => fynna(namn, `skriptfel: ${e.message}`));
  page.on("console", (m) => {
    if (m.type() === "error") fynna(namn, `konsolen: ${m.text()}`);
  });
  page.on("response", (r) => {
    if (r.status() >= 400 && !r.url().endsWith("favicon.ico")) fynna(namn, `${r.status()} ${r.url()}`);
  });
  // Webbläsarna står på tornets dator, där introt leder till värdsidan;
  // ?spela går in som spelledaren gör med Gå in själv.
  await page.goto(bas + "/?spela");
  return s;
}

const dator = await spelare("datorn", { viewport: { width: 1440, height: 900 } });
const telefon = await spelare("telefonen", { ...devices["iPhone 13"], defaultBrowserType: undefined });
const alla = [dator, telefon];

// Läget i sidan: tillståndet från tornet och vad som syns.
const lage = (s) => s.page.evaluate(() => {
  const st = typeof state !== "undefined" ? state : null;
  const me = st && st.players ? st.players.find((p) => p.id === st.you) : null;
  return {
    joined: typeof joined !== "undefined" && joined,
    intro: document.getElementById("start").classList.contains("active"),
    phase: st ? st.phase : "",
    round: st ? st.round : -1,
    action: st ? st.action : "",
    inne: !!me && me.status === "inne",
    namn: st && st.players ? st.players.map((p) => p.name) : [],
  };
});

// granska letar i det som syns efter trasiga texter, en sida som är
// bredare än skärmen och tangenter vid dörrarna som ligger över varandra.
async function granska(s, fas) {
  const g = await s.page.evaluate(() => {
    const syns = (el) => {
      const r = el.getBoundingClientRect();
      const cs = getComputedStyle(el);
      return r.width > 0 && r.height > 0 && cs.visibility !== "hidden" && cs.display !== "none";
    };
    const tangenter = [...document.querySelectorAll(".isotangenter span[data-key]")]
      .filter(syns).map((el) => {
        const r = el.getBoundingClientRect();
        return { k: el.dataset.key, x: r.x, y: r.y, w: r.width, h: r.height };
      });
    const kartan = document.getElementById("mapbox").getBoundingClientRect();
    // Raden som loggen skriver ut just nu (.cursor) är inte färdig.
    const skrivs = [...document.querySelectorAll(".cursor")]
      .flatMap((e) => e.innerText.split("\n").map((r) => r.trim()));
    return {
      text: document.body.innerText,
      skrivs,
      bredd: document.documentElement.scrollWidth,
      skarm: innerWidth,
      tangenter,
      kartan: { x: kartan.x, y: kartan.y, w: kartan.width, h: kartan.height },
      utanWebgl: !document.getElementById("mapNone").hidden,
    };
  });
  const skrivs = new Set(g.skrivs);
  for (const rad of g.text.split("\n").map((r) => r.trim()).filter((r) => r.includes(" "))) {
    if (s.sett.has(rad) || skrivs.has(rad)) continue;
    s.sett.add(rad);
    const fel = kontrollera(rad, sprak);
    if (fel.length) fynna(s.namn, `trasig text i ${fas} [${fel.join(", ")}]: ${rad}`);
  }
  if (g.bredd > g.skarm + 1) fynna(s.namn, `sidan är ${g.bredd} bildpunkter bred på en skärm på ${g.skarm} i ${fas}`);
  if (g.utanWebgl) fynna(s.namn, "kartan ritades inte (ingen WebGL i webbläsaren)");
  const t = g.tangenter;
  for (let i = 0; i < t.length; i++) {
    for (let j = i + 1; j < t.length; j++) {
      const a = t[i], b = t[j];
      if (a.x < b.x + b.w && b.x < a.x + a.w && a.y < b.y + b.h && b.y < a.y + a.h) {
        fynna(s.namn, `tangenterna ${a.k} och ${b.k} ligger över varandra i ${fas}`);
      }
    }
  }
  return g;
}

async function bild(s, namn) {
  if (s.bilder.has(namn)) return;
  s.bilder.add(namn);
  await s.page.screenshot({ path: path.join(mapp, `${s.namn}-${namn}.png`) });
}

// tryck trycker på något på telefonen. Ligger något annat över det, som
// det stora kortet, väntar det högst tre sekunder och skriver vad som låg
// i vägen (en iakttagelse, inget fynd), och svarar false.
async function tryck(s, loc, vad) {
  try {
    await loc.tap({ timeout: 3000 });
    return true;
  } catch (e) {
    // Det som ligger över mitten av knappen, med närmaste id (#spelkort,
    // #storkort och så vidare).
    const over = await loc.evaluate((el) => {
      const r = el.getBoundingClientRect();
      const top = document.elementFromPoint(r.x + r.width / 2, r.y + r.height / 2);
      if (!top || el.contains(top)) return "";
      const med = top.closest("[id]");
      return med ? "#" + med.id : top.tagName.toLowerCase();
    }).catch(() => "");
    const rad = `OBS ${s.namn}: ${vad} gick inte att trycka på` + (over ? `; ${over} låg över` : "");
    if (!rapport.includes(rad)) logga(rad);
    return false;
  }
}
const tangent = (s, k) => s.page.evaluate((k) => pressKey(k), k);

// gaIn går förbi introt: Enter när knappen Gå in är redo.
async function gaIn(s) {
  // Med ?spela, eller efter en omladdning, går sidan ibland in av sig själv.
  await s.page.waitForFunction(() => joined || !document.getElementById("enter").classList.contains("waiting"), null, { timeout: 20000 });
  if (await s.page.evaluate(() => joined)) return;
  await bild(s, "1-introt");
  if (s.mobil) await s.page.tap("#enter");
  else await s.page.keyboard.press("Enter");
  await s.page.waitForFunction(() => typeof joined !== "undefined" && joined, null, { timeout: 15000 });
}

// draget för en spelare: datorn trycker en pil mot en öppen dörr, och
// telefonen trycker på tangenten vid en dörr (eller pilknapparna).
async function spela(s) {
  const val = await s.page.evaluate(() => {
    const me = state.players.find((p) => p.id === state.you);
    const rum = state.rooms.find((r) => r.x === me.x && r.y === me.y);
    const pilar = ["ArrowUp", "ArrowRight", "ArrowDown", "ArrowLeft"];
    const oppna = rum ? pilar.filter((_, i) => rum.d[i]) : [];
    return oppna.length ? oppna[Math.floor(Math.random() * oppna.length)] : " ";
  });
  if (s.mobil) {
    const vid = s.page.locator(`.isotangenter span[data-key="${val}"]`);
    if (await vid.count() && await vid.first().isVisible()) {
      if (!(await tryck(s, vid.first(), "tangenten vid dörren"))) await tangent(s, val);
      return;
    }
    const knapp = s.page.locator(`#pad button[data-key="${val}"]`);
    if (await knapp.count() && await knapp.first().isVisible()) {
      if (!(await tryck(s, knapp.first(), "pilknappen"))) await tangent(s, val);
      return;
    }
  }
  await s.page.keyboard.press(val);
}

// minikortet: ett klick på ett minikort i loggen ska öppna kortvisningen,
// ett klick till stänga den, och sedan ska mellanslag spela (vänta).
async function provaMinikort(s) {
  const mini = s.page.locator("#log .minikort[data-kort], #mapleft .minikort[data-kort]");
  const n = await mini.count();
  let synlig = null;
  for (let i = n - 1; i >= 0; i--) {
    if (await mini.nth(i).isVisible()) { synlig = mini.nth(i); break; }
  }
  if (!synlig) return false;
  // Rundan ska ha tid kvar, annars tar den slut innan mellanslaget syns.
  const kvar = await s.page.evaluate(() => state.deadlineMs - (performance.now() - state.receivedAt));
  if (kvar < 5000) return false;
  const runda0 = await s.page.evaluate(() => state.round);
  await synlig.click();
  if (await s.page.locator("#kortvisning").isHidden()) {
    fynna(s.namn, "ett klick på ett minikort öppnade inte kortvisningen");
    return true;
  }
  await bild(s, "4-kortvisningen");
  await s.page.click("#kortvisning");
  await s.page.keyboard.press(" ");
  try {
    // Draget syns, eller så var det det sista som behövdes och rundan
    // gick vidare.
    await s.page.waitForFunction((r0) => state.action === "vanta" || state.round > r0, runda0, { timeout: 4000 });
  } catch {
    const l = await s.page.evaluate((fore) => {
      const a = document.activeElement;
      return {
        fokus: a ? a.tagName.toLowerCase() + (a.id ? "#" + a.id : "") + (a.className ? "." + String(a.className).split(" ")[0] : "") : "",
        drag: state.action, runda: state.round, fore,
        kvar: Math.round(state.deadlineMs - (performance.now() - state.receivedAt)),
      };
    }, runda0);
    fynna(s.namn, `mellanslag spelade inte efter ett klick på ett minikort (fokus ${l.fokus}, drag "${l.drag}", runda ${l.fore} till ${l.runda}, ${l.kvar} ms kvar)`);
    await bild(s, "fel-minikortet");
  }
  if (await s.page.locator("#kortvisning").isVisible()) {
    fynna(s.namn, "kortvisningen öppnades igen av mellanslag");
  }
  return true;
}

let klar = false;
const slut = setTimeout(() => {
  fynna("provet", "vandringen blev inte klar på fem minuter");
  klar = true;
}, 5 * 60 * 1000);

try {
  await Promise.all(alla.map(gaIn));
  logga("båda är inne i tornet");
  const gjort = { minikort: false, omladdning: false };
  const runda = new Map();
  while (!klar) {
    let slutfas = 0;
    for (const s of alla) {
      const l = await lage(s);
      if (l.intro && !l.joined) { await gaIn(s); continue; }
      if (l.phase === "lobby") {
        await granska(s, "lobbyn");
        await bild(s, "2-lobbyn");
        if (!s.redo) {
          s.redo = true;
          if (s.mobil) { if (!(await tryck(s, s.page.locator("#padEnter"), "Redo"))) await tangent(s, "Enter"); }
          else await s.page.keyboard.press("Enter");
        }
      } else if (l.phase === "klassval") {
        await granska(s, "klassvalet");
        await vanta(1500); // korten vänds upp ett i taget
        await bild(s, "3-klassvalet");
        if (!s.klass) {
          s.klass = true;
          if (s.mobil) { if (!(await tryck(s, s.page.locator("#cards button.card").nth(1), "klasskortet"))) await tangent(s, "2"); }
          else await s.page.keyboard.press(s === dator ? "1" : "2");
        }
      } else if (l.phase === "spel") {
        await granska(s, `runda ${l.round}`);
        if (l.round === 2) await bild(s, "4-vandringen");
        if (!l.inne || l.action || runda.get(s) === l.round) continue;
        runda.set(s, l.round);
        if (s === dator && !gjort.minikort && l.round >= 2) {
          if (await provaMinikort(s)) { gjort.minikort = true; continue; }
        }
        if (s === dator && !gjort.omladdning && l.round >= 4) {
          gjort.omladdning = true;
          await s.page.reload();
          await s.page.waitForFunction(() => typeof state !== "undefined" && state && state.phase === "spel", null, { timeout: 15000 })
            .catch(async () => { if ((await lage(s)).intro) await gaIn(s); });
          // Loggen skriver raderna en i taget; rummet ska stå i den senaste
          // rundan inom några sekunder (längre på en dator med mycket att göra).
          const rummet = await s.page.waitForFunction(() => {
            const g = document.querySelector("#log section.round");
            return !!(g && g.querySelector(".roomblock"));
          }, null, { timeout: 20000 }).then(() => true, () => false);
          // I mörkret och i källarvalven ser vandraren inget, och då
          // berättas rummet inte (retellRoom i here.go).
          const blind = await s.page.evaluate(() => state.dark > 0 || state.cellar > 0);
          if (!rummet && !blind) fynna(s.namn, "efter en omladdning berättas inte rummet igen i loggen");
          await bild(s, "5-omladdningen");
          runda.delete(s);
          continue;
        }
        await spela(s);
      } else if (l.phase === "slut") {
        slutfas++;
        if (!s.slut) {
          await vanta(1500);
          await granska(s, "slutet");
          await bild(s, "6-slutet");
          const res = await s.page.locator("#chronResults").innerText().catch(() => "");
          for (const n of l.namn) if (!res.includes(n)) fynna(s.namn, `${n} saknas i resultatet`);
          s.slut = true;
        }
      }
    }
    if (slutfas === alla.length) break;
    await vanta(300);
  }
  if (!gjort.minikort) logga("minikorten provades inte: inget minikort syntes i loggen när rundan hade tid kvar");
} catch (e) {
  fynna("provet", `stannade: ${e.message.split("\n")[0]}`);
  for (const s of alla) await s.page.screenshot({ path: path.join(mapp, `${s.namn}-fel.png`) }).catch(() => {});
}
clearTimeout(slut);
await browser.close();
logga(fynd.length ? `${fynd.length} fynd. Bilderna står i ${mapp}` : `inga fynd. Bilderna står i ${mapp}`);
stang(fynd.length ? 1 : 0);
