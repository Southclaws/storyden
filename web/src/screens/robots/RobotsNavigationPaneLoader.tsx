"use client";

import {
  useRobotSessionsList,
  useRobotsList,
} from "@/api/openapi-client/robots";

import { RobotsNavigationPane } from "./RobotsNavigationPane";

export function RobotsNavigationPaneLoader() {
  const robots = useRobotsList();
  const sessions = useRobotSessionsList();

  return (
    <RobotsNavigationPane
      robots={robots.data?.robots ?? []}
      sessions={sessions.data?.sessions ?? []}
    />
  );
}
