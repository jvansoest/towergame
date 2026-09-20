// Hand-authored pixel art: a palette plus rows
// of palette keys, one character per pixel,
// '.' meaning clear. Drawn once to its own
// canvas, then stamped scaled and crisp.

export type Pixels = { palette: Record<string, string>; rows: string[] };

const natives = new WeakMap<object, HTMLCanvasElement>();

// The sprite at one pixel per art cell.
export function pixelCanvas(art: Pixels): HTMLCanvasElement {
  let cv = natives.get(art);
  if (cv) return cv;
  cv = document.createElement("canvas");
  cv.width = art.rows[0].length;
  cv.height = art.rows.length;
  const g = cv.getContext("2d");
  if (!g) return cv;
  art.rows.forEach((row, r) => {
    for (let c = 0; c < row.length; c++) {
      const col = art.palette[row[c]];
      if (!col) continue;
      g.fillStyle = col;
      g.fillRect(c, r, 1, 1);
    }
  });
  natives.set(art, cv);
  return cv;
}

// Stamps a sprite onto a room texture,
// stretched to fill, pixels kept square.
export function stampPixels(
  g: CanvasRenderingContext2D,
  art: Pixels,
  w: number,
  h: number,
) {
  g.imageSmoothingEnabled = false;
  g.drawImage(pixelCanvas(art), 0, 0, w, h);
  g.imageSmoothingEnabled = true;
}

// The izakaya wall: wood planks, a wide slit
// window, noren, kanji lanterns, a bottle shelf
// and a barrel. Built on a grid, 50x33.
const IZ_W = 50;
const IZ_H = 33;

// 酒, five by six, X marks ink.
const SAKE = ["..XXX", "X.X.X", "..XXX", ".XX.X", "X.X.X", "..XXX"];

// A char grid with rect and glyph painters.
function gridOf(w: number, h: number, fill: string) {
  const g: string[][] = Array.from({ length: h }, () =>
    Array<string>(w).fill(fill),
  );
  const rect = (x: number, y: number, rw: number, rh: number, ch: string) => {
    for (let r = y; r < y + rh; r++) {
      for (let c = x; c < x + rw; c++) {
        if (r >= 0 && r < h && c >= 0 && c < w) g[r][c] = ch;
      }
    }
  };
  const blit = (x: number, y: number, rows: string[], ink: string) => {
    rows.forEach((row, r) => {
      for (let c = 0; c < row.length; c++) {
        if (row[c] === "X") rect(x + c, y + r, 1, 1, ink);
      }
    });
  };
  return { g, rect, blit };
}

function izakayaGrid(): string[] {
  const { g, rect, blit } = gridOf(IZ_W, IZ_H, "W");

  // Planks, beam, wainscot, floor.
  for (let c = 0; c < IZ_W; c += 9) rect(c, 4, 1, 24, "V");
  rect(0, 0, IZ_W, 3, "B");
  rect(0, 3, IZ_W, 1, "T");
  rect(0, 22, IZ_W, 1, "T");
  rect(0, 23, IZ_W, 5, "P");
  rect(0, 28, IZ_W, 5, "F");
  rect(0, 28, IZ_W, 1, "S");
  for (let c = 6; c < IZ_W; c += 12) rect(c, 29, 1, 4, "f");

  // A wide, narrow window with a sill.
  rect(17, 4, 32, 5, "D");
  rect(18, 5, 30, 3, "G");
  rect(18, 5, 30, 1, "H");
  for (const c of [27, 38]) rect(c, 5, 1, 3, "D");
  rect(16, 9, 34, 1, "S");

  // A noren curtain: two flaps, a white mark.
  for (let i = 0; i < 2; i++) {
    const x = 2 + i * 5;
    rect(x, 4, 4, 12, "N");
    rect(x + 1, 8, 2, 2, "n");
  }
  rect(1, 4, 11, 1, "B");

  // A shelf of bottles behind the server.
  const bottles = ["g", "b", "c", "g", "b"];
  rect(37, 16, 12, 1, "S");
  bottles.forEach((k, i) => {
    const x = 38 + i * 2;
    rect(x, 11, 2, 5, k);
    rect(x, 10, 1, 1, k);
  });

  // A sake barrel with two hoops.
  rect(5, 23, 10, 5, "O");
  rect(5, 24, 10, 1, "h");
  rect(5, 26, 10, 1, "h");
  rect(5, 23, 1, 5, "o");
  rect(14, 23, 1, 5, "o");

  // Kanji lanterns hang in front of the window.
  for (const x of [19, 29, 39]) {
    rect(x + 3, 3, 1, 2, "B");
    rect(x + 1, 5, 5, 1, "B");
    rect(x, 6, 7, 8, "R");
    rect(x, 6, 1, 8, "r");
    rect(x + 6, 6, 1, 8, "r");
    blit(x + 1, 7, SAKE, "Y");
    rect(x + 1, 14, 5, 1, "B");
  }
  return g.map((r) => r.join(""));
}

export const IZAKAYA_WALL: Pixels = {
  palette: {
    W: "#5b3a29",
    V: "#4d3123",
    B: "#2b1a12",
    T: "#7a4f30",
    P: "#4a2c1f",
    F: "#3a2416",
    f: "#4a3020",
    S: "#8b5a33",
    D: "#2b1a12",
    G: "rgba(143,196,224,0.22)",
    H: "rgba(169,215,239,0.4)",
    N: "#2a3f6b",
    n: "#e8e2d0",
    M: "#2b1a12",
    y: "#e8e2d0",
    g: "#6fae7a",
    b: "#c9a86a",
    c: "#9cc9e0",
    O: "#b98a4e",
    o: "#8c6236",
    h: "#2b1a12",
    R: "#d94f34",
    r: "#a83a24",
    Y: "#f6e0a8",
  },
  rows: izakayaGrid(),
};

// The counter-joint walls share one layout: a
// striped awning, a sign, a wide slit window, a
// menu (or oven), wainscot, a checked floor. 50x33.
type Rect = (x: number, y: number, w: number, h: number, ch: string) => void;

function jointGrid(
  sign: (r: Rect) => void,
  decor: (r: Rect) => void,
): string[] {
  const { g, rect } = gridOf(IZ_W, IZ_H, "W");
  for (let c = 0; c < IZ_W; c += 6) rect(c, 5, 3, 17, "V");
  rect(0, 22, IZ_W, 1, "a");
  rect(0, 23, IZ_W, 5, "P");
  for (let y = 28; y < IZ_H; y += 2) {
    for (let x = 0; x < IZ_W; x += 4) {
      rect(x, y, 4, 2, ((x / 4 + (y - 28) / 2) & 1) === 0 ? "a" : "f");
    }
  }

  // Awning: stripes with a scalloped edge.
  for (let x = 0; x < IZ_W; x += 4) {
    const ch = (x / 4) % 2 === 0 ? "A" : "a";
    rect(x, 0, 4, 5, ch);
    rect(x, 5, 4, 1, ch);
    rect(x + 1, 6, 2, 1, ch);
  }

  // A wide, narrow window.
  rect(17, 8, 19, 5, "a");
  rect(18, 9, 17, 3, "G");
  rect(18, 9, 17, 1, "H");
  rect(26, 9, 1, 3, "a");
  rect(16, 13, 21, 1, "S");

  sign(rect);
  decor(rect);
  return g.map((r) => r.join(""));
}

// A chalk menu behind the server.
const chalkMenu = (rect: Rect) => {
  rect(37, 7, 12, 11, "M");
  rect(38, 8, 10, 9, "C");
  for (const y of [10, 12, 14, 16]) {
    rect(39, y, 4, 1, "y");
    rect(44, y, 3, 1, "z");
  }
  rect(41, 8, 4, 1, "z");
};

// A pink scoop on a waffle cone.
const coneSign = (rect: Rect) => {
  rect(8, 7, 1, 2, "M");
  rect(6, 8, 5, 1, "p");
  rect(5, 9, 7, 1, "p");
  rect(4, 10, 9, 3, "p");
  rect(4, 10, 3, 1, "q");
  rect(5, 13, 7, 1, "q");
  for (let i = 0; i < 8; i++) {
    rect(4 + Math.floor(i / 2), 14 + i, 9 - i, 1, i % 2 ? "T" : "t");
  }
};

// A burger: bun, lettuce, cheese, patty, bun.
const burgerSign = (rect: Rect) => {
  rect(8, 7, 1, 2, "M");
  rect(6, 8, 5, 1, "t");
  rect(5, 9, 7, 3, "t");
  rect(7, 10, 1, 1, "y");
  rect(10, 10, 1, 1, "y");
  rect(8, 9, 1, 1, "y");
  rect(4, 12, 9, 1, "g");
  rect(4, 13, 9, 1, "z");
  rect(5, 14, 7, 2, "q");
  rect(5, 16, 7, 2, "t");
  rect(6, 18, 5, 1, "T");
};

// A pizza slice with pepperoni.
const pizzaSign = (rect: Rect) => {
  rect(8, 7, 1, 1, "M");
  rect(4, 8, 9, 2, "T");
  for (let i = 0; i < 9; i++) {
    rect(4 + Math.floor(i / 2), 10 + i, 9 - i, 1, "p");
  }
  rect(6, 11, 2, 2, "q");
  rect(9, 12, 2, 2, "q");
  rect(7, 14, 2, 2, "q");
};

// A steaming bowl with chopsticks.
const noodleSign = (rect: Rect) => {
  rect(8, 7, 1, 1, "M");
  for (let r = 8; r < 12; r++) {
    rect(6 + (r % 2), r, 1, 1, "y");
    rect(9 + ((r + 1) % 2), r, 1, 1, "y");
  }
  rect(5, 12, 7, 2, "t");
  rect(4, 14, 9, 1, "a");
  rect(5, 15, 7, 3, "B");
  rect(6, 18, 5, 1, "B");
  rect(12, 8, 1, 5, "M");
  rect(14, 8, 1, 5, "M");
};

// A brick oven with a fire, behind the server.
const brickOven = (rect: Rect) => {
  rect(37, 7, 12, 12, "O");
  for (const y of [9, 12, 15]) rect(37, y, 12, 1, "o");
  rect(40, 11, 6, 8, "k");
  rect(41, 12, 4, 1, "k");
  rect(41, 15, 4, 4, "e");
  rect(42, 16, 2, 3, "z");
};

const jointBase = {
  G: "rgba(143,196,224,0.22)",
  H: "rgba(169,215,239,0.4)",
  a: "#fbf6ee",
  M: "#8b5a33",
  y: "#fbf6ee",
};

export const ICECREAM_WALL: Pixels = {
  palette: {
    ...jointBase,
    W: "#cfeee0",
    V: "#c2e6d6",
    A: "#ee8fb0",
    P: "#f4b6c8",
    f: "#f2a3bd",
    S: "#e8dcc8",
    C: "#2f4a3f",
    z: "#f7a8c8",
    p: "#f7a8c8",
    q: "#e0709a",
    t: "#d9a25a",
    T: "#b9803a",
  },
  rows: jointGrid(coneSign, chalkMenu),
};

export const BURGER_WALL: Pixels = {
  palette: {
    ...jointBase,
    W: "#f6e9cf",
    V: "#efdcb8",
    A: "#d94f34",
    P: "#b83a2a",
    f: "#3c3f46",
    S: "#e8dcc8",
    C: "#2b2b2b",
    z: "#f2c230",
    g: "#6fae3f",
    q: "#6b3b22",
    t: "#e0a458",
    T: "#c98a3a",
  },
  rows: jointGrid(burgerSign, chalkMenu),
};

export const PIZZA_WALL: Pixels = {
  palette: {
    ...jointBase,
    W: "#e6b48c",
    V: "#dea67b",
    A: "#2f7d4a",
    P: "#2a6a3f",
    f: "#c0392b",
    S: "#e8dcc8",
    z: "#f6c453",
    p: "#f2b134",
    q: "#c0392b",
    T: "#c98a3a",
    O: "#a8523a",
    o: "#8c4230",
    k: "#2b1a12",
    e: "#e8742a",
  },
  rows: jointGrid(pizzaSign, brickOven),
};

export const NOODLE_WALL: Pixels = {
  palette: {
    ...jointBase,
    W: "#efe6d2",
    V: "#e6dbc2",
    A: "#c0392b",
    P: "#2f5d62",
    f: "#d8b47c",
    S: "#e8dcc8",
    C: "#3a2a22",
    z: "#e8b04a",
    t: "#e8c47a",
    B: "#3b6ea5",
  },
  rows: jointGrid(noodleSign, chalkMenu),
};
