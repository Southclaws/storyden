import { ark } from "@ark-ui/react/factory";
import type { ComponentProps } from "react";
import { forwardRef, useMemo } from "react";

import { styled } from "@/styled-system/jsx";
import { ButtonVariantProps, button } from "@/styled-system/recipes";

import { useButtonPropsContext } from "./button";
import { Spinner } from "./Spinner";

export type StyledIconButtonProps = ComponentProps<typeof StyledIconButton>;
export const StyledIconButton = styled(ark.button, button, {
  defaultProps: { px: "0" } as ButtonVariantProps,
});

interface IconButtonLoadingProps {
  loading?: boolean;
  loadingText?: React.ReactNode;
}

export interface IconButtonProps
  extends StyledIconButtonProps,
    IconButtonLoadingProps {}

export const IconButton = forwardRef<HTMLButtonElement, IconButtonProps>(
  (props, ref) => {
    const propsContext = useButtonPropsContext();
    const iconButtonProps = useMemo(
      () => ({ ...propsContext, ...props }),
      [propsContext, props],
    );
    const { loading, disabled, loadingText, children, ...rest } =
      iconButtonProps;

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
