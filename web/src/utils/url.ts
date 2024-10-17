import { WEB_ADDRESS } from "@/config";

export function isExternalURL(url: string) {
  if (url.startsWith("/")) {
    return false;
  }

  try {
    const u = new URL(url);

    const host = u.host;

    if (WEB_ADDRESS.includes(host)) {
      return false;
    }

    return true;
  } catch (_) {
    return false;
  }
}
