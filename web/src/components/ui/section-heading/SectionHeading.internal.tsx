import { styled } from "@/styled-system/jsx";

export const SectionHeading = styled(
  "h2",
  {
    base: {
      color: "text.muted",
      fontSize: "md",
      fontWeight: "bold",
      lineHeight: "normal",
    },
  },
  {
    defaultProps: {
      className: "section-heading",
    },
  },
);
