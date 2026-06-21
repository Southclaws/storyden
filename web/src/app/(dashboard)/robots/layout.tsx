import { notFound } from "next/navigation";
import { PropsWithChildren } from "react";

import { InstanceCapability } from "@/api/openapi-schema";
import { hasCapability } from "@/lib/settings/capabilities";
import { getSettings } from "@/lib/settings/settings-server";

export default async function Layout({ children }: PropsWithChildren) {
  const { capabilities } = await getSettings();

  if (!hasCapability(InstanceCapability.robots, capabilities)) {
    notFound();
  }

  return <>{children}</>;
}
