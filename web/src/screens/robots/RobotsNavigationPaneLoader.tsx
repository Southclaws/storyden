"use client";

import {
  useRobotSessionsList,
  useRobotsList,
} from "@/api/openapi-client/robots";
import { useTrailList } from "@/api/openapi-client/trails";
import { Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { hasPermission } from "@/utils/permissions";

import { RobotsNavigationPane } from "./RobotsNavigationPane";

export function RobotsNavigationPaneLoader() {
  const session = useSession();
  const canUseRobots = hasPermission(
    session,
    Permission.USE_ROBOTS,
    Permission.MANAGE_ROBOTS,
  );
  const canManageTrails = hasPermission(session, Permission.MANAGE_TRAILS);

  const robots = useRobotsList(undefined, {
    swr: { enabled: canUseRobots },
  });
  const sessions = useRobotSessionsList(undefined, {
    swr: { enabled: canUseRobots },
  });
  const trails = useTrailList({ swr: { enabled: canManageTrails } });

  return (
    <RobotsNavigationPane
      robots={robots.data?.robots ?? []}
      sessions={sessions.data?.sessions ?? []}
      trails={trails.data?.trails ?? []}
      canUseRobots={canUseRobots}
      canManageTrails={canManageTrails}
    />
  );
}
