import { redirect } from "next/navigation";
import { LoginScreen } from "src/screens/auth/LoginScreen/LoginScreen";

import { getServerSession } from "@/auth/server-session";
import { OAuthProviderList } from "@/components/auth/OAuthProviderList";
import { UnreadyBanner } from "@/components/site/Unready";
import { getProviders } from "@/lib/auth/providers";

export async function LoginPage() {
  const session = await getServerSession();
  if (session) {
    redirect("/");
  }

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
