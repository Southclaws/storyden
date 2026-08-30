"use client";

import { usePathname } from "next/navigation";
import { useEffect, useRef } from "react";

export function ThemeRuntime() {
  const pathname = usePathname();
  const ready = useRef(false);

  useEffect(() => {
    if (!ready.current) {
      ready.current = true;
      document.dispatchEvent(new CustomEvent("storyden:ready"));
    }

    document.dispatchEvent(
      new CustomEvent("storyden:navigate", {
        detail: { pathname },
      }),
    );
  }, [pathname]);

  return null;
}
