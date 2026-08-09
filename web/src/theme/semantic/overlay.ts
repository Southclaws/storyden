export const overlayBaseStyles = {
  background: "background.overlay",
  backdropBlur: "subtle",
  backdropFilter: "auto",
  borderColor: "border.strong",
  borderWidth: "thin",
  color: "text.default",
} as const;

export const overlayContentStyles = {
  ...overlayBaseStyles,
  boxShadow: "overlay",
} as const;

export const floatingOverlayStyles = {
  ...overlayBaseStyles,
  boxShadow: "floating",
} as const;
