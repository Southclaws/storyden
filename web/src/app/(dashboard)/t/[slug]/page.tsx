import { ThreadScreen } from "src/screens/thread/ThreadScreen";
import { getInfo } from "src/utils/info";

import { threadGet } from "@/api/openapi-server/threads";

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
  const { data } = await threadGet(params.slug);

  return {
    title: `${data.title} | ${info.title}`,
    description: data.description,
  };
}
