import { robotSessionsList, robotsList } from "@/api/openapi-server/robots";

import { RobotsNavigationPane } from "./RobotsNavigationPane";

export async function RobotsNavigationPaneLoader() {
  const [robotsResult, sessionsResult] = await Promise.allSettled([
    robotsList(),
    robotSessionsList(),
  ]);

  const robots =
    robotsResult.status === "fulfilled" ? robotsResult.value.data.robots : [];
  const sessions =
    sessionsResult.status === "fulfilled"
      ? sessionsResult.value.data.sessions
      : [];

  return <RobotsNavigationPane robots={robots} sessions={sessions} />;
}
