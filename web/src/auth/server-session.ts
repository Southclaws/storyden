"use server";

import { cookies } from "next/headers";

import { accountGet } from "@/api/openapi-server/accounts";

export async function getServerSession() {
  const session = cookies().get("storyden-session");

  if (!session) return;

  try {
    const { data } = await accountGet();

    return data;
  } catch (e) {
    console.error(e);
    return;
  }
}
