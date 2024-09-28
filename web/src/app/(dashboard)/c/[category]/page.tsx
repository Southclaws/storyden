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
    const { data } = await threadList({
      categories: [props.params.category],
    });

    return (
      <FeedScreenClient
        params={{
          categories: [props.params.category],
        }}
        initialData={data}
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
