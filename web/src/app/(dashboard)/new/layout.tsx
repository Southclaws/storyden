import { Metadata } from "next";
import { PropsWithChildren } from "react";

import { getSettings } from "@/lib/settings/settings-server";

export default async function Layout({ children }: PropsWithChildren) {
  return <>{children}</>;
}

export async function generateMetadata(): Promise<Metadata> {
  const settings = await getSettings();

  return {
    title: `Draft a new post on ${settings.title}`,
    description: `Compose a new masterpice and share it with the community on ${settings.title}`,
  };
}
