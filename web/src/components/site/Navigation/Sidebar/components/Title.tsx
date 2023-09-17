import { Heading } from "@chakra-ui/react";
import { PropsWithChildren } from "react";

export function Title({ children }: PropsWithChildren) {
  return (
    <Heading size="sm" role="navigation" wordBreak="keep-all">
      {children}
    </Heading>
  );
}
