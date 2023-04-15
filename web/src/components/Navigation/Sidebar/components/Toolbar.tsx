import { HStack } from "@chakra-ui/react";
import { Home } from "src/components/Action/Action";

export function Toolbar() {
  return (
    <HStack gap={2} pb={2}>
      <Home />
    </HStack>
  );
}
