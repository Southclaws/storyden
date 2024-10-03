import { linkList } from "@/api/openapi-server/links";

import { Client } from "./Client";
import { Props } from "./useLinkIndexScreen";

export async function LinkIndexScreen(props: Omit<Props, "links">) {
  const { data } = await linkList({
    q: props.query,
    page: props.page?.toString(),
  });

  return <Client {...props} links={data} />;
}
