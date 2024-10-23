import { LoginScreen } from "src/screens/auth/LoginScreen/LoginScreen";

import { getSettings } from "@/lib/settings/settings-server";

export default function Page() {
  return <LoginScreen />;
}

export async function generateMetadata() {
  const settings = await getSettings();
  return {
    title: `Login to ${settings.title}`,
    description: `Log in or sign up to ${settings.title} - powered by Storyden`,
  };
}
