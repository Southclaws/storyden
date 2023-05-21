import { HStack } from "@chakra-ui/react";
import { useSession } from "src/auth";
import { Bell, Create, Home, Login } from "src/components/Action/Action";

export function Toolbar() {
  const account = useSession();
  return (
    <HStack gap={2} pb={2}>
      {account ? (
        <>
          <Home />
          <Bell />
          <Create />
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
