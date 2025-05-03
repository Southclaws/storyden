import { RegisterScreen } from "src/screens/auth/RegisterScreen/RegisterScreen";

import { OAuthProviderList } from "@/components/auth/OAuthProviderList";
import { UnreadyBanner } from "@/components/site/Unready";
import { getProviders } from "@/lib/auth/providers";
import { getSettings } from "@/lib/settings/settings-server";

// NOTE: We don't want any caching for data fetching here. OAuth URLs need to be
// generated freshly for each page render.
export const dynamic = "force-dynamic";

export default async function Page() {
  try {
    const { oauth } = await getProviders();

    return (
      <>
        <RegisterScreen />
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
    title: `Join the community at ${settings.title}`,
    description: `Log in or sign up to ${settings.title} - powered by Storyden`,
  };
}
