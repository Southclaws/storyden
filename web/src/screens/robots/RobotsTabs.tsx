"use client";

import { usePathname, useRouter } from "next/navigation";
import { type PropsWithChildren } from "react";

import { Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import * as Tabs from "@/components/ui/tabs";
import { hasPermission } from "@/utils/permissions";

const routes = {
  robots: "/robots",
  toolsets: "/robots/toolsets",
  sessions: "/robots/chats",
  trails: "/robots/trails",
} as const;

type RobotTab = keyof typeof routes;

export function RobotsTabs({ children }: PropsWithChildren) {
  const pathname = usePathname();
  const router = useRouter();
  const session = useSession();
  const canUseRobots = hasPermission(
    session,
    Permission.USE_ROBOTS,
    Permission.MANAGE_ROBOTS,
  );
  const canManageTrails = hasPermission(session, Permission.MANAGE_TRAILS);

  const activeTab: RobotTab = pathname.startsWith(routes.trails)
    ? "trails"
    : pathname.startsWith(routes.sessions)
      ? "sessions"
      : pathname.startsWith(routes.toolsets)
        ? "toolsets"
        : "robots";

  return (
    <Tabs.Root
      className="robots-tabs"
      activationMode="manual"
      size="sm"
      value={activeTab}
      variant="line"
      onValueChange={({ value }) => router.push(routes[value as RobotTab])}
    >
      <Tabs.List className="robots-tabs__list" aria-label="Robots sections">
        {canUseRobots && (
          <>
            <Tabs.Trigger value="robots">Robots</Tabs.Trigger>
            <Tabs.Trigger value="toolsets">Toolsets</Tabs.Trigger>
            <Tabs.Trigger value="sessions">Sessions</Tabs.Trigger>
          </>
        )}
        {canManageTrails && <Tabs.Trigger value="trails">Trails</Tabs.Trigger>}
        <Tabs.Indicator />
      </Tabs.List>

      <Tabs.Content className="robots-tabs__content" value={activeTab}>
        {children}
      </Tabs.Content>
    </Tabs.Root>
  );
}
