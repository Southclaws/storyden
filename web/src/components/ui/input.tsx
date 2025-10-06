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
      px="3"
      py="2"
      bg="bg.muted"
      borderStyle="solid"
      borderColor="border.default"
      borderWidth="thin"
      borderRightStyle="none"
      borderTopLeftRadius="l2"
      borderBottomLeftRadius="l2"
      display="flex"
      alignItems="center"
      color="fg.muted"
      fontSize="sm"
      {...(props as any)}
    >
      {children}
    </styled.div>
  );
}
