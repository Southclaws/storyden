import { ark } from "@ark-ui/react/factory";
import type { ComponentProps } from "react";
import { forwardRef, useMemo } from "react";

import { styled } from "@/styled-system/jsx";
import { ButtonVariantProps, button } from "@/styled-system/recipes";

import { useButtonPropsContext } from "../button";
import { Spinner } from "../spinner";

export type StyledIconButtonProps = ComponentProps<typeof StyledIconButton>;
export const StyledIconButton = styled(ark.button, button, {
  defaultProps: { px: "0" } as ButtonVariantProps,
});

interface IconButtonLoadingProps {
  loading?: boolean;
  loadingText?: React.ReactNode;
}

export interface IconButtonProps
  extends Omit<StyledIconButtonProps, "colorPalette">, IconButtonLoadingProps {}

export const IconButton = forwardRef<HTMLButtonElement, IconButtonProps>(
  (props, ref) => {
    const propsContext = useButtonPropsContext();
    const iconButtonProps = useMemo(
      () => ({ ...propsContext, ...props }),
      [propsContext, props],
    );
    const {
      loading,
      disabled,
      loadingText,
      children,
      "aria-label": ariaLabel,
      title = ariaLabel,
      ...rest
    } = iconButtonProps;

    const trulyDisabled = loading || disabled;

    return (
      <StyledIconButton
        aria-label={ariaLabel}
        disabled={trulyDisabled}
        ref={ref}
        title={title}
        {...rest}
      >
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
