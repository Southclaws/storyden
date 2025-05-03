import { LoginScreen } from "src/screens/auth/LoginScreen/LoginScreen";

import { OAuthProviderList } from "@/components/auth/OAuthProviderList";
import { UnreadyBanner } from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";

// NOTE: We don't want any caching for data fetching here. OAuth URLs need to be
// generated freshly for each page render.
export const dynamic = "force-dynamic";

export default function Page() {
  try {
    return (
      <>
        <LoginScreen />
        <OAuthProviderList />
      </>
    );
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}

export async function generateMetadata() {
  const settings = await getSettings();
  return {
    title: `Login to ${settings.title}`,
    description: `Log in or sign up to ${settings.title} - powered by Storyden`,
  };
}
