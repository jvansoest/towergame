import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { getTypes, whenTypes } from "./catalog";
import { roomCardURL } from "./roomCard";
import { variantCount } from "./roomArt";
import "./index.css";

// /card lists every card; /card/<type> shows
// one big, with the whole set small below.
function cardRequest(): { type: string; stair: boolean } | null {
  const path = window.location.pathname;
  if (/^\/card\/?$/.test(path)) return { type: "", stair: false };
  const m = path.match(/^\/card\/([\w-]+)\/?$/);
  if (!m) return null;
  return {
    type: m[1],
    stair: new URLSearchParams(window.location.search).has("stair"),
  };
}

// One small card with its id, linking to its page.
function thumb(id: string): HTMLAnchorElement | null {
  const url = roomCardURL(id);
  if (!url) return null;
  const img = document.createElement("img");
  img.src = url;
  img.title = id;
  img.style.height = "48px";
  img.style.imageRendering = "pixelated";
  const label = document.createElement("div");
  label.textContent = id;
  label.style.fontSize = "10px";
  label.style.textAlign = "center";
  label.style.color = "#6b7280";
  const cell = document.createElement("div");
  cell.appendChild(img);
  cell.appendChild(label);
  const a = document.createElement("a");
  a.href = `/card/${id}`;
  a.appendChild(cell);
  return a;
}

// Every card, small, along the bottom edge.
function strip(): HTMLDivElement {
  const bar = document.createElement("div");
  bar.style.display = "flex";
  bar.style.flexWrap = "wrap";
  bar.style.gap = "10px";
  bar.style.justifyContent = "center";
  bar.style.alignItems = "flex-end";
  bar.style.padding = "10px";
  bar.style.borderTop = "1px solid #e5e7eb";
  for (const t of getTypes()) {
    const a = thumb(t.id);
    if (a) bar.appendChild(a);
  }
  return bar;
}

// A bare page: the card up top, the set below.
function renderCard(type: string, stair: boolean) {
  document.title = type ? `card: ${type}` : "room cards";
  document.body.style.margin = "0";
  document.body.style.minHeight = "100vh";
  document.body.style.display = "grid";
  document.body.style.gridTemplateRows = type ? "1fr auto" : "auto";
  document.body.style.background = "#ffffff";

  const mount = () => {
    document.body.textContent = "";
    const stage = document.createElement("div");
    stage.style.display = "grid";
    stage.style.placeItems = "center";
    if (type) {
      // One big card per look of the art.
      const looks = variantCount(type);
      for (let v = 0; v < looks; v++) {
        const url = roomCardURL(type, { stair, variant: v });
        if (!url) {
          stage.textContent = `no card for "${type}"`;
          break;
        }
        const img = document.createElement("img");
        img.src = url;
        img.style.maxWidth = "92vw";
        img.style.imageRendering = "pixelated";
        img.style.margin = "6px";
        stage.appendChild(img);
      }
      document.body.appendChild(stage);
      document.body.appendChild(strip());
    } else {
      stage.appendChild(strip());
      document.body.appendChild(stage);
    }
  };
  // Cards need the catalog, which rides the
  // first snapshot from the game server.
  if (getTypes().length) mount();
  else whenTypes(mount);
}

const card = cardRequest();
const container = document.getElementById("root");
if (!container) throw new Error("Root element #root not found");

if (card) {
  renderCard(card.type, card.stair);
} else {
  // The game bundle stays off the card pages.
  void import("./App").then(({ default: App }) => {
    createRoot(container).render(
      <StrictMode>
        <App />
      </StrictMode>,
    );
  });
}
