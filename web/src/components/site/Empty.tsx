import { CubeTransparentIcon } from "@heroicons/react/24/outline";
import { PropsWithChildren } from "react";

import { css } from "@/styled-system/css";
import { HStack, styled } from "@/styled-system/jsx";

const iconStyles = css({ width: "6" });

export function Empty({ children }: PropsWithChildren) {
  return (
    <HStack alignItems="center" color="fg.muted">
      <CubeTransparentIcon className={iconStyles} />
      <styled.p fontStyle="italic" textWrap="nowrap">
        {children}
      </styled.p>
    </HStack>
  );
}
