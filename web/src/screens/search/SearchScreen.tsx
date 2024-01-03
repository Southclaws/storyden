import { cookies } from "next/headers";

import {
  DatagraphSearchOKResponse,
  DatagraphSearchParams,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";

import { Client } from "./Client";

type Props = {
  query: string;
};

export async function SearchScreen(props: Props) {
  const session = cookies().get("storyden-session");
  const cookie = session ? `${session.name}=${session.value}` : "";

  const data = await server<DatagraphSearchOKResponse>({
    url: `/v1/datagraph`,
    params: {
      q: props.query,
    } as DatagraphSearchParams,
    cookie,
  });

  return <Client query={props.query} results={data} />;
}
