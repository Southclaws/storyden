import { UnreadyBanner } from "src/components/site/Unready";

import { threadList } from "@/api/openapi-server/threads";
import { FeedScreen } from "@/screens/feed/FeedScreen";

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
      <FeedScreen
        params={{
          categories: [props.params.category],
        }}
        initialData={data}
      />
    );
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
