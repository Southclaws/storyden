export function isShareEnabled() {
  if (!isClient()) {
    return false;
  }

  return Boolean(navigator.share);
}

export function isClient() {
  return typeof window !== "undefined";
}
