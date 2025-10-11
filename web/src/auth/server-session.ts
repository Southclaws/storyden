"use server";

import { cookies } from "next/headers";

import { RequestError } from "@/api/common";
import { accountGet } from "@/api/openapi-server/accounts";

export async function getServerSession() {
  const session = (await cookies()).get("storyden-session");

  if (!session) return;

  try {
    const { data } = await accountGet({
      next: { tags: ["accounts"] },
    });

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
