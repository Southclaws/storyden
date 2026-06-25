import { RegisterScreen } from "@/screens/auth/RegisterScreen/RegisterScreen";

import { RegistrationMode } from "@/api/openapi-schema";
import { OAuthProviderList } from "@/components/auth/OAuthProviderList";
import { UnreadyBanner } from "@/components/site/Unready";
import { getProviders } from "@/lib/auth/providers";
import { allowsPublicRegistration } from "@/lib/settings/registration";
import { getSettings } from "@/lib/settings/settings-server";

type Props = {
  searchParams: Promise<{
    invitation_id?: string;
  }>;
};

export default async function Page({ searchParams }: Props) {
  try {
    const [params, settings] = await Promise.all([searchParams, getSettings()]);
    const { oauth } = await getProviders();
    const invitationID = params.invitation_id;
    const canCreateOAuthAccount =
      settings.registration_mode === RegistrationMode.public;

    return (
      <>
        <RegisterScreen
          invitationID={invitationID}
          registrationMode={settings.registration_mode}
        />
        {canCreateOAuthAccount && <OAuthProviderList providers={oauth} />}
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
