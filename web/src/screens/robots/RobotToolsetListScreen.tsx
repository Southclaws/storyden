"use client";

import { useRobotToolsetsList } from "@/api/openapi-client/robots";
import { RobotToolsetCard } from "@/components/robots/RobotToolsetCard";
import { EmptyState } from "@/components/site/EmptyState";
import { UnreadyBanner } from "@/components/site/Unready";
import { ToolsetIcon } from "@/components/ui/icons/Toolset";
import { CardRows } from "@/components/ui/surface";

export function RobotToolsetListScreen() {
  const { data, error } = useRobotToolsetsList();

  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  if (data.toolsets.length === 0) {
    return (
      <EmptyState icon={<ToolsetIcon />}>
        No Toolsets are available yet.
      </EmptyState>
    );
  }

  return (
    <CardRows>
      {data.toolsets.map((toolset) => (
        <RobotToolsetCard key={toolset.id} toolset={toolset} />
      ))}
    </CardRows>
  );
}
