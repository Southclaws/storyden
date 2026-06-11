import { useLayoutEffect, useRef, useState } from "react";

export function useIsTextWrapping<T extends HTMLElement>() {
  const ref = useRef<T>(null);
  const [wrapped, setWrapped] = useState(false);

  useLayoutEffect(() => {
    const el = ref.current;
    if (!el) return;

    let frame: number | null = null;

    const check = () => {
      if (frame !== null) cancelAnimationFrame(frame);

      frame = requestAnimationFrame(() => {
        const range = document.createRange();
        range.selectNodeContents(el);

        setWrapped(range.getClientRects().length > 1);
      });
    };

    check();

    const ro = new ResizeObserver(check);

    // Usually the parent/container width is what changes
    if (el.parentElement) {
      ro.observe(el.parentElement);
    } else {
      ro.observe(el);
    }

    return () => {
      if (frame !== null) cancelAnimationFrame(frame);
      ro.disconnect();
    };
  }, []);

  return [ref, wrapped] as const;
}
