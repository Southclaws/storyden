"use client";

import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import { Text } from "@/components/ui/text";
import RobotListScreen from "@/screens/admin/RobotsSettingsScreen/RobotListScreen";
import { LStack, WStack } from "@/styled-system/jsx";

export default function Page() {
  return (
    <LStack>
      <WStack>
        <Heading size="md">Robots</Heading>

        <LinkButton variant="subtle" size="xs" href="/robots/chats">
          Chats
        </LinkButton>
      </WStack>

      <Text color="fg.muted">
        Robots are language-model driven automations for organising your
        community.
      </Text>

      <RobotListScreen />
    </LStack>
  );
}
