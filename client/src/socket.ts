// WebSocket wrapper with emit/on/off.
// Every frame is { type, ...payload }.

type Handler = (payload: Record<string, unknown>) => void;

const URL = import.meta.env.VITE_SERVER_URL ?? "ws://localhost:7777/ws";

class GameSocket {
  private ws: WebSocket | null = null;
  private handlers = new Map<string, Set<Handler>>();
  private outbox: string[] = [];
  private reconnectDelay = 500;

  constructor(private url: string) {
    this.connect();
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
    } else {
      this.outbox.push(frame);
    }
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
