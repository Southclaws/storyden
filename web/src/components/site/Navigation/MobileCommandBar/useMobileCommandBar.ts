"use client";

import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";

import { useSession } from "@/auth";

export function useMobileCommandBar() {
  const pathname = usePathname();
  const [isExpanded, setExpanded] = useState(false);
  const account = useSession();

  // Close the menu for either navigation events or outside clicks/taps:

  useEffect(() => setExpanded(false), [pathname]);

  function onExpand() {
    setExpanded(true);
  }

  function onClose() {
    setExpanded(false);
  }

  return {
    isExpanded,
    onExpand,
    onClose,
    account,
  };
}
