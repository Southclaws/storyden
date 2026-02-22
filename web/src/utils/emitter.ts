export type Emitter<TEvents extends Record<string | symbol, unknown>> = {
  on<K extends keyof TEvents>(
    type: K,
    handler: (event: TEvents[K]) => void,
  ): void;
  off<K extends keyof TEvents>(
    type: K,
    handler: (event: TEvents[K]) => void,
  ): void;
  emit<K extends keyof TEvents>(type: K, event: TEvents[K]): void;
};

export function createEmitter<
  TEvents extends Record<string | symbol, unknown>,
>(): Emitter<TEvents> {
  const listeners = new Map<keyof TEvents, Set<(event: unknown) => void>>();

  function on<K extends keyof TEvents>(
    type: K,
    handler: (event: TEvents[K]) => void,
  ) {
    let set = listeners.get(type);
    if (!set) {
      set = new Set();
      listeners.set(type, set);
    }
    set.add(handler as (event: unknown) => void);
  }

  function off<K extends keyof TEvents>(
    type: K,
    handler: (event: TEvents[K]) => void,
  ) {
    const set = listeners.get(type);
    if (!set) return;
    set.delete(handler as (event: unknown) => void);
    if (set.size === 0) listeners.delete(type);
  }

  function emit<K extends keyof TEvents>(type: K, event: TEvents[K]) {
    const set = listeners.get(type);
    if (!set) return;
    // Iterate a snapshot, but consult the live set so removals during dispatch
    // can prevent later handlers from firing in this cycle.
    [...set].forEach((fn) => {
      if (set.has(fn)) {
        fn(event);
      }
    });
  }

  return { on, off, emit };
}
