import { threadGet } from "@/api/openapi-server/threads";

import { Client } from "./Client";

type Props = {
  slug: string;
};

export async function ThreadScreen(props: Props) {
  const { data } = await threadGet(props.slug);

  return <Client slug={props.slug} thread={data} />;
}
