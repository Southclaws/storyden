import { ark } from "@ark-ui/react/factory";
import type { ComponentProps } from "react";
import { forwardRef } from "react";

import { styled } from "@/styled-system/jsx";
import { ButtonVariantProps, button } from "@/styled-system/recipes";

import { Spinner } from "./Spinner";

export type StyledIconButtonProps = ComponentProps<typeof StyledIconButton>;
export const StyledIconButton = styled(ark.button, button, {
  defaultProps: { px: "0" } as ButtonVariantProps,
});

interface IconButtonLoadingProps {
  loading?: boolean;
  loadingText?: React.ReactNode;
}

export interface ButtonProps
  extends StyledIconButtonProps,
    IconButtonLoadingProps {}

export const IconButton = forwardRef<HTMLButtonElement, ButtonProps>(
  (props, ref) => {
    const { loading, disabled, loadingText, children, ...rest } = props;

    const trulyDisabled = loading || disabled;

    return (
      <StyledIconButton disabled={trulyDisabled} ref={ref} {...rest}>
        {loading && !loadingText ? (
          <>
            <Spinner />
          </>
        ) : loadingText ? (
          loadingText
        ) : (
          children
        )}
      </StyledIconButton>
    );
  },
);

IconButton.displayName = "IconButton";
