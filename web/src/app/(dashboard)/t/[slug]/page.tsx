import { getInfo } from "src/utils/info";

import { threadGet } from "@/api/openapi-server/threads";
import { UnreadyBanner } from "@/components/site/Unready";
import { ThreadScreen } from "@/screens/thread/ThreadScreen/ThreadScreen";

export type Props = {
  params: {
    slug: string;
  };
};

export default async function Page(props: Props) {
  const { slug } = props.params;

  try {
    const { data } = await threadGet(slug);

    return <ThreadScreen slug={slug} thread={data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}

export async function generateMetadata({ params }: Props) {
  try {
    const info = await getInfo();
    const { data } = await threadGet(params.slug);

    return {
      title: `${data.title} | ${info.title}`,
      description: data.description,
    };
  } catch (e) {
    return {
      title: "Thread Not Found",
      description: "The thread you are looking for does not exist.",
    };
  }
}
