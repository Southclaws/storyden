import { RegisterScreen } from "src/screens/auth/RegisterScreen/RegisterScreen";

import { getSettings } from "@/lib/settings/settings-server";

export default function Page() {
  return <RegisterScreen />;
}

export async function generateMetadata() {
  const settings = await getSettings();
  return {
    title: `Join the community at ${settings.title}`,
    description: `Log in or sign up to ${settings.title} - powered by Storyden`,
  };
}
