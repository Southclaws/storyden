const SAFE_URL_SCHEMES = ["http:", "https:"];

/**
 * Validates that a URL has a safe scheme (http/https only).
 * Rejects dangerous schemes like javascript:, data:, ftp:, etc.
 *
 * @param url - The URL to validate
 * @returns true if the URL has a safe scheme or no explicit scheme, false otherwise
 */
function isValidScheme(url: string): boolean {
  try {
    const parsed = new URL(url);
    return SAFE_URL_SCHEMES.includes(parsed.protocol);
  } catch {
    const hasUnsafeProtocol = /^[a-zA-Z][a-zA-Z0-9+.-]*:/.test(url);
    return !hasUnsafeProtocol;
  }
}

/**
 * Checks if input looks like a valid link-like string.
 * A "link-like" string is a domain-only input that users type without a protocol,
 * such as "example.com" or "subdomain.example.com", which can be normalized to a full URL.
 *
 * Requirements:
 * - Must contain at least one dot (to distinguish from single words)
 * - Cannot contain whitespace
 * - Must have a safe scheme if a protocol is specified
 *
 * @param input - The string to validate
 * @returns true if the input looks like a valid link
 *
 * @example
 * isValidLinkLike("example.com") // true
 * isValidLinkLike("https://example.com") // true
 * isValidLinkLike("javascript:alert(1)") // false
 * isValidLinkLike("hello world") // false (contains space)
 * isValidLinkLike("example") // false (no dot)
 */
export function isValidLinkLike(input: string): boolean {
  const trimmed = input.trim();
  if (!trimmed) {
    return false;
  }

  if (/\s/.test(input)) {
    return false;
  }

  if (!input.includes(".")) {
    return false;
  }

  return isValidScheme(trimmed);
}

/**
 * Normalizes a link-like string into a valid URL.
 * If the input is already a full URL, validates it and returns as-is.
 * If the input is a domain-only string (e.g., "example.com"), prepends "https://" and normalizes.
 *
 * Security: Only allows http: and https: schemes. Rejects javascript:, data:, ftp:, etc.
 *
 * @param input - The link-like string to normalize
 * @returns A normalized URL string, or undefined if invalid
 *
 * @example
 * normalizeLink("example.com") // "https://example.com/"
 * normalizeLink("https://example.com") // "https://example.com/"
 * normalizeLink("javascript:alert(1)") // undefined (unsafe scheme)
 * normalizeLink("   ") // undefined (empty after trim)
 * normalizeLink("not a link") // undefined (contains space)
 */
export function normalizeLink(input: string | undefined): string | undefined {
  if (!input) {
    return undefined;
  }

  if (!isValidLinkLike(input)) {
    return undefined;
  }

  const trimmed = input.trim();

  if (!trimmed) {
    return undefined;
  }

  try {
    const normalized = new URL(trimmed).toString();
    return isValidScheme(normalized) ? normalized : undefined;
  } catch {
    //
  }

  try {
    return new URL(`https://${trimmed}`).toString();
  } catch {
    //
  }

  return undefined;
}
