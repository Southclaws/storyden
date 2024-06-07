import { Presence } from "@ark-ui/react";
import type { PropsWithChildren } from "react";

import { css } from "@/styled-system/css";
import { VStack, styled } from "@/styled-system/jsx";
import { AdmonitionVariantProps, admonition } from "@/styled-system/recipes";

import { Button } from "./button";

export const _Admonition = styled("aside", admonition);

export type AdmonitionProps = {
  value: boolean;
  onChange?: (visible: boolean) => void;
  title?: string;
} & AdmonitionVariantProps;

export function Admonition(props: PropsWithChildren<AdmonitionProps>) {
  function handleClose() {
    props.onChange?.(false);
  }

  const [admonitionVariantProps] = admonition.splitVariantProps(props);

  return (
    <Presence
      className={css({ w: "full" })}
      present={props.value}
      lazyMount
      unmountOnExit
    >
      <_Admonition {...admonitionVariantProps}>
        <VStack alignItems="start">
          {props.title && (
            <styled.h1 fontWeight="bold">{props.title}</styled.h1>
          )}
          {props.children}
        </VStack>
        <Button size="xs" variant="ghost" onClick={handleClose}>
          Close
        </Button>
      </_Admonition>
    </Presence>
  );
}
