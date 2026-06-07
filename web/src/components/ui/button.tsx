import { ark } from "@ark-ui/react/factory";
import { createContext, mergeProps } from "@ark-ui/react/utils";
import type { ComponentProps } from "react";
import { forwardRef, useMemo } from "react";

import { styled } from "@/styled-system/jsx";
import {
  ButtonVariantProps,
  button,
  group,
} from "@/styled-system/recipes";
import { cx } from "@/styled-system/css";

import { Spinner } from "./Spinner";
import { Group, GroupProps } from "./group";

export const StyledButton = styled(ark.button, button);
export interface StyledButtonProps extends ComponentProps<
  typeof StyledButton
> {}

interface ButtonLoadingProps {
  loading?: boolean;
  loadingText?: React.ReactNode;
}

export interface ButtonProps extends StyledButtonProps, ButtonLoadingProps {}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (props, ref) => {
    const propsContext = useButtonPropsContext();
    const buttonProps = useMemo(
      () => mergeProps<ButtonProps>(propsContext, props),
      [propsContext, props],
    );
    const { loading, disabled, loadingText, children, ...rest } = buttonProps;

    const trulyDisabled = loading || disabled;

    return (
      <StyledButton disabled={trulyDisabled} ref={ref} {...rest}>
        {loading ? (
          loadingText ? (
            loadingText
          ) : (
            <>
              <Spinner />
            </>
          )
        ) : (
          children
        )}
      </StyledButton>
    );
  },
);

Button.displayName = "Button";

export interface ButtonGroupProps extends GroupProps, ButtonVariantProps {}

export const ButtonGroup = forwardRef<HTMLDivElement, ButtonGroupProps>(
  function ButtonGroup(props, ref) {
    const [buttonVariantProps, groupVariantProps, otherProps] = useMemo(
      () => {
        const [buttonVariantProps, groupProps] =
          button.splitVariantProps(props);
        const [groupVariantProps, otherProps] =
          group.splitVariantProps(groupProps);

        return [buttonVariantProps, groupVariantProps, otherProps];
      },
      [props],
    );

    return (
      <ButtonPropsProvider value={buttonVariantProps}>
        <Group
          ref={ref}
          {...otherProps}
          className={cx(group(groupVariantProps), otherProps.className)}
        />
      </ButtonPropsProvider>
    );
  },
);

const [ButtonPropsProvider, useButtonPropsContext] =
  createContext<ButtonVariantProps>({
    name: "ButtonPropsContext",
    hookName: "useButtonPropsContext",
    providerName: "<ButtonPropsProvider />",
    strict: false,
  });

export { useButtonPropsContext };
