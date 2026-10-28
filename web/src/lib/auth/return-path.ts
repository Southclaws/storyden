export function getSafeReturnPath(value: string | null | undefined): string {
  // Keep post-auth navigation on this origin, including after URL parsing.
  if (
    !value?.startsWith("/") ||
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
