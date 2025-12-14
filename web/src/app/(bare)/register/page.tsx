import { redirect } from "next/navigation";
import { RegisterScreen } from "src/screens/auth/RegisterScreen/RegisterScreen";

import { getServerSession } from "@/auth/server-session";
import { OAuthProviderList } from "@/components/auth/OAuthProviderList";
import { UnreadyBanner } from "@/components/site/Unready";
import { getProviders } from "@/lib/auth/providers";
import { getSettings } from "@/lib/settings/settings-server";

export default async function Page() {
  const session = await getServerSession();
  if (session) {
    redirect("/");
  }

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
