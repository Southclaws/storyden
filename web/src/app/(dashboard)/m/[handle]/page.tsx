import { UnreadyBanner } from "src/components/site/Unready";
import { ProfileScreen } from "src/screens/profile/ProfileScreen";

import { profileGet } from "@/api/openapi-server/profiles";

type Props = {
  params: { handle: string };
};

export default async function Page({ params }: Props) {
  try {
    const { handle } = params;
    const { data } = await profileGet(handle);
    return <ProfileScreen profile={data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
