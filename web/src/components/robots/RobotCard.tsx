import { Robot } from "@/api/openapi-schema";
import { Badge } from "@/components/ui/badge";
import { RobotIcon } from "@/components/ui/icons/Robot";
import { Card } from "@/components/ui/surface";
import { HStack } from "@/styled-system/jsx";

type Props = {
  robot: Robot;
  editHref?: string;
};

export function RobotCard({ robot, editHref = `/robots/${robot.id}` }: Props) {
  const toolsetCount = robot.toolsets.length;
  const toolsetCountLabel = `${toolsetCount} Toolset${toolsetCount === 1 ? "" : "s"}`;

  return (
    <Card
      id={robot.id}
      shape="row"
      title={robot.name}
      titleIcon={<RobotIcon />}
      text={robot.description}
      content={robot.playbook}
      url={editHref}
    >
      <HStack gap="2">
        <Badge size="sm" variant="outline">
          {toolsetCountLabel}
        </Badge>
        <Badge size="sm" variant="outline">
          {robot.model}
        </Badge>
      </HStack>
    </Card>
  );
}
