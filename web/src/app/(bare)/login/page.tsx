import { LoginScreen } from "src/screens/auth/LoginScreen/LoginScreen";

import { UnreadyBanner } from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";

export default function Page() {
  try {
    return <LoginScreen />;
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
