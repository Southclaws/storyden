import {
  LinkListOKResponse,
  NodeListOKResponse,
  ThreadListOKResponse,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { Unready } from "src/components/site/Unready";
import { FeedScreenClient } from "src/screens/feed/FeedScreenClient";

export default async function Page() {
  try {
    // NOTE: This is very unoptimised but GSD. Long term we want an actual API
    // for feeds at /v1/feed which delivers a list of DatagraphNodeReference
    // objects of all kinds based on a set of heuristics such as what's hot,
    // what's relevant to the account (if any) and what's been featured.

    const [threads, nodes, links] = await Promise.all([
      server<ThreadListOKResponse>({ url: "/v1/threads" }),
      server<NodeListOKResponse>({ url: "/v1/nodes" }),
      server<LinkListOKResponse>({ url: "/v1/links" }),
    ]);

    return (
      <FeedScreenClient
        initialData={{
          threads,
          nodes,
          links,
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
