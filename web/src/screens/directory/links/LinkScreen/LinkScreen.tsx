import { LinkListOKResponse, LinkListParams } from "src/api/openapi/schemas";
import { server } from "src/api/server";

import { Client } from "./Client";
import { Props } from "./useLinkScreen";

export async function LinkIndexScreen(props: Props) {
  const response = await server<LinkListOKResponse>({
    url: "/v1/links",
    params: {} as LinkListParams,
  });

  return <Client {...props} links={response.links} />;
}
