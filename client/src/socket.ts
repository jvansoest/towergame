// WebSocket wrapper with emit/on/off.
// Every frame is { type, ...payload }.

type Handler = (payload: Record<string, unknown>) => void;

const URL = import.meta.env.VITE_SERVER_URL ?? "ws://localhost:7777/ws";

// Silence this long means a dead link.
const DEAD_AFTER_MS = 20_000;
// How often to check for silence.
const WATCHDOG_EVERY_MS = 5_000;

class GameSocket {
  private ws: WebSocket | null = null;
  private handlers = new Map<string, Set<Handler>>();
  private outbox: string[] = [];
  private reconnectDelay = 500;
  private lastRx = Date.now();

  constructor(private url: string) {
    this.connect();
    this.watchdog();
    this.resyncOnFocus();
  }

  private connect() {
    const ws = new WebSocket(this.url);
    this.ws = ws;

    ws.onopen = () => {
      this.reconnectDelay = 500;
      for (const frame of this.outbox) ws.send(frame);
      this.outbox = [];
    };

    ws.onmessage = (ev) => {
      this.lastRx = Date.now();
      let msg: Record<string, unknown>;
      try {
        msg = JSON.parse(ev.data as string);
      } catch {
        return;
      }
      const type = msg.type as string | undefined;
      if (!type) return;
      this.handlers.get(type)?.forEach((h) => h(msg));
    };

    ws.onclose = () => {
      // Reconnect with backoff.
      setTimeout(() => this.connect(), this.reconnectDelay);
      this.reconnectDelay = Math.min(this.reconnectDelay * 2, 8000);
    };

    ws.onerror = () => ws.close();
  }

  // Send an intent to the server.
  emit(type: string, payload: Record<string, unknown> = {}) {
    const frame = JSON.stringify({ type, ...payload });
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(frame);
    } else if (type !== "resync") {
      // A fresh registration resyncs anyway.
      this.outbox.push(frame);
    }
  }

  // Hang up on silence, so onclose redials.
  private watchdog() {
    setInterval(() => {
      if (Date.now() - this.lastRx < DEAD_AFTER_MS) return;
      // Re-arm now; backoff paces the retries.
      this.lastRx = Date.now();
      this.ws?.close();
    }, WATCHDOG_EVERY_MS);
  }

  // Wake from a hidden tab with fresh state.
  private resyncOnFocus() {
    document.addEventListener("visibilitychange", () => {
      if (document.visibilityState === "visible") this.emit("resync");
    });
  }

  on(type: string, handler: Handler) {
    let set = this.handlers.get(type);
    if (!set) {
      set = new Set();
      this.handlers.set(type, set);
    }
    set.add(handler);
  }

  off(type: string, handler: Handler) {
    this.handlers.get(type)?.delete(handler);
  }
}

const socket = new GameSocket(URL);

export default socket;
