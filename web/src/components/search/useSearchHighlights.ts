"use client";

import { escapeRegExp, uniqBy } from "lodash";
import { useCallback } from "react";

const highlightName = "storyden-search-match";
const highlights = new WeakMap<Document, Highlight>();
const excludedText =
  'input, textarea, select, script, style, [hidden], [aria-hidden="true"], [contenteditable]:not([contenteditable="false"])';

function findSearchRanges(root: HTMLElement, pattern: RegExp): Range[] {
  const ranges: Range[] = [];

  for (const target of root.querySelectorAll("[data-search-highlight]")) {
    const walker = root.ownerDocument.createTreeWalker(
      target,
      NodeFilter.SHOW_TEXT,
    );

    while (walker.nextNode()) {
      const node = walker.currentNode;
      if (node.parentElement?.closest(excludedText)) continue;

      for (const match of (node.textContent ?? "").matchAll(pattern)) {
        const range = root.ownerDocument.createRange();
        range.setStart(node, match.index);
        range.setEnd(node, match.index + match[0].length);
        ranges.push(range);
      }
    }
  }

  return ranges;
}

/** Paint matches in marked descendants without changing their content or layout. */
export function useSearchHighlights(query: string, enabled = true) {
  return useCallback(
    (root: HTMLElement | null) => {
      if (
        !root ||
        !enabled ||
        !query.trim() ||
        typeof Highlight === "undefined" ||
        typeof CSS === "undefined" ||
        !CSS.highlights ||
        !CSS.supports?.(`selector(::highlight(${highlightName}))`)
      ) {
        return;
      }

      const terms = uniqBy(query.trim().split(/\s+/u), (term) =>
        term.toLowerCase(),
      );
      const pattern = new RegExp(
        terms
          .sort((a, b) => b.length - a.length)
          .map(escapeRegExp)
          .join("|"),
        "giu",
      );
      const document = root.ownerDocument;
      const container = root;
      const highlight = highlights.get(document) ?? new Highlight();
      highlights.set(document, highlight);
      let ranges: Range[] = [];
      let frame: number | undefined;

      function removeRanges() {
        for (const range of ranges) highlight.delete(range);
        if (highlight.size === 0) CSS.highlights.delete(highlightName);
      }

      function update() {
        frame = undefined;
        removeRanges();
        ranges = findSearchRanges(container, pattern);
        for (const range of ranges) highlight.add(range);
        if (highlight.size > 0) CSS.highlights.set(highlightName, highlight);
      }

      update();
      const observer = new MutationObserver(() => {
        frame ??= requestAnimationFrame(update);
      });
      observer.observe(root, {
        childList: true,
        characterData: true,
        subtree: true,
      });

      return () => {
        observer.disconnect();
        if (frame !== undefined) cancelAnimationFrame(frame);
        removeRanges();
      };
    },
    [query, enabled],
  );
}
