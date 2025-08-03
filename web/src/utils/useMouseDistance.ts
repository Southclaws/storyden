import { useEffect, useRef } from "react";

export function useMouseDistance<T extends HTMLElement = HTMLElement>() {
  const elementRef = useRef<T>(null);
  const distanceRef = useRef<{ x: number; y: number; d: number }>({
    x: 0,
    y: 0,
    d: Infinity,
  });

  useEffect(() => {
    const update = (e: MouseEvent) => {
      const el = elementRef.current;
      if (!el) return;

      const rect = el.getBoundingClientRect();
      const x = e.clientX;
      const y = e.clientY;

      const cx = Math.max(rect.left, Math.min(x, rect.right));
      const cy = Math.max(rect.top, Math.min(y, rect.bottom));

      const dx = x - cx;
      const dy = y - cy;
      const d = Math.sqrt(dx * dx + dy * dy);

      distanceRef.current = { x: dx, y: dy, d };
    };

    window.addEventListener("mousemove", update);
    return () => window.removeEventListener("mousemove", update);
  }, []);

  return { elementRef, distanceRef };
}
