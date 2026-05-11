import { LoginScreen } from "src/screens/auth/LoginScreen/LoginScreen";

import { OAuthProviderList } from "@/components/auth/OAuthProviderList";
import { UnreadyBanner } from "@/components/site/Unready";
import { tServer } from "@/i18n/server";
import { getProviders } from "@/lib/auth/providers";
import { getSettings } from "@/lib/settings/settings-server";

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
  const login = await tServer("Login");
  const register = await tServer("Register");

  return {
    title: `${login} ${settings.title}`,
    description: `${login} / ${register} ${settings.title} - Storyden`,
  };
}
