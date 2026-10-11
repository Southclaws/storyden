export interface ContentOutlineEntry {
  key: string;
  kind: "section" | "heading" | "reply";
  label: string;
  depth: number;
  description?: string;
  createdAt?: string;
  target: HTMLElement;
}

export type CollectOutlineEntries = (
  root: HTMLElement,
) => ContentOutlineEntry[];

/** Read only the content root supplied by the screen, never the whole page. */
export function collectContentHeadings(root: HTMLElement, prefix: string) {
  const headings = root.querySelectorAll<HTMLElement>("h1, h2, h3, h4, h5, h6");
  const entries: ContentOutlineEntry[] = [];
  const levels: number[] = [];

  for (const target of headings) {
    const label = target.textContent?.trim();

    if (!label) {
      continue;
    }

    const visible = target.getClientRects().length > 0;

    if (!visible) {
      continue;
    }

    const embeddedContainer = target.closest(
      '[data-no-typography], [data-display="card"]',
    );

    if (embeddedContainer) {
      continue;
    }

    const level = Number(target.tagName.slice(1));

    while (levels.length > 0 && levels[levels.length - 1]! >= level) {
      levels.pop();
    }

    levels.push(level);

    entries.push({
      key: `${prefix}:heading:${entries.length}`,
      kind: "heading",
      label,
      depth: levels.length,
      target,
    });
  }

  return entries;
}

export function getOutlineReadingLine() {
  const header = document.querySelector(".navigation__topbar");
  const headerBottom = header?.getBoundingClientRect().bottom ?? 48;

  return Math.max(0, headerBottom) + 24;
}

export function getActiveOutlineKey(
  entries: ContentOutlineEntry[],
  readingLine: number,
) {
  const viewportBottom = window.scrollY + window.innerHeight;
  const documentBottom = document.documentElement.scrollHeight;
  const atBottom = viewportBottom >= documentBottom - 2;
  const hasScrolled = window.scrollY > 0;

  if (hasScrolled && atBottom) {
    return entries.at(-1)?.key;
  }

  let active = entries[0]?.key;

  for (const entry of entries) {
    const top = entry.target.getBoundingClientRect().top;

    // Smooth scrolling can settle on a fractional CSS pixel below the target.
    if (top > readingLine + 2) {
      break;
    }

    active = entry.key;
  }

  return active;
}

export function scrollToOutlineEntry(entry: ContentOutlineEntry) {
  const { target } = entry;

  if (!target.isConnected) {
    return;
  }

  const targetTop = target.getBoundingClientRect().top + window.scrollY;
  const top = targetTop - getOutlineReadingLine();
  const reducedMotion = window.matchMedia(
    "(prefers-reduced-motion: reduce)",
  ).matches;

  // A temporary focus target does not change the stored editor document.
  if (!target.hasAttribute("tabindex")) {
    target.setAttribute("tabindex", "-1");
    target.addEventListener("blur", () => target.removeAttribute("tabindex"), {
      once: true,
    });
  }

  target.focus({ preventScroll: true });

  window.scrollTo({
    top,
    behavior: reducedMotion ? "instant" : "smooth",
  });
}
