import {
  ProfileListOKResponse,
  ProfileListParams,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";

import { Client } from "./Client";
import { Props } from "./useMemberIndexScreen";

export async function MemberIndexScreen(props: Omit<Props, "profiles">) {
  const response = await server<ProfileListOKResponse>({
    url: "/v1/profiles",
    params: {
      q: props.query,
      page: props.page,
    } as ProfileListParams,
  });

  return <Client {...props} profiles={response} />;
}
