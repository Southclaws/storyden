"use server";

import { cookies } from "next/headers";
import { cache } from "react";

import { RequestError } from "@/api/common";
import { accountGet } from "@/api/openapi-server/accounts";

const getSessionCached = cache(async () => {
  return await accountGet({
    cache: "default",
  });
});

export async function getServerSession(options?: RequestInit) {
  const session = (await cookies()).get("storyden-session");

  if (!session) return;

  try {
    const { data } = options
      ? await accountGet(options)
      : await getSessionCached();

    return data;
  } catch (e) {
    if (e instanceof RequestError) {
      if (e.status === 401) {
        console.debug("user not authenticated:", e);
        return;
      } else if (e.status === 403) {
        console.debug("user not authorised:", e);
        return;
      }
    }

    console.error("get server session failed:", e);
    return;
  }
}
