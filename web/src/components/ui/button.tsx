import { ark } from "@ark-ui/react/factory";
import type { ComponentProps } from "react";
import { forwardRef } from "react";

import { styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { Spinner } from "./Spinner";

export const StyledButton = styled(ark.button, button);
export interface StyledButtonProps
  extends ComponentProps<typeof StyledButton> {}

interface ButtonLoadingProps {
  loading?: boolean;
  loadingText?: React.ReactNode;
}

export interface ButtonProps extends StyledButtonProps, ButtonLoadingProps {}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (props, ref) => {
    const { loading, disabled, loadingText, children, ...rest } = props;

    const trulyDisabled = loading || disabled;

    return (
      <StyledButton disabled={trulyDisabled} ref={ref} {...rest}>
        {loading && !loadingText ? (
          <>
            <Spinner />
          </>
        ) : loadingText ? (
          loadingText
        ) : (
          children
        )}
      </StyledButton>
    );
  },
);

Button.displayName = "Button";
