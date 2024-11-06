"use server";

import { cookies } from "next/headers";

import { NAVIGATION_SIDEBAR_STATE_KEY } from "@/local/state-keys";

import { parseSidebarCookie } from "./shared";

export async function getServerSidebarState() {
  const serverSidebarCookieState = (await cookies()).get(
    NAVIGATION_SIDEBAR_STATE_KEY,
  );

  const initialSidebarState = parseSidebarCookie(
    serverSidebarCookieState?.value,
  );

  return initialSidebarState;
}
