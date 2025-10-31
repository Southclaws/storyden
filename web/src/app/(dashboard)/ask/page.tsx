import { redirect } from "next/navigation";

import { getServerSession } from "@/auth/server-session";
import { AskScreen } from "@/screens/ask/AskScreen";

export default async function Page() {
  const session = await getServerSession();

  return <AskScreen session={session} />;
}
