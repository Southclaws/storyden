"use client";

import { useRouter } from "next/navigation";

import { useRobotsList } from "@/api/openapi-client/robots";
import { RobotCard } from "@/components/robots/RobotCard";
import { Unready } from "@/components/site/Unready";
import { LStack } from "@/styled-system/jsx";

export default function RobotListScreen() {
  const router = useRouter();
  const { data, error } = useRobotsList();

  if (!data) {
    return <Unready error={error} />;
  }

  return (
    <LStack>
      {data.robots.map((robot) => (
        <RobotCard
          key={robot.id}
          robot={robot}
          onClick={() => router.push(`/robots/${robot.id}`)}
        />
      ))}
    </LStack>
  );
}
