"use client";

import { useRobotsList } from "@/api/openapi-client/robots";
import { RobotCard } from "@/components/robots/RobotCard";
import { EmptyState } from "@/components/site/EmptyState";
import { Unready } from "@/components/site/Unready";
import { RobotIcon } from "@/components/ui/icons/Robot";
import { CardRows } from "@/components/ui/surface";

export default function RobotListScreen() {
  const { data, error } = useRobotsList();

  if (!data) {
    return <Unready error={error} />;
  }

  if (data.robots.length === 0) {
    return (
      <EmptyState icon={<RobotIcon />}>
        No specialist Robots yet. Denbot can handle any conversation and
        delegate once you create one.
      </EmptyState>
    );
  }

  return (
    <CardRows>
      {data.robots.map((robot) => (
        <RobotCard key={robot.id} robot={robot} />
      ))}
    </CardRows>
  );
}
