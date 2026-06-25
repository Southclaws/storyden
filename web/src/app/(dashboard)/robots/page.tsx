"use client";

import { ButtonGroup } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import RobotListScreen from "@/screens/robots/RobotListScreen";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

export default function Page() {
  return (
    <LStack>
      <WStack>
        <HStack gap="2">
          <Heading size="md">Robots</Heading>
        </HStack>

        <ButtonGroup attached>
          <LinkButton variant="subtle" size="xs" href="/robots/chats">
            Chats
          </LinkButton>

          <LinkButton href="/robots/new" variant="subtle" size="xs">
            New
          </LinkButton>
        </ButtonGroup>
      </WStack>

      <RobotListScreen />
    </LStack>
  );
}
