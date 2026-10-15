import { describe, expect, test } from "vitest";

import { contrast, oklch, parseColour } from "./colour-core";

describe("colour conversions", () => {
  test.each([
    ["#f08", "#ff0088"],
    ["00ff88", "#00ff88"],
    ["#ff008080", "#ff008080"],
    ["red", "#ff0000"],
    ["rebeccapurple", "#663399"],
    ["rgb(255 0 128 / 0.5)", "#ff008080"],
    ["hsl(120deg 100% 50%)", "#00ff00"],
    ["oklch(70% 0.1 40)", "#d4896e"],
    ["transparent", "#00000000"],
  ])("preserves %s", (input, expected) => {
    expect(parseColour(input).hex()).toBe(expected);
  });

  test("rejects invalid input for callers to handle", () => {
    expect(() => parseColour("not-a-colour")).toThrow();
    expect(() => parseColour("")).toThrow();
  });

  test("keeps standard contrast ratios and numeric OKLCH conversion", () => {
    expect(contrast("black", "white")).toBe(21);
    expect(contrast("white", "white")).toBe(1);
    expect(oklch(1, 0, 0).hex()).toBe("#ffffff");
    expect(oklch(0, 0, 0).hex()).toBe("#000000");
  });
});
