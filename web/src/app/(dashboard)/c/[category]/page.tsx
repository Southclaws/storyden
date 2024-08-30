import { Unready } from "src/components/site/Unready";
import { FeedScreenClient } from "src/screens/feed/FeedScreenClient";

import { threadList } from "@/api/openapi-server/threads";

type Props = {
  params: {
    category: string;
  };
};

export default async function Page(props: Props) {
  try {
    // NOTE: This is very unoptimised but GSD. Long term we want an actual API
    // for feeds at /feed which delivers a list of DatagraphNodeReference
    // objects of all kinds based on a set of heuristics such as what's hot,
    // what's relevant to the account (if any) and what's been featured.

    const { data } = await threadList({
      categories: [props.params.category],
    });

    return (
      <FeedScreenClient
        params={{
          threads: {
            categories: [props.params.category],
          },
        }}
        initialData={{ threads: data }}
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
