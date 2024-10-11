import Cookies from "js-cookie";
import { useState } from "react";

import { NAVIGATION_SIDEBAR_STATE_KEY } from "@/local/state-keys";

import { parseSidebarCookie } from "./shared";

export function useSidebar(initialValue: boolean) {
  const clientSideCookieValue = Cookies.get(NAVIGATION_SIDEBAR_STATE_KEY);
  const initialState =
    initialValue ?? parseSidebarCookie(clientSideCookieValue);

  const [showLeftBar, setLocalState] = useState(initialState);

  if (typeof window === "undefined") {
    return {
      showLeftBar: initialValue,
      setShowLeftBar: () => {},
    };
  }

  function setShowLeftBar() {
    const next = !showLeftBar;

    Cookies.set(NAVIGATION_SIDEBAR_STATE_KEY, next ? "true" : "false", {
      secure: true,
      sameSite: "lax",
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
