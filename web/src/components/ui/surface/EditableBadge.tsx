import { PropsWithChildren } from "react";

import { styled } from "@/styled-system/jsx";

import { Badge, BadgeProps } from "../badge";

type Props = Omit<BadgeProps, "children"> & {
  editing: boolean;
  value: string;
  inputLabel: string;
  placeholder?: string;
  onValueChange: (value: string) => void;
};

// Keep the Badge shell in both modes. The hidden value sizes the native editor
// without an arbitrary field width or a second copy of the Badge recipe.
export function EditableBadge({
  editing,
  value,
  inputLabel,
  placeholder,
  onValueChange,
  children,
  ...props
}: PropsWithChildren<Props>) {
  return (
    <Badge
      {...props}
      position="relative"
      zIndex={editing ? "[5]" : undefined}
      minW="0"
      maxW="full"
      whiteSpace="nowrap"
    >
      <styled.span
        aria-hidden={editing || undefined}
        visibility={editing ? "hidden" : undefined}
        minW="0"
        truncate
      >
        {editing ? value || placeholder || "\u00a0" : (children ?? value)}
      </styled.span>
      {editing && (
        <styled.input
          type="text"
          aria-label={inputLabel}
          value={value}
          placeholder={placeholder}
          position="absolute"
          inset="0"
          w="full"
          h="full"
          minW="0"
          minH="0"
          appearance="none"
          background="transparent"
          border="none"
          borderRadius="inherit"
          padding="inherit"
          font="inherit"
          color="inherit"
          textOverflow="ellipsis"
          overflow="hidden"
          userSelect="text"
          cursor="text"
          outline="none"
          _hover={{
            outlineWidth: "thin",
            outlineStyle: "solid",
            outlineColor: "border.strong",
          }}
          _focusVisible={{
            outlineWidth: "medium",
            outlineStyle: "solid",
            outlineColor: "accent.text",
            outlineOffset: "0.5",
          }}
          _placeholder={{ color: "text.muted" }}
          onChange={(event) =>
            onValueChange(event.target.value.replace(/[\r\n]+/g, " "))
          }
          onPaste={(event) => {
            const text = event.clipboardData.getData("text");

            if (!/[\r\n]/.test(text)) {
              return;
            }

            event.preventDefault();
            const { selectionStart, selectionEnd } = event.currentTarget;
            onValueChange(
              value.slice(0, selectionStart ?? 0) +
                text.replace(/[\r\n]+/g, " ") +
                value.slice(selectionEnd ?? value.length),
            );
          }}
          onKeyDown={(event) => {
            if (event.key === "Enter" && !event.nativeEvent.isComposing) {
              event.preventDefault();
              event.currentTarget.blur();
            }
          }}
        />
      )}
    </Badge>
  );
}
