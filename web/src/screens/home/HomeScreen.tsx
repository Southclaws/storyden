import { server } from "src/api/client";
import { ThreadListOKResponse } from "src/api/openapi/schemas";
import { getThreadListKey } from "src/api/openapi/threads";

import { Client } from "./Client";

export async function HomeScreen() {
  const key = getThreadListKey()[0];

  const data = await server<ThreadListOKResponse>(key);

  return (
    <>
      <Client threads={data.threads} />
    </>
  );
}
