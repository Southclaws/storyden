import { defineRecipe } from "@pandacss/dev";

export const skeleton = defineRecipe({
  className: "skeleton",
  base: {
    display: "inline-block",
    height: "1em",
    position: "relative",
    overflow: "hidden",
    backgroundColor: "#DDDBDD",
    _after: {
      animation: "shimmer 5s infinite",
      content: "''",
      position: "absolute",
      top: "0",
      right: "0",
      bottom: "0",
      left: "0",
      transform: "translateX(-100%)",
      backgroundImage:
        "linear-gradient(90deg, rgba(255, 255, 255, 0) 0, rgba(255, 255, 255, 0.2) 25%, rgba(255, 255, 255, 0.5) 66%, rgba(255, 255, 255, 0))",
    },
  },
});
