import { z } from "zod";

import { SidebarDefaultState } from "@/lib/settings/sidebar";

const CookieSchema = z.string().transform((value) => value === "true");

export function parseSidebarCookie(
  cookieValue?: string,
  defaultState: SidebarDefaultState = "closed",
) {
  const { success, data } = CookieSchema.safeParse(cookieValue);

  if (!success) {
    return defaultState === "open";
  }

  return data;
}
