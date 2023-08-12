import { StackProps, VStack } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

export function SettingsSection({
  children,
  ...props
}: PropsWithChildren<StackProps>) {
  return (
    <VStack
      w="full"
      gap={2}
      borderWidth={1}
      borderStyle="solid"
      borderColor="blackAlpha.100"
      borderRadius={10}
      p={4}
      alignItems="start"
      {...props}
    >
      {children}
    </VStack>
  );
}
