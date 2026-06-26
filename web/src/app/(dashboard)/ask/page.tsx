import { getServerSession } from "@/auth/server-session";
import { getSettings } from "@/lib/settings/settings-server";
import { AskScreen } from "@/screens/ask/AskScreen";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export default async function Page() {
  const [session, settings] = await Promise.all([
    getServerSession(),
    getSettings(),
  ]);

  return <AskScreen session={session} initialSettings={settings} />;
}
