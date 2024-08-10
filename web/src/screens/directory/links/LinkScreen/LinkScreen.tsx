import { linkGet } from "@/api/openapi-server/links";

import { Client } from "./Client";
import { Props } from "./useLinkScreen";

export async function LinkScreen(props: Omit<Props, "link">) {
  const { data } = await linkGet(props.slug);

  return <Client {...props} link={data} />;
}
