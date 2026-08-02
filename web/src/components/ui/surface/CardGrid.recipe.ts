import { defineSlotRecipe } from "@pandacss/dev";

export const cardGrid = defineSlotRecipe({
  className: "card-grid",
  slots: ["container", "grid"],
  base: {
    container: {
      containerType: "inline-size",
      containerName: "card-grid",
      width: "full",
    },
    grid: {
      display: "grid",
      width: "full",
      gap: "2",
      gridTemplateColumns: {
        base: "minmax(0, 1fr)",
        _containerMedium: "repeat(2, minmax(0, 1fr))",
        _containerLarge: "repeat(6, minmax(0, 1fr))",
      },
      "& > *": {
        minWidth: "0",
        _containerLarge: {
          gridColumn: "span 2",
        },
      },
      '&[data-balance="one"] > *': {
        _containerMedium: {
          gridColumn: "1 / -1",
        },
        _containerLarge: {
          gridColumn: "1 / -1",
        },
      },
      '&[data-balance="last-two"] > :nth-last-child(-n+2)': {
        _containerLarge: {
          gridColumn: "span 3",
        },
      },
      '&[data-balance="last-four"] > :nth-last-child(-n+4)': {
        _containerLarge: {
          gridColumn: "span 3",
        },
      },
    },
  },
});
