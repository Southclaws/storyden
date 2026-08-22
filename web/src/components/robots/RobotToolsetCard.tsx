import { RobotToolset } from "@/api/openapi-schema";
import { Badge } from "@/components/ui/badge";
import { ToolsetIcon } from "@/components/ui/icons/Toolset";
import { Card } from "@/components/ui/surface";
import { HStack, WStack } from "@/styled-system/jsx";
import { pluralise } from "@/utils/text";

export function RobotToolsetCard({ toolset }: { toolset: RobotToolset }) {
  const toolCount = toolset.tools.length;

  return (
    <Card
      id={toolset.id}
      shape="row"
      title={toolset.name}
      titleIcon={<ToolsetIcon />}
      text={toolset.description}
      url={`/robots/toolsets/${toolset.id}`}
    >
      <WStack w="full">
        <HStack gap="1">
          <Badge size="sm" variant="outline">
            {toolset.source}
          </Badge>
          <Badge size="sm" variant="outline">
            {toolCount} {pluralise(toolCount, "tool")}
          </Badge>
        </HStack>
        {toolset.usage_count > 0 && (
          <Badge size="sm" variant="subtle">
            {toolset.usage_count} {pluralise(toolset.usage_count, "Robot")}
          </Badge>
        )}
      </WStack>
    </Card>
  );
}
