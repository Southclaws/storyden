import { getServerSession } from "@/auth/server-session";

import { MemberActions } from "./MemberActions";

export async function MemberActionsServer() {
  const session = await getServerSession();

  return <MemberActions session={session} />;
}
