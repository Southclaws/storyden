import { ark } from "@ark-ui/react";
import {
  type ComponentPropsWithoutRef,
  ForwardedRef,
  type PropsWithChildren,
  forwardRef,
} from "react";
import { type TitleInputVariantProps, titleInput } from "styled-system/recipes";

import { cx } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

export type TitleInputProps = TitleInputVariantProps &
  ComponentPropsWithoutRef<typeof ark.input>;

function _TitleInput(
  props: PropsWithChildren<TitleInputProps>,
  ref: ForwardedRef<HTMLSpanElement>,
) {
  const { children, ...rest } = props;
  const [recipeProps, componentProps] = titleInput.splitVariantProps(rest);

  return (
    <styled.span
      contentEditable
      suppressContentEditableWarning
      suppressHydrationWarning
      className={cx(titleInput({ ...recipeProps }))}
      {...(componentProps as any)}
      ref={ref}
    >
      {children}
    </styled.span>
  );
}

const TitleInput = forwardRef(_TitleInput);

export { TitleInput };
