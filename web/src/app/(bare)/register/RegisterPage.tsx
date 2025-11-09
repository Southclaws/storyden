import { redirect } from "next/navigation";
import { RegisterScreen } from "src/screens/auth/RegisterScreen/RegisterScreen";

import { getServerSession } from "@/auth/server-session";
import { OAuthProviderList } from "@/components/auth/OAuthProviderList";
import { UnreadyBanner } from "@/components/site/Unready";
import { getProviders } from "@/lib/auth/providers";

export async function RegisterPage() {
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
