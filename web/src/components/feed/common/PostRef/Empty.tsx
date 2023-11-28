import { CubeTransparentIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { css } from "@/styled-system/css";
import { HStack, styled } from "@/styled-system/jsx";

export function Empty({ children }: PropsWithChildren) {
  return (
    <HStack alignItems="center">
      <CubeTransparentIcon className={css({ width: "6", color: "bg.muted" })} />
      <styled.p fontStyle="italic" color="gray.500">
        {children ?? "no posts"}
      </styled.p>
    </HStack>
  );
}
