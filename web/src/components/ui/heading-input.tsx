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
    const element = internalRef.current;
    if (element && value !== undefined) {
      const currentText = element.textContent || '';
      const newValue = value.toString();
      
      // Only update if content differs and element is not focused
      // This prevents cursor jumping during user input
      if (currentText !== newValue && document.activeElement !== element) {
        // Save current selection/cursor position
        const selection = window.getSelection();
        const range = selection && selection.rangeCount > 0 ? selection.getRangeAt(0) : null;
        const isSelectionInElement = range && element.contains(range.commonAncestorContainer);
        
        // Calculate cursor position relative to text length
        let cursorPosition = 0;
        if (isSelectionInElement && range) {
          const preCaretRange = range.cloneRange();
          preCaretRange.selectNodeContents(element);
          preCaretRange.setEnd(range.endContainer, range.endOffset);
          cursorPosition = preCaretRange.toString().length;
        }
        
        // Update content
        element.textContent = newValue;
        
        // Restore cursor position if element was previously focused
        if (isSelectionInElement && cursorPosition <= newValue.length) {
          const newRange = document.createRange();
          const walker = document.createTreeWalker(
            element,
            NodeFilter.SHOW_TEXT,
            null
          );
          
          let currentPos = 0;
          let textNode = walker.nextNode();
          
          while (textNode && currentPos + textNode.textContent!.length < cursorPosition) {
            currentPos += textNode.textContent!.length;
            textNode = walker.nextNode();
          }
          
          if (textNode) {
            const offset = cursorPosition - currentPos;
            newRange.setStart(textNode, Math.min(offset, textNode.textContent!.length));
            newRange.setEnd(textNode, Math.min(offset, textNode.textContent!.length));
            
            if (selection) {
              selection.removeAllRanges();
              selection.addRange(newRange);
            }
          }
        }
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

    const selection = window.getSelection();
    if (selection && selection.rangeCount > 0) {
      const range = selection.getRangeAt(0);
      range.deleteContents();
      
      const textNode = document.createTextNode(stripped);
      range.insertNode(textNode);
      
      // Move cursor to end of inserted text
      range.setStartAfter(textNode);
      range.setEndAfter(textNode);
      selection.removeAllRanges();
      selection.addRange(range);
      
      // Trigger input event for React to detect the change
      const inputEvent = new InputEvent('input', {
        bubbles: true,
        cancelable: true,
        inputType: 'insertText',
        data: stripped
      });
      e.currentTarget.dispatchEvent(inputEvent);
    }
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
