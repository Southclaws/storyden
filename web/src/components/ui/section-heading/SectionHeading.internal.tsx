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
    variants: {
      emphasis: {
        strong: {
          color: "text.default",
          fontSize: "lg",
          fontWeight: "semibold",
        },
        subtle: {},
      },
    },
    defaultVariants: {
      emphasis: "subtle",
    },
  },
  {
    defaultProps: {
      className: "section-heading",
    },
  },
);
