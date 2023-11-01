import { PropsWithChildren } from "react";

import { Heading } from "src/theme/components";

export function Title({ children }: PropsWithChildren) {
  return (
    <Heading size="sm" role="navigation" wordBreak="keep-all">
      {children}
    </Heading>
  );
}
