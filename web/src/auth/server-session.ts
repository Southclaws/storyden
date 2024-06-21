"use server";

import { cookies } from "next/headers";

import { AccountGetOKResponse } from "src/api/openapi/schemas";
import { server } from "src/api/server";

export async function getServerSession() {
  const session = cookies().get("storyden-session");
  const cookie = session ? `${session.name}=${session.value}` : "";

  if (!session) return;

  try {

    const data = await server<AccountGetOKResponse>({
      url: `/v1/accounts`,
      cookie,
    });
    
    return data;
  } catch(e) {
    return
  }
}
