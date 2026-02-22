const validSlugRegex = /^[^\s/?#%]+$/;
const disallowedChars = /[/?#%]/g;
const spaceRegex = /\s+/g;
const multiDashRegex = /-+/g;
const trailingDashRegex = /-+$/;

/**
 * Processes a raw string input to create a slug suitable for use as a Mark. It
 * is designed to be used with an <input> value prop, allowing the typing state
 * to continue in a flowing way while outputting possibly invalid slugs.
 *
 * Invalid slugs are still stored in state, they are still processed partially
 * in order to allow the user to type "my-page" by not stripping the trailing
 * hyphen when they are at the "my-" stage of typing. Even though this will be
 * rejected by the API, the pre-patch function further processes the value by
 * removing the trailing hyphen, so that the user can continue typing easily.
 *
 * @param raw raw string input to be processed into a slug from an input box.
 * @returns the partially well-formed slug, spaces replaced by hyphens, multiple
 *          hyphens collapsed into one, disallowed characters removed, and
 *          converted to lowercase.
 */
export function processMarkInput(raw: string): string {
  const noSpaces = raw.replace(spaceRegex, "-");
  const cleaned = noSpaces.replace(disallowedChars, "");
  const collapsed = cleaned.replace(multiDashRegex, "-");
  const lowercase = collapsed.toLowerCase();
  return lowercase;
}

export function isSlugReady(slug: string): boolean {
  const hasNoTrailingDash = !trailingDashRegex.test(slug);
  const matches = validSlugRegex.test(slug);
  const isValid = matches && hasNoTrailingDash;
  return isValid;
}
