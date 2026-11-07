import type { KeyboardEvent } from "react";

export function submitFormOnModEnter(event: KeyboardEvent<HTMLElement>) {
  const isSubmitShortcut =
    event.key === "Enter" &&
    (event.metaKey || event.ctrlKey) &&
    !event.shiftKey &&
    !event.altKey;

  if (!isSubmitShortcut || event.defaultPrevented) {
    return;
  }

  if (event.nativeEvent.isComposing || event.nativeEvent.keyCode === 229) {
    return;
  }

  const target = event.target;
  if (!(target instanceof HTMLElement)) {
    return;
  }

  if (!target.closest("textarea, [contenteditable='true']")) {
    return;
  }

  const form = target.closest("form");
  if (!form) {
    return;
  }

  const submitter = Array.from(form.elements).find(
    (element): element is HTMLButtonElement | HTMLInputElement =>
      (element instanceof HTMLButtonElement ||
        element instanceof HTMLInputElement) &&
      element.type === "submit" &&
      !element.matches(":disabled"),
  );

  if (!submitter) {
    return;
  }

  // Capture before TipTap handles Mod+Enter as an editor command.
  event.preventDefault();
  event.stopPropagation();

  if (!event.repeat) {
    // Clicking preserves validation, submitter values and disabled states.
    submitter.click();
  }
}
