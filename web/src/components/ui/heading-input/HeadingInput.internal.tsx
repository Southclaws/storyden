import { ark } from "@ark-ui/react";
import {
  ClipboardEvent,
  type ComponentPropsWithoutRef,
  FormEvent,
  ForwardedRef,
  KeyboardEvent,
  forwardRef,
  useCallback,
  useEffect,
  useImperativeHandle,
  useRef,
} from "react";
import {
  type HeadingInputVariantProps,
  type TypographyHeadingVariantProps,
  headingInput,
} from "styled-system/recipes";

import { cx } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";
import { typographyHeading } from "@/styled-system/recipes";
import { JsxStyleProps } from "@/styled-system/types";

type CustomProps = {
  onValueChange: (s: string) => void;
};

export type HeadingInputProps = JsxStyleProps &
  HeadingInputVariantProps &
  TypographyHeadingVariantProps &
  ComponentPropsWithoutRef<typeof ark.input> &
  CustomProps;

function HeadingInputWithRef(
  props: HeadingInputProps,
  ref: ForwardedRef<HTMLSpanElement>,
) {
  const { onValueChange, defaultValue, value, ...rest } = props;
  const [recipeProps, componentProps] = headingInput.splitVariantProps(rest);
  const internalRef = useRef<HTMLSpanElement>(null);

  useImperativeHandle(ref, () => internalRef.current as any);

  useEffect(() => {
    if (internalRef.current && value !== undefined) {
      if (internalRef.current.textContent !== value) {
        internalRef.current.textContent = value.toString();
      }
    }
  }, [value]);

  useEffect(() => {
    if (internalRef.current && defaultValue) {
      internalRef.current.textContent = defaultValue.toString();
    }
  }, [defaultValue]);

  const [headingProps] = typographyHeading.splitVariantProps(rest);

  const handleInput = useCallback(
    (e: FormEvent<HTMLSpanElement>) => {
      const text = (e.target as any).textContent;
      onValueChange(text);
    },
    [onValueChange],
  );

  const handleKeyDown = useCallback((e: KeyboardEvent<HTMLSpanElement>) => {
    if (e.code === "Enter") {
      e.preventDefault();
      e.stopPropagation();
    }
  }, []);

  const handlePaste = useCallback((e: ClipboardEvent<HTMLSpanElement>) => {
    e.preventDefault();

    const text = e.clipboardData.getData("text/plain");

    const stripped = text.replace(/(\r\n|\n|\r)/gm, " ");

    document.execCommand("insertText", false, stripped);
  }, []);

  return (
    <styled.span
      {...(componentProps as any)}
      ref={internalRef}
      className={cx(
        headingInput({ ...recipeProps }),
        typographyHeading({ ...headingProps }),
      )}
      //
      // NOTE: We're doing a bit of a hack here in order to make this
      // field look nice and behave like the Substack title editor.
      //
      // More info:
      //
      // https://medium.com/programming-essentials/good-to-know-about-the-state-management-of-a-contenteditable-element-in-react-adb4f933df12
      //
      contentEditable
      suppressContentEditableWarning
      suppressHydrationWarning
      spellCheck={false}
      onInput={handleInput}
      onKeyDown={handleKeyDown}
      onPaste={handlePaste}
    >
      {defaultValue}
    </styled.span>
  );
}

const HeadingInput = forwardRef(HeadingInputWithRef);

export { HeadingInput };
