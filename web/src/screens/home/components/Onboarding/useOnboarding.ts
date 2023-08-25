import { useDisclosure } from "@chakra-ui/react";

import { useSession } from "src/auth";

export function useOnboarding() {
  const { onOpen, isOpen, onClose } = useDisclosure();
  const account = useSession();

  return { onOpen, isOpen, onClose, isLoggedIn: Boolean(account) };
}
