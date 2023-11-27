import { LinkGetOKResponse } from "src/api/openapi/schemas";
import { server } from "src/api/server";

import { Client } from "./Client";
import { Props } from "./useLinkScreen";

export async function LinkScreen(props: Omit<Props, "link">) {
  const response = await server<LinkGetOKResponse>({
    url: `/v1/links/${props.slug}`,
  });

  return <Client {...props} link={response} />;
}
