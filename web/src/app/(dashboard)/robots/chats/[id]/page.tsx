"use client";

import { useParams, useSearchParams } from "next/navigation";

import { RobotSessionScreen } from "@/screens/robots/RobotSessionScreen";

export default function Page() {
  const { id } = useParams<{ id: string }>();
  const searchParams = useSearchParams();

  return (
    <RobotSessionScreen
      sessionId={id === "new" ? undefined : id}
      initialChatBefore={searchParams.get("before") ?? undefined}
      initialChatLimit={searchParams.get("limit") ?? undefined}
      initialSelectedRobotID={searchParams.get("robot") ?? undefined}
    />
  );
}
