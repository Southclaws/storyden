import { notFound } from "next/navigation";
import { PropsWithChildren } from "react";

import { InstanceCapability } from "@/api/openapi-schema";
import { hasCapability } from "@/lib/settings/capabilities";
import { getSettings } from "@/lib/settings/settings-server";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

export default async function Layout({ children }: PropsWithChildren) {
  const { capabilities } = await getSettings();

  if (!hasCapability(InstanceCapability.robots, capabilities)) {
    notFound();
  }

  return <>{children}</>;
}
