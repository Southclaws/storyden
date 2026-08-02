"use client";

import { LinkButton } from "@/components/ui/link-button";
import { PageHeading } from "@/components/ui/page-heading";
import RobotListScreen from "@/screens/robots/RobotListScreen";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

export default function Page() {
  return (
    <LStack>
      <WStack>
        <HStack gap="2">
          <PageHeading>Robots</PageHeading>
        </HStack>

        <LinkButton href="/robots/new">New</LinkButton>
      </WStack>

      <RobotListScreen />
    </LStack>
  );
}
