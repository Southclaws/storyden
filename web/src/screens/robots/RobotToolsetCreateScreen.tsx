"use client";

import { useRouter } from "next/navigation";

import { RobotToolset } from "@/api/openapi-schema";
import { RobotToolsetConfigurationForm } from "@/components/robots/RobotToolsetConfiguration/RobotToolsetConfigurationForm";
import { BackAction } from "@/components/site/Action/Back";
import { PageHeading } from "@/components/ui/page-heading";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

export function RobotToolsetCreateScreen() {
  const router = useRouter();

  function handleSave(toolset: RobotToolset) {
    router.push(`/robots/toolsets/${toolset.id}`);
  }

  return (
    <LStack gap="4" w="full">
      <WStack>
        <HStack gap="2">
          <BackAction fallbackHref="/robots/toolsets" />
          <PageHeading>New Toolset</PageHeading>
        </HStack>
      </WStack>
      <RobotToolsetConfigurationForm onSave={handleSave} />
    </LStack>
  );
}
