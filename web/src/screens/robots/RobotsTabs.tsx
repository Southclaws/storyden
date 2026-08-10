"use client";

import { usePathname, useRouter } from "next/navigation";
import { type PropsWithChildren } from "react";

import * as Tabs from "@/components/ui/tabs";

const routes = {
  robots: "/robots",
  toolsets: "/robots/toolsets",
  sessions: "/robots/chats",
} as const;

type RobotTab = keyof typeof routes;

export function RobotsTabs({ children }: PropsWithChildren) {
  const pathname = usePathname();
  const router = useRouter();
  const activeTab: RobotTab = pathname.startsWith(routes.sessions)
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
        <Tabs.Trigger value="robots">Robots</Tabs.Trigger>
        <Tabs.Trigger value="toolsets">Toolsets</Tabs.Trigger>
        <Tabs.Trigger value="sessions">Sessions</Tabs.Trigger>
        <Tabs.Indicator />
      </Tabs.List>

      <Tabs.Content className="robots-tabs__content" value={activeTab}>
        {children}
      </Tabs.Content>
    </Tabs.Root>
  );
}
