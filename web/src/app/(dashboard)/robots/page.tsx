"use client";

import { BackAction } from "@/components/site/Action/Back";
import { LinkButton } from "@/components/ui/link-button";
import { PageHeading } from "@/components/ui/page-heading";
import RobotListScreen from "@/screens/robots/RobotListScreen";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

export default function Page() {
  return (
    <LStack>
      <WStack>
        <HStack gap="2">
          <BackAction className="robots-tabs__back" />
          <PageHeading>Robots</PageHeading>
        </HStack>

        <LinkButton href="/robots/new">New Robot</LinkButton>
      </WStack>

      <RobotListScreen />
    </LStack>
  );
}
