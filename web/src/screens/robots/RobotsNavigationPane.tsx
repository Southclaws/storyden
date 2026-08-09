"use client";

import { usePathname } from "next/navigation";

import { type Robot, type RobotSessionRef } from "@/api/openapi-schema";
import {
  SidebarNavigationLink,
  SidebarNavigationPane,
  SidebarNavigationSection,
} from "@/components/site/Navigation/SidebarNavigation";
import { AddIcon } from "@/components/ui/icons/Add";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { RobotIcon } from "@/components/ui/icons/Robot";
import { LinkButton } from "@/components/ui/link-button";

type Props = {
  robots: Robot[];
  sessions: RobotSessionRef[];
};

export function RobotsNavigationPane({ robots, sessions }: Props) {
  const pathname = usePathname();

  return (
    <SidebarNavigationPane label="Robots navigation" title="Robots">
      <SidebarNavigationSection>
        <SidebarNavigationLink
          href="/robots"
          label="All robots"
          icon={<RobotIcon />}
          active={pathname === "/robots"}
          action={
            <LinkButton
              aria-label="New robot"
              href="/robots/new"
              px="0"
              size="sm"
              variant="plain"
            >
              <AddIcon />
            </LinkButton>
          }
        />
        <SidebarNavigationLink
          href="/robots/chats"
          label="All chats"
          icon={<DiscussionIcon />}
          active={pathname === "/robots/chats"}
          action={
            <LinkButton
              aria-label="New chat"
              href="/robots/chats/new"
              px="0"
              size="sm"
              variant="plain"
            >
              <AddIcon />
            </LinkButton>
          }
        />
      </SidebarNavigationSection>

      {robots.length > 0 && (
        <SidebarNavigationSection label="Robots">
          {robots.map((robot) => (
            <SidebarNavigationLink
              key={robot.id}
              href={`/robots/${robot.id}`}
              label={robot.name}
              icon={<RobotIcon />}
              active={pathname === `/robots/${robot.id}`}
            />
          ))}
        </SidebarNavigationSection>
      )}

      {sessions.length > 0 && (
        <SidebarNavigationSection label="Recent chats">
          {sessions.slice(0, 6).map((session) => (
            <SidebarNavigationLink
              key={session.id}
              href={`/robots/chats/${session.id}`}
              label={session.name}
              icon={<DiscussionIcon />}
              active={pathname === `/robots/chats/${session.id}`}
            />
          ))}
        </SidebarNavigationSection>
      )}
    </SidebarNavigationPane>
  );
}
