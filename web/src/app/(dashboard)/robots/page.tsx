"use client";

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

        <LinkButton href="/robots/new">New</LinkButton>
      </WStack>

      <RobotListScreen />
    </LStack>
  );
}
