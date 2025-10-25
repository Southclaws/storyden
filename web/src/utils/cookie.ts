export type CookieSetOptions = {
  days?: number; // expiration in days
  path?: string;
  domain?: string;
  sameSite?: "lax" | "strict" | "none";
  secure?: boolean;
};

export function getCookie(name: string): string | undefined {
  if (typeof document === "undefined") return undefined;
  const nameEQ = `${encodeURIComponent(name)}=`;
  const parts = document.cookie.split("; ");
  for (const part of parts) {
    if (part.startsWith(nameEQ)) {
      return decodeURIComponent(part.substring(nameEQ.length));
    }
  }
  return undefined;
}

export function setCookie(
  name: string,
  value: string,
  options: CookieSetOptions = {},
): void {
  if (typeof document === "undefined") return;

  const {
    days,
    path = "/",
    domain,
    sameSite = "lax",
    secure = typeof window !== "undefined" &&
      window.location.protocol === "https:",
  } = options;

  let cookie = `${encodeURIComponent(name)}=${encodeURIComponent(value)}`;
  if (days != null) {
    const date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    cookie += `; Expires=${date.toUTCString()}`;
  }
  if (path) cookie += `; Path=${path}`;
  if (domain) cookie += `; Domain=${domain}`;
  if (sameSite)
    cookie += `; SameSite=${sameSite.charAt(0).toUpperCase()}${sameSite.slice(1)}`;
  // Ensure Secure when SameSite=None to satisfy browser requirements
  const effectiveSecure = secure || sameSite === "none";
  if (effectiveSecure) cookie += `; Secure`;

  document.cookie = cookie;
}
