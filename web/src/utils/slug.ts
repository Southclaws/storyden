/**
 * Simple slug utility for generating URL-safe slugs
 * Uses String.prototype.normalize, regex replacement, and ASCII folding
 */
export function slugify(text: string): string {
  return (
    text
      // Normalize Unicode characters (NFD decomposition)
      .normalize("NFD")
      // Convert to lowercase
      .toLowerCase()
      // Remove diacritics/accents (anything that's not a basic Latin letter)
      .replace(/[\u0300-\u036f]/g, "")
      // Replace spaces and underscores with hyphens
      .replace(/[\s_]+/g, "-")
      // Remove all non-alphanumeric characters except hyphens
      .replace(/[^a-z0-9-]/g, "")
      // Replace multiple consecutive hyphens with single hyphen
      .replace(/-+/g, "-")
      // Remove leading and trailing hyphens
      .replace(/^-+|-+$/g, "")
  );
}
