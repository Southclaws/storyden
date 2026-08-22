export function capitalise(value: string): string {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

export function pluralise(
  count: number,
  singular: string,
  plural = `${singular}s`,
): string {
  return count === 1 ? singular : plural;
}

export function humanise(value: string): string {
  return capitalise(value.replaceAll("_", " "));
}

export function truncateText(
  value: string | undefined,
  maxLength = 280,
): string | undefined {
  if (!value) return undefined;

  const trimmed = value.trim();
  if (trimmed.length <= maxLength) return trimmed;

  return `${trimmed.slice(0, maxLength).trimEnd()}…`;
}
