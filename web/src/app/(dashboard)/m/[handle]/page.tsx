import { UnreadyBanner } from "@/components/site/Unready";
import { ProfileScreen } from "@/screens/profile/ProfileScreen";

import { profileGet } from "@/api/openapi-server/profiles";
import { getServerSession } from "@/auth/server-session";
import { getSettings } from "@/lib/settings/settings-server";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

type Props = {
  params: Promise<{ handle: string }>;
};

export default async function Page(props: Props) {
  const params = await props.params;
  try {
    const { handle } = params;

    const session = await getServerSession();
    const settings = await getSettings();
    const { data } = await profileGet(handle);

    return (
      <ProfileScreen
        initialSession={session}
        profile={data}
        initialSignatureConfig={settings.metadata.signatures}
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
