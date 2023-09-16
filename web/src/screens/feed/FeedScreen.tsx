import { server } from "src/api/client";
import { ThreadListOKResponse } from "src/api/openapi/schemas";
import { getThreadListKey } from "src/api/openapi/threads";

import { Client } from "./Client";

type Props = {
  category: string;
};

export async function FeedScreen(props: Props) {
  const key = getThreadListKey({ categories: [props.category] })[0];

  const data = await server<ThreadListOKResponse>(key);

  return <Client category={props.category} threads={data.threads} />;
}
