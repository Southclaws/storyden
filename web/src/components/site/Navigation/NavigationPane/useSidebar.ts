import { useState } from "react";

import { NAVIGATION_SIDEBAR_STATE_KEY } from "@/local/state-keys";
import { SidebarDefaultState } from "@/lib/settings/sidebar";
import { getCookie, setCookie } from "@/utils/cookie";

import { parseSidebarCookie } from "./shared";

export function useSidebar(
  initialValue: boolean,
  defaultState: SidebarDefaultState = "closed",
) {
  const clientSideCookieValue = getCookie(NAVIGATION_SIDEBAR_STATE_KEY);
  const initialState =
    initialValue ?? parseSidebarCookie(clientSideCookieValue, defaultState);

  const [showLeftBar, setLocalState] = useState(initialState);

  if (typeof window === "undefined") {
    return {
      showLeftBar: initialValue,
      setShowLeftBar: () => {},
    };
  }

  function setShowLeftBar() {
    const next = !showLeftBar;

    setCookie(NAVIGATION_SIDEBAR_STATE_KEY, next ? "true" : "false", {
      days: 180,
    });
    setLocalState(next);

    // Manipulate the DOM directly to show/hide the left bar.
    // This is done because the component itself is a server-side component
    // and cannot be re-rendered on the client side. There's no simple way to
    // have this component "listen" for cookie state changes.
    document
      .querySelector("#navigation__container")
      ?.setAttribute("data-leftbar-shown", next.toString());
  }

  return {
    showLeftBar,
    setShowLeftBar,
  };
}
