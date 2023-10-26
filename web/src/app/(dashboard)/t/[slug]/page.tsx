import { server } from "src/api/client";
import { ThreadGetResponse } from "src/api/openapi/schemas";
import { ThreadScreen } from "src/screens/thread/ThreadScreen";
import { getInfo } from "src/utils/info";

export type Props = {
  params: {
    slug: string;
  };
};

export default function Page(props: Props) {
  return <ThreadScreen slug={props.params.slug} />;
}

export async function generateMetadata({ params }: Props) {
  const info = await getInfo();
  const data = await server<ThreadGetResponse>({
    url: `/v1/threads/${params.slug}`,
    params: {
      slug: [params.slug],
    },
  });

  return {
    title: `${data.title} | ${info.title}`,
    description: data.short,
  };
}
