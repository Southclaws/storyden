import { Suspense } from "react";

import { Unready } from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";

import { RegisterPage } from "./RegisterPage";

export default function Page() {
  return (
    <Suspense fallback={<Unready />}>
      <RegisterPage />
    </Suspense>
  );
}

export async function generateMetadata() {
  const settings = await getSettings();
  return {
    title: `Join the community at ${settings.title}`,
    description: `Log in or sign up to ${settings.title} - powered by Storyden`,
  };
}
