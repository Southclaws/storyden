import { styled } from "@/styled-system/jsx";

export const PageHeading = styled(
  "h1",
  {
    base: {
      color: "text.default",
      fontSize: "lg",
      fontWeight: "semibold",
      lineHeight: "normal",
    },
  },
  {
    defaultProps: {
      className: "page-heading",
    },
  },
);
