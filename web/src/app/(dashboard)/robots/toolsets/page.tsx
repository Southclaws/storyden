"use client";

import { BackAction } from "@/components/site/Action/Back";
import { LinkButton } from "@/components/ui/link-button";
import { PageHeading } from "@/components/ui/page-heading";
import { RobotToolsetListScreen } from "@/screens/robots/RobotToolsetListScreen";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

export default function Page() {
  return (
    <LStack gap="4">
      <WStack>
        <HStack gap="2">
          <BackAction className="robots-tabs__back" />
          <PageHeading>Toolsets</PageHeading>
        </HStack>
        <LinkButton href="/robots/toolsets/new">New Toolset</LinkButton>
      </WStack>
      <RobotToolsetListScreen />
    </LStack>
  );
}
