"use client";

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
        background="bg.default"
        p="4"
        zIndex="banner"
        boxShadow="2xl"
        borderRadius="2xl"
        w="full"
        maxH="breakpoint-sm"
        overflowY="scroll"
      >
        <Checklist onboardingStatus={onboardingStatus} onFinish={onFinish} />
      </Box>
    </Box>
  );
}
