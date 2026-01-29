import { Robot } from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { Badge } from "@/components/ui/badge";
import { Heading } from "@/components/ui/heading";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

type Props = {
  robot: Robot;
  onClick: () => void;
};

export function RobotCard({ robot, onClick }: Props) {
  const toolCount = robot.tools.length;
  const toolCountLabel = `${toolCount} tool${toolCount === 1 ? "" : "s"}`;

  return (
    <CardBox
      cursor="pointer"
      onClick={onClick}
      _hover={{ background: "bg.emphasized" }}
    >
      <LStack gap="2">
        <WStack alignItems="start">
          <LStack gap="1">
            <Heading size="sm">{robot.name}</Heading>
            <styled.p fontSize="sm" color="fg.muted">
              {robot.description}
            </styled.p>
          </LStack>
        </WStack>

        <styled.p
          fontSize="sm"
          color="fg.subtle"
          lineClamp={2}
          fontFamily="mono"
        >
          {robot.playbook}
        </styled.p>

        <HStack justify="space-between">
          <MemberIdent profile={robot.author} size="sm" name="handle" />

          <Badge size="sm" variant="outline">
            {toolCountLabel}
          </Badge>
        </HStack>
      </LStack>
    </CardBox>
  );
}
