// Matches the Go implementation in app/resources/mark/slugify.go
const nonLetterNumberPattern = /[^\p{L}\p{M}\p{N}\-_]+/gu;
const multiHyphenPattern = /-+/g;
const trailingSeparatorPattern = /[^\p{L}\p{M}\p{N}]+$/u;

export function slugify(input: string): string {
  // Trim leading and trailing whitespace
  const trimmed = input.trim();

  // NFKC normalization
  const normalized = trimmed.normalize("NFKC");

  // Lowercase
  const lowercased = normalized.toLowerCase();

  // Replace non-letter/number chars with hyphens
  const lettersReplaced = lowercased.replace(nonLetterNumberPattern, "-");

  // Collapse multiple hyphens
  const collapsed = lettersReplaced.replace(multiHyphenPattern, "-");

  // Trim leading and trailing hyphens/underscores
  const trimmedDividers = collapsed.replace(/^[-_]+|[-_]+$/g, "");

  return trimmedDividers;
}

export function slugifyDraft(input: string): string {
  const slug = slugify(input);
  const normalized = input.normalize("NFKC");

  // Keep one pending divider so typing the next word cannot join both words.
  if (slug && trailingSeparatorPattern.test(normalized)) {
    const divider = normalized.endsWith("_") ? "_" : "-";
    return `${slug}${divider}`;
  }

  return slug;
}

export function isSlug(input: string): boolean {
  return slugify(input) === input;
}
