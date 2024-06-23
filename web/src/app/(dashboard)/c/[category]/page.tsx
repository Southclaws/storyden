import {
  LinkListOKResponse,
  NodeListOKResponse,
  ThreadListOKResponse,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { Unready } from "src/components/site/Unready";
import { FeedScreenClient } from "src/screens/feed/FeedScreenClient";

type Props = {
  params: {
    category: string;
  };
};

export default async function Page(props: Props) {
  try {
    // NOTE: This is very unoptimised but GSD. Long term we want an actual API
    // for feeds at /v1/feed which delivers a list of DatagraphNodeReference
    // objects of all kinds based on a set of heuristics such as what's hot,
    // what's relevant to the account (if any) and what's been featured.

    const threads = await server<ThreadListOKResponse>({
      url: "/v1/threads",
      params: {
        category: props.params.category,
      },
    });

    return (
      <FeedScreenClient
        params={{
          threads: {
            categories: [props.params.category],
          },
        }}
        initialData={{
          threads,
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
