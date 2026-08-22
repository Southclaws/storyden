"use client";

import { usePathname } from "next/navigation";

import {
  type Robot,
  type RobotSessionRef,
  type Trail,
} from "@/api/openapi-schema";
import {
  SidebarNavigationLink,
  SidebarNavigationPane,
  SidebarNavigationSection,
} from "@/components/site/Navigation/SidebarNavigation";
import { AddIcon } from "@/components/ui/icons/Add";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { RobotIcon } from "@/components/ui/icons/Robot";
import { ToolsetIcon } from "@/components/ui/icons/Toolset";
import { TrailIcon } from "@/components/ui/icons/Trail";
import { LinkButton } from "@/components/ui/link-button";

type Props = {
  robots: Robot[];
  sessions: RobotSessionRef[];
  trails: Trail[];
  canUseRobots: boolean;
  canManageTrails: boolean;
};

export function RobotsNavigationPane({
  robots,
  sessions,
  trails,
  canUseRobots,
  canManageTrails,
}: Props) {
  const pathname = usePathname();

  return (
    <SidebarNavigationPane label="Robots navigation" title="Robots">
      <SidebarNavigationSection>
        {canUseRobots && (
          <>
            <SidebarNavigationLink
              href="/robots"
              label="Robots"
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
              href="/robots/toolsets"
              label="Toolsets"
              icon={<ToolsetIcon />}
              active={pathname.startsWith("/robots/toolsets")}
              action={
                <LinkButton
                  aria-label="New Toolset"
                  href="/robots/toolsets/new"
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
              label="Chats"
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
          </>
        )}
        {canManageTrails && (
          <SidebarNavigationLink
            href="/robots/trails"
            label="Trails"
            icon={<TrailIcon />}
            active={pathname.startsWith("/robots/trails")}
            action={
              <LinkButton
                aria-label="New Trail"
                href="/robots/trails/new"
                px="0"
                size="sm"
                variant="plain"
              >
                <AddIcon />
              </LinkButton>
            }
          />
        )}
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

      {canManageTrails && trails.length > 0 && (
        <SidebarNavigationSection label="Recent Trails">
          {trails.slice(0, 5).map((trail) => (
            <SidebarNavigationLink
              key={trail.id}
              href={`/robots/trails/${trail.id}`}
              label={trail.name}
              icon={<TrailIcon />}
              active={pathname === `/robots/trails/${trail.id}`}
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
