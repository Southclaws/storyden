import { profileList } from "@/api/openapi-server/profiles";

import { Client } from "./Client";
import { Props } from "./useMemberIndexScreen";

export async function MemberIndexScreen(props: Omit<Props, "profiles">) {
  const { data } = await profileList({
    q: props.query,
    page: props.page?.toString(),
  });

  return <Client {...props} profiles={data} />;
}
