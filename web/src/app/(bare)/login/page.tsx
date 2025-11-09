import { Suspense } from "react";

import { Unready } from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";

import { LoginPage } from "./LoginPage";

export default function Page() {
  return (
    <Suspense fallback={<Unready />}>
      <LoginPage />
    </Suspense>
  );
}

export async function generateMetadata() {
  const settings = await getSettings();
  return {
    title: `Login to ${settings.title}`,
    description: `Log in or sign up to ${settings.title} - powered by Storyden`,
  };
}
