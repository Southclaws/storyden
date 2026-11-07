let pending: Promise<unknown> = Promise.resolve();

export function withPushLock<T>(action: () => Promise<T>): Promise<T> {
  if (navigator.locks) {
    return navigator.locks.request("storyden-webpush", action);
  }

  const next = pending.then(action, action);
  pending = next;
  return next;
}
