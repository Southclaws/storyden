import { Unready } from "src/components/site/Unready";
import { FeedScreenClient } from "src/screens/feed/FeedScreenClient";

import { nodeList } from "@/api/openapi-server/nodes";
import { threadList } from "@/api/openapi-server/threads";

export default async function Page() {
  try {
    // NOTE: This is very unoptimised but GSD. Long term we want an actual API
    // for feeds at /v1/feed which delivers a list of DatagraphNodeReference
    // objects of all kinds based on a set of heuristics such as what's hot,
    // what's relevant to the account (if any) and what's been featured.

    const [threads, nodes] = await Promise.all([threadList(), nodeList()]);

    return (
      <FeedScreenClient
        initialData={{
          threads: threads.data,
          nodes: nodes.data,
        }}
      />
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
