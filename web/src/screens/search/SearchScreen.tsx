import { datagraphSearch } from "@/api/openapi-server/datagraph";

import { Client } from "./Client";

type Props = {
  query: string;
};

export async function SearchScreen(props: Props) {
  const { data } = await datagraphSearch({ q: props.query });

  return <Client query={props.query} results={data} />;
}
