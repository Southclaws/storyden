/** Returns a safe origin-relative path, or the homepage for invalid input. */
export function getSafeReturnPath(
  value: string | string[] | null | undefined,
): string {
  if (
    typeof value !== "string" ||
    !value.startsWith("/") ||
    value.startsWith("//") ||
    value.includes("\\") ||
    Array.from(value).some((char) => {
      const code = char.charCodeAt(0);
      return code < 0x20 || code === 0x7f;
    })
  ) {
    return "/";
  }

  return value;
}
