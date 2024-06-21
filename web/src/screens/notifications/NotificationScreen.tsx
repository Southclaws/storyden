"use client";

import { Mailbox } from "src/components/graphics/Mailbox";

import { Heading } from "@/components/ui/heading";
import { Box } from "@/styled-system/jsx";

export function NotificationScreen() {
  return (
    <Box p="4">
      <Box>
        <Mailbox width="4em" height="auto" />
      </Box>

      <Box>
        <Heading size="sm">Notifications</Heading>
        <p>You have no notifications.</p>
      </Box>
    </Box>
  );
}
