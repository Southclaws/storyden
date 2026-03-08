"use client";

import Link from "next/link";

import { useRobotGet } from "@/api/openapi-client/robots";
import { RobotConfigurationForm } from "@/components/robots/RobotConfiguration/RobotConfigurationForm";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { IconButton } from "@/components/ui/icon-button";
import { ArrowLeftIcon } from "@/components/ui/icons/Arrow";
import { LinkButton } from "@/components/ui/link-button";
import { HStack, LStack, WStack } from "@/styled-system/jsx";

type Props = {
  robotId: string;
};

export function RobotConfigurationScreen({ robotId }: Props) {
  const { data, error } = useRobotGet(robotId);

  if (!data) {
    return <UnreadyBanner error={error} />;
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
            {data.name}
          </Heading>
        </HStack>

        <LinkButton variant="subtle" size="xs" href="/robots/chats">
          Chats
        </LinkButton>
      </WStack>

      <RobotConfigurationForm robot={data} />
    </LStack>
  );
}
