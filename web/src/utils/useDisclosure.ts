import { UseDisclosureProps, useDisclosure } from "@chakra-ui/react";

// Disclosure
// TODO: Copy into our codebase:
// https://github.com/chakra-ui/chakra-ui/blob/main/packages/hooks/use-disclosure/src/index.ts

export { useDisclosure };
export type { UseDisclosureProps };
export type WithDisclosure<T> = UseDisclosureProps & T;
