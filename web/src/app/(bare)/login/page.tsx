import { LoginScreen } from "@/screens/auth/LoginScreen/LoginScreen";

import { OAuthProviderList } from "@/components/auth/OAuthProviderList";
import { UnreadyBanner } from "@/components/site/Unready";
import { getProviders } from "@/lib/auth/providers";
import { getSettings } from "@/lib/settings/settings-server";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export default async function Page() {
  try {
    const { oauth } = await getProviders();

    return (
      <>
        <LoginScreen />
        <OAuthProviderList providers={oauth} />
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
