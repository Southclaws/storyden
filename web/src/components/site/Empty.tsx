import { CubeTransparentIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { css } from "@/styled-system/css";
import { HStack, styled } from "@/styled-system/jsx";

export function Empty({ children }: PropsWithChildren) {
  return (
    <HStack alignItems="center" color="fg.muted">
      <CubeTransparentIcon className={css({ width: "6" })} />
      <styled.p fontStyle="italic">{children}</styled.p>
    </HStack>
  );
}
