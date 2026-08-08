import { useEffect, useRef } from "react";

const OVERFLOW_THRESHOLD = 1;

export function useOverflowGradient<T extends HTMLElement>() {
  const ref = useRef<T>(null);

  useEffect(() => {
    const viewport = ref.current;
    if (!viewport) return;

    let frame: number | undefined;

    const update = () => {
      if (frame !== undefined) cancelAnimationFrame(frame);

      frame = requestAnimationFrame(() => {
        const remaining =
          viewport.scrollHeight - viewport.clientHeight - viewport.scrollTop;

        viewport.dataset["overflowBlockEnd"] = String(
          remaining > OVERFLOW_THRESHOLD,
        );
      });
    };

    const resizeObserver = new ResizeObserver(update);
    const observeLayout = () => {
      resizeObserver.disconnect();
      resizeObserver.observe(viewport);

      for (const child of viewport.children) {
        resizeObserver.observe(child);
      }
    };

    const mutationObserver = new MutationObserver(() => {
      observeLayout();
      update();
    });

    observeLayout();
    mutationObserver.observe(viewport, {
      childList: true,
      subtree: true,
    });
    viewport.addEventListener("scroll", update, { passive: true });
    update();

    return () => {
      if (frame !== undefined) cancelAnimationFrame(frame);
      resizeObserver.disconnect();
      mutationObserver.disconnect();
      viewport.removeEventListener("scroll", update);
      delete viewport.dataset["overflowBlockEnd"];
    };
  }, []);

  return ref;
}
