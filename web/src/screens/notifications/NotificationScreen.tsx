import { Box, Heading, Text } from "@chakra-ui/react";
import { Mailbox } from "src/components/graphics/Mailbox";

export function NotificationScreen() {
  return (
    <Box p={4}>
      <Box>
        <Mailbox width="4em" height="auto" />
      </Box>

      <Box>
        <Heading size="sm">Notifications</Heading>
        <Text>You have no notifications.</Text>
      </Box>
    </Box>
  );
}
