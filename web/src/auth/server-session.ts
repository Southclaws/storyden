"use server";

import { cookies } from "next/headers";
import { cache } from "react";

import { RequestError } from "@/api/common";
import { getSession } from "@/api/openapi-server/misc";

const getSessionCached = cache(async (sessionCookie: string) => {
  void sessionCookie;

  const { data } = await getSession({
    cache: "no-store",
  });

  return data;
});

export async function getServerSession() {
  const session = (await cookies()).get("storyden-session");

  if (!session) return;

  try {
    const data = await getSessionCached(session.value);

    return data.account;
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
