import { cookies } from "next/headers";

import { AccountGetOKResponse } from "src/api/openapi/schemas";
import { server } from "src/api/server";

export async function useServerSession() {
  const session = cookies().get("storyden-session");
  const cookie = session ? `${session.name}=${session.value}` : "";

  if (!session) return;

  // NOTE: Throws when no session
  const data = await server<AccountGetOKResponse>({
    url: `/v1/accounts`,
    cookie,
  });

  return data;
}
