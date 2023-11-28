import { LinkListOKResponse, LinkListParams } from "src/api/openapi/schemas";
import { server } from "src/api/server";

import { Client } from "./Client";
import { Props } from "./useLinkIndexScreen";

export async function LinkIndexScreen(props: Omit<Props, "links">) {
  const response = await server<LinkListOKResponse>({
    url: "/v1/links",
    params: {
      q: props.query,
      page: props.page,
    } as LinkListParams,
  });

  return <Client {...props} links={response.links} />;
}
