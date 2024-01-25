import { Box } from "@/styled-system/jsx";

import { Checklist } from "./Checklist/Checklist";
import { useOnboarding } from "./useOnboarding";

export function Onboarding() {
  const { onFinish, onboardingStatus, showOnboarding } = useOnboarding();

  if (!showOnboarding) return null;

  return (
    <Box position="relative">
      <Box
        position="absolute"
        // NOTE: not dark mode ready! need a variable
        background="bg.default"
        p="4"
      >
        <Checklist onboardingStatus={onboardingStatus} onFinish={onFinish} />
      </Box>
    </Box>
  );
}
