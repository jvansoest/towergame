// Remembers what the inspector asked for, so the
// card can ask again and show fresh numbers.

import socket from "./socket";

type Ask = { type: string; payload: Record<string, unknown> };

let last: Ask | null = null;

// Asks about one thing and remembers the question.
export function askInspect(type: string, payload: Record<string, unknown>) {
  last = { type, payload };
  socket.emit(type, payload);
}

// Asks the same question again, without an error on failure.
export function repeatInspect() {
  if (!last) return;
  socket.emit(last.type, { ...last.payload, quiet: true });
}

// Forgets the question when the card closes.
export function stopInspect() {
  last = null;
}
