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

  return true;
}

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
    return new URL(trimmed).toString();
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
