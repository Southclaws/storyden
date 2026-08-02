import { ark } from "@ark-ui/react/factory";
import type { ComponentProps, PropsWithChildren } from "react";
import { HstackProps, styled } from "styled-system/jsx";
import { input } from "styled-system/recipes";

export const Input = styled(ark.input, input);
export interface InputProps extends ComponentProps<typeof Input> {}

// Shove this to the left side of <Input /> in a HStack to make a lil prefix box
export function InputPrefix({
  children,
  ...props
}: PropsWithChildren<HstackProps>) {
  return (
    <styled.div
      boxSizing="border-box"
      h="6"
      px="2"
      bg="background.control"
      borderStyle="solid"
      borderColor="border.default"
      borderWidth="thin"
      borderRightStyle="none"
      borderTopLeftRadius="sm"
      borderBottomLeftRadius="sm"
      display="flex"
      alignItems="center"
      color="text.subtle"
      fontSize="xs"
      lineHeight="1.125rem"
      whiteSpace="nowrap"
      {...(props as any)}
    >
      {children}
    </styled.div>
  );
}
