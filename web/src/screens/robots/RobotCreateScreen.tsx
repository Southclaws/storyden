"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";

import { RobotCreateOKResponse } from "@/api/openapi-schema";
import { RobotConfigurationForm } from "@/components/robots/RobotConfiguration/RobotConfigurationForm";
import { Heading } from "@/components/ui/heading";
import { IconButton } from "@/components/ui/icon-button";
import { ArrowLeftIcon } from "@/components/ui/icons/Arrow";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

export function RobotCreateScreen() {
  const router = useRouter();

  function handleSave(robot: RobotCreateOKResponse) {
    router.push(`/robots/${robot.id}`);
  }

  return (
    <LStack gap="4" w="full">
      <WStack>
        <HStack gap="2">
          <Link href="/robots">
            <IconButton variant="ghost" size="xs">
              <ArrowLeftIcon />
            </IconButton>
          </Link>

          <Heading size="md" lineClamp="1">
            New robot
          </Heading>
        </HStack>
      </WStack>

      <RobotConfigurationForm onSave={handleSave} />
    </LStack>
  );
}
