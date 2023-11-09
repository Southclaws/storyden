export function isWebauthnAvailable() {
  if (navigator === undefined) {
    return false;
  }
  if (window === undefined) {
    return false;
  }

  if (window.PublicKeyCredential === undefined) {
    return false;
  }

  const ua = navigator.userAgent.toLowerCase();

  // Disable on all Android devices until it's ready for regular users.
  if (ua.includes("android")) {
    return false;
  }

  return true;
}
