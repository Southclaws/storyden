import { ThreadGetResponse } from "src/api/openapi/schemas";
import { server } from "src/api/server";

import { Client } from "./Client";

type Props = {
  slug: string;
};

export async function ThreadScreen(props: Props) {
  const data = await server<ThreadGetResponse>({
    url: `/v1/threads/${props.slug}`,
    params: {
      slug: [props.slug],
    },
  });

  return <Client slug={props.slug} thread={data} />;
}
