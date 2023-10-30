import { ThreadListOKResponse } from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { Unready } from "src/components/site/Unready";

import { Client } from "./Client";

export async function HomeScreen() {
  try {
    const data = await server<ThreadListOKResponse>({ url: "/v1/threads" });

    return (
      <>
        <Client threads={data.threads} />
      </>
    );
  } catch (error) {
    return (
      <Unready
        message={"Content failed to load"}
        error={(error as Error).message}
        metadata={JSON.parse(JSON.stringify(error))}
      />
    );
  }
}
