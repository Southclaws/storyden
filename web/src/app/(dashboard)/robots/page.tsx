"use client";

import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import { Text } from "@/components/ui/text";
import RobotListScreen from "@/screens/robots/RobotListScreen";
import { LStack, WStack } from "@/styled-system/jsx";

export default function Page() {
  return (
    <LStack>
      <WStack>
        <LinkButton variant="subtle" size="xs" href="/robots/chats">
          Chats
        </LinkButton>
      </WStack>

      <RobotListScreen />
    </LStack>
  );
}
