import { z } from "zod";

const DEFAULT_SIDEBAR_STATE = false;

const CookieSchema = z.string().transform((value) => value === "true");

export function parseSidebarCookie(cookieValue?: string) {
  const { success, data } = CookieSchema.safeParse(cookieValue);

  if (!success) {
    return DEFAULT_SIDEBAR_STATE;
  }

  return data;
}
