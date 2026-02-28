import { UnreadyBanner } from "src/components/site/Unready";
import { ProfileScreen } from "src/screens/profile/ProfileScreen";

import { profileGet } from "@/api/openapi-server/profiles";
import { getServerSession } from "@/auth/server-session";
import { getSettings } from "@/lib/settings/settings-server";

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
