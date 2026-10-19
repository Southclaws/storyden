import { type RefObject, useLayoutEffect, useState } from "react";
import { createPortal } from "react-dom";

import { type dragTree } from "@/styled-system/recipes";

/** A visual companion to the tree's live status; it never takes keyboard focus. */
export function DragTreeMovePreview({
  rootRef,
  focusedId,
  label,
  preview,
  pending,
  styles,
}: {
  rootRef: RefObject<HTMLDivElement | null>;
  focusedId: string | null;
  label: string;
  preview: string;
  pending: boolean;
  styles: ReturnType<typeof dragTree>;
}) {
  const [position, setPosition] = useState({ x: 0, y: 0, width: 240 });
  useLayoutEffect(() => {
    function update() {
      const anchor = Array.from(
        rootRef.current?.querySelectorAll<HTMLAnchorElement>(
          "a[data-tree-node-id]",
        ) ?? [],
      ).find((element) => element.dataset["treeNodeId"] === focusedId);
      const row = anchor?.parentElement;
      if (!row) return;
      const rect = row.getBoundingClientRect();
      const width = Math.min(280, window.innerWidth - 24);
      const fitsRight = rect.right + width + 24 <= window.innerWidth;
      const y = fitsRight
        ? rect.top
        : rect.bottom + 128 <= window.innerHeight
          ? rect.bottom + 8
          : rect.top - 120;
      setPosition({
        x: Math.max(
          12,
          fitsRight
            ? rect.right + 12
            : Math.min(rect.left, window.innerWidth - width - 12),
        ),
        y: Math.max(12, Math.min(y, window.innerHeight - 120)),
        width,
      });
    }
    update();
    window.addEventListener("scroll", update, true);
    window.addEventListener("resize", update);
    return () => {
      window.removeEventListener("scroll", update, true);
      window.removeEventListener("resize", update);
    };
  }, [focusedId, rootRef]);

  return createPortal(
    <div
      aria-hidden="true"
      className={styles.movePreview}
      style={{
        width: position.width,
        translate: `${position.x}px ${position.y}px`,
      }}
    >
      <strong className={styles.moveTitle}>{label}</strong>
      <span>{pending ? "Placing…" : preview}</span>
      <span className={styles.moveHelp}>
        Arrows: browse · Space: placement · Enter: place · Esc: cancel
      </span>
    </div>,
    document.body,
  );
}
