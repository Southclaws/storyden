"use client";

import { Mailbox } from "src/components/graphics/Mailbox";
import { Heading1 } from "src/theme/components/Heading/Index";

import { Box } from "@/styled-system/jsx";

export function NotificationScreen() {
  return (
    <Box p="4">
      <Box>
        <Mailbox width="4em" height="auto" />
      </Box>

      <Box>
        <Heading1 size="sm">Notifications</Heading1>
        <p>You have no notifications.</p>
      </Box>
    </Box>
  );
}
