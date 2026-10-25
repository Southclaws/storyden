import { useEffect, useRef, useState } from "react";

import {
  type CollectOutlineEntries,
  type ContentOutlineEntry,
  getActiveOutlineKey,
  getOutlineReadingHeight,
  getOutlineReadingLine,
} from "./outline";

function sameEntries(
  previous: ContentOutlineEntry[],
  next: ContentOutlineEntry[],
) {
  if (previous.length !== next.length) {
    return false;
  }

  for (let index = 0; index < previous.length; index++) {
    const previousEntry = previous[index]!;
    const nextEntry = next[index]!;

    if (previousEntry.key !== nextEntry.key) {
      return false;
    }

    if (previousEntry.target !== nextEntry.target) {
      return false;
    }

    if (previousEntry.kind !== nextEntry.kind) {
      return false;
    }

    if (previousEntry.label !== nextEntry.label) {
      return false;
    }

    if (previousEntry.description !== nextEntry.description) {
      return false;
    }

    if (previousEntry.createdAt !== nextEntry.createdAt) {
      return false;
    }

    if (previousEntry.depth !== nextEntry.depth) {
      return false;
    }
  }

  return true;
}

export function useContentOutline(
  collectEntries: CollectOutlineEntries,
  disabled: boolean,
) {
  const rootRef = useRef<HTMLDivElement>(null);
  const [entries, setEntries] = useState<ContentOutlineEntry[]>([]);
  const [activeKey, setActiveKey] = useState<string>();
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const root = rootRef.current;

    if (!root) {
      return;
    }

    if (disabled) {
      return;
    }

    const desktop = window.matchMedia("(min-width: 1024px)");
    let frame = 0;

    const collect = () => {
      cancelAnimationFrame(frame);

      frame = requestAnimationFrame(() => {
        const next = desktop.matches ? collectEntries(root) : [];

        setVisible(
          next.length > 1 &&
            getOutlineReadingHeight(root, next) >= 2 * window.innerHeight,
        );

        setEntries((previous) => {
          if (sameEntries(previous, next)) {
            return previous;
          }

          return next;
        });

        setActiveKey(getActiveOutlineKey(next, getOutlineReadingLine()));
      });
    };

    const mutations = new MutationObserver(collect);
    const resize = new ResizeObserver(collect);

    const observe = () => {
      mutations.disconnect();
      resize.disconnect();

      if (desktop.matches) {
        mutations.observe(root, {
          childList: true,
          subtree: true,
          characterData: true,
          attributes: true,
          attributeFilter: [
            "data-outline-label",
            "data-outline-description",
            "data-outline-created-at",
          ],
        });

        resize.observe(root);
      }

      collect();
    };

    observe();

    desktop.addEventListener("change", observe);
    window.addEventListener("resize", collect);

    return () => {
      cancelAnimationFrame(frame);

      mutations.disconnect();
      resize.disconnect();

      desktop.removeEventListener("change", observe);
      window.removeEventListener("resize", collect);
    };
  }, [collectEntries, disabled]);

  useEffect(() => {
    if (disabled) {
      return;
    }

    if (!visible) {
      return;
    }

    let frame = 0;

    const update = () => {
      cancelAnimationFrame(frame);

      frame = requestAnimationFrame(() => {
        const readingLine = getOutlineReadingLine();
        const activeKey = getActiveOutlineKey(entries, readingLine);

        setActiveKey(activeKey);
      });
    };

    update();

    window.addEventListener("scroll", update, { passive: true });
    window.addEventListener("resize", update);

    return () => {
      cancelAnimationFrame(frame);

      window.removeEventListener("scroll", update);
      window.removeEventListener("resize", update);
    };
  }, [disabled, entries, visible]);

  return { rootRef, entries, activeKey, visible: !disabled && visible };
}
