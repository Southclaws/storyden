export function isLikelyInternalIP(ip: string): boolean {
  const trimmed = ip.trim();
  if (!trimmed) return false;

  if (trimmed === "::1") return true;
  if (trimmed.includes(":")) {
    const lower = trimmed.toLowerCase();
    return (
      lower.startsWith("fc") ||
      lower.startsWith("fd") ||
      lower.startsWith("fe80:")
    );
  }

  const parts = trimmed.split(".");
  if (parts.length !== 4) return false;

  const octets = parts.map((p) => Number.parseInt(p, 10));
  if (octets.some((n) => Number.isNaN(n) || n < 0 || n > 255)) return false;

  const a = octets[0] ?? -1;
  const b = octets[1] ?? -1;
  if (a === 10) return true;
  if (a === 127) return true;
  if (a === 192 && b === 168) return true;
  if (a === 169 && b === 254) return true;
  if (a === 172 && b >= 16 && b <= 31) return true;
  if (a === 100 && b >= 64 && b <= 127) return true;

  return false;
}

export function isPublicSingleHost(cidr: string): boolean {
  const [address = "", bits] = cidr.trim().split("/");

  if (!address || isLikelyInternalIP(address)) {
    return false;
  }

  if (bits === undefined) {
    return true;
  }

  return address.includes(":") ? bits === "128" : bits === "32";
}

export function parseForwardedFor(value: string | undefined): string[] {
  return (value ?? "")
    .split(",")
    .map((hop) => hop.trim())
    .filter(Boolean);
}
