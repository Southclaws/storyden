import { server } from "src/api/client";
import { ThreadListOKResponse } from "src/api/openapi/schemas";

import { Client } from "./Client";
import { Fallback } from "./Fallback";

export async function Content() {
  if (typeof window === "undefined") {
    const data = await server<ThreadListOKResponse>("v1/threads");

    return <Fallback threads={data.threads} />;
  }

  return <Client showEmptyState />;
}
