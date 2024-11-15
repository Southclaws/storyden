import { RegisterScreen } from "src/screens/auth/RegisterScreen/RegisterScreen";

import { UnreadyBanner } from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";

export default function Page() {
  try {
    return <RegisterScreen />;
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
