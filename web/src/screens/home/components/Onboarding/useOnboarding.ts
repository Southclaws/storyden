import { useDisclosure } from "@chakra-ui/react";

import { OnboardingStatus } from "src/api/openapi/schemas";
import { useSession } from "src/auth";

export function useOnboarding() {
  const { onOpen, isOpen, onClose } = useDisclosure();
  const account = useSession();

  return { onOpen, isOpen, onClose, isLoggedIn: Boolean(account) };
}

export function isOnboarding(status?: OnboardingStatus) {
  switch (status) {
    // NOTE: explicit exhaustivity here because we want to default to content
    // if something went wrong with the info data. Last resort: show content!
    case "requires_first_account":
    case "requires_category":
    case "requires_more_accounts":
    case "requires_first_post":
      return true;

    case "complete":
    default:
      return false;
  }
}
