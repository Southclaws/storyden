import { Robot } from "@/api/openapi-schema";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { Badge } from "@/components/ui/badge";
import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import { CardBox, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

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

        <WStack>
          <MemberIdent profile={robot.author} size="sm" name="handle" />

          <HStack gap="2">
            <Badge size="sm" variant="outline">
              {toolCountLabel}
            </Badge>

            <LinkButton href={editHref} size="xs" variant="subtle">
              Edit
            </LinkButton>
          </HStack>
        </WStack>
      </LStack>
    </CardBox>
  );
}
