import { redirect } from "next/navigation";

import { getServerSession } from "@/auth/server-session";
import { AskScreen } from "@/screens/ask/AskScreen";

export const dynamic = "force-dynamic";

export default async function Page() {
  const session = await getServerSession();

  return <AskScreen session={session} />;
}
