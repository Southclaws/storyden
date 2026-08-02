import { Robot } from "@/api/openapi-schema";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { Badge } from "@/components/ui/badge";
import { CardBox } from "@/components/ui/card-box";
import { LinkButton } from "@/components/ui/link-button";
import { Text } from "@/components/ui/text";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

type Props = {
  robot: Robot;
  editHref?: string;
};

export function RobotCard({ robot, editHref = `/robots/${robot.id}` }: Props) {
  const toolCount = robot.tools.length;
  const toolCountLabel = `${toolCount} tool${toolCount === 1 ? "" : "s"}`;

  return (
    <CardBox>
      <LStack gap="2">
        <WStack alignItems="start">
          <LStack gap="1">
            <Text
              variant="supporting"
              color="text.default"
              fontWeight="semibold"
            >
              {robot.name}
            </Text>
            <Text variant="supporting">{robot.description}</Text>
          </LStack>
        </WStack>

        <Text variant="supporting" lineClamp={2} fontFamily="mono">
          {robot.playbook}
        </Text>

        <WStack>
          <MemberIdent profile={robot.author} size="sm" name="handle" />

          <HStack gap="2">
            <Badge size="sm" variant="outline">
              {toolCountLabel}
            </Badge>

            <LinkButton href={editHref} variant="subtle">
              Edit
            </LinkButton>
          </HStack>
        </WStack>
      </LStack>
    </CardBox>
  );
}
