import { OnboardingStatus } from "src/api/openapi-schema";
import { useSession } from "src/auth";
import { useDisclosure } from "src/utils/useDisclosure";

export type Step = 1 | 2 | 3 | 4 | 5;

export const statusToStep: Record<OnboardingStatus, Step> = {
  requires_first_account: 1,
  requires_category: 2,
  requires_first_post: 3,
  requires_more_accounts: 4,
  complete: 5,
};

export function useChecklist() {
  const { onOpen, isOpen, onClose } = useDisclosure();
  const account = useSession();

  return {
    onOpen,
    isOpen,
    onClose,
    isLoggedIn: Boolean(account),
  };
}

export function isComplete(step: Step, status: OnboardingStatus) {
  return statusToStep[status] > step;
}
