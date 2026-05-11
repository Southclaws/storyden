import { Presence } from "@ark-ui/react";
import type { PropsWithChildren } from "react";

import { css } from "@/styled-system/css";
import { VStack, styled } from "@/styled-system/jsx";
import { AdmonitionVariantProps, admonition } from "@/styled-system/recipes";
import type { HTMLStyledProps } from "@/styled-system/types";

import { IconButton } from "./icon-button";
import { CloseIcon } from "./icons/Close";

export const _Admonition = styled("aside", admonition);

type AdmonitionControlProps = {
  value: boolean;
  onChange?: (visible: boolean) => void;
  title?: string;
};

export type AdmonitionProps = {
  [K in keyof AdmonitionControlProps]: AdmonitionControlProps[K];
} & AdmonitionVariantProps &
  Omit<
    HTMLStyledProps<"aside">,
    keyof AdmonitionVariantProps | keyof AdmonitionControlProps
  >;

export function Admonition(props: PropsWithChildren<AdmonitionProps>) {
  const { children, onChange, title, value, ...rest } = props;

  function handleClose() {
    onChange?.(false);
  }

  const [admonitionVariantProps, elementProps] =
    admonition.splitVariantProps(rest);

  return (
    <Presence
      className={css({ w: "full" })}
      present={value}
      lazyMount
      unmountOnExit
    >
      <_Admonition {...admonitionVariantProps} {...elementProps}>
        <VStack alignItems="start">
          {title && <styled.h1 fontWeight="bold">{title}</styled.h1>}
          {children}
        </VStack>
        <IconButton
          type="button"
          size="xs"
          variant="ghost"
          aria-label="Close"
          onClick={handleClose}
        >
          <CloseIcon />
        </IconButton>
      </_Admonition>
    </Presence>
  );
}
