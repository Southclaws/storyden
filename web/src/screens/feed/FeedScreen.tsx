import {
  ThreadListOKResponse,
  ThreadListParams,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";

import { Client } from "./Client";

type Props = {
  category: string;
};

export async function FeedScreen(props: Props) {
  const data = await server<ThreadListOKResponse>({
    url: `/v1/threads`,
    params: {
      categories: [props.category],
    } as ThreadListParams,
  });

  return <Client category={props.category} threads={data.threads} />;
}
