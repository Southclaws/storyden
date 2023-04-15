import { HStack } from "@chakra-ui/react";
import { Bell, Home, Login } from "src/components/Action/Action";

type Props = { isAuthenticated: boolean };

export function Toolbar(props: Props) {
  return (
    <HStack gap={2} pb={2}>
      {props.isAuthenticated ? (
        <>
          <Home />
          <Bell />
        </>
      ) : (
        <>
          <Home />
          <Login />
        </>
      )}
    </HStack>
  );
}
