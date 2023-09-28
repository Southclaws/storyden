import { server } from "src/api/client";
import { ThreadListOKResponse } from "src/api/openapi/schemas";

import { Client } from "./Client";

export async function HomeScreen() {
  const data = await server<ThreadListOKResponse>({ url: "/v1/threads" });

  return (
    <>
      <Client threads={data.threads} />
    </>
  );
}
