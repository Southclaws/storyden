import { PropsWithChildren } from "react";

import { VStack, VstackProps } from "@/styled-system/jsx";

export function SettingsSection({
  children,
  ...props
}: PropsWithChildren<VstackProps>) {
  return (
    <VStack
      w="full"
      gap="2"
      borderWidth="thin"
      borderStyle="solid"
      borderColor="blackAlpha.100"
      borderRadius="lg"
      p="4"
      alignItems="start"
      {...(props as any)}
    >
      {children}
    </VStack>
  );
}
