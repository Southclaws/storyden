// The initial state of contenteditable or content.Content html parse of empty.
const EMPTY_CONTENT_INITIAL = "<body></body>";

// Sometimes pops up also, usually if the user backspaced every character.
const EMPTY_CONTENT_MANUAL = "<body><p></p></body>";

// TipTap emits this when an editable rich-text field is cleared.
const EMPTY_CONTENT_RICH = "<p></p>";

export function isContentEmpty(s: string | undefined) {
  if (s === undefined) {
    return true;
  }

  const trimmed = s.trim();

  if (trimmed === "") {
    return true;
  }

  if (trimmed === EMPTY_CONTENT_INITIAL) {
    return true;
  }

  if (trimmed === EMPTY_CONTENT_MANUAL) {
    return true;
  }

  if (trimmed === EMPTY_CONTENT_RICH) {
    return true;
  }

  return false;
}
