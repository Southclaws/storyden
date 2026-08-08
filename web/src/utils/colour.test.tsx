import { getContrast } from "polished";
import { describe, expect, test } from "vitest";

import { getReadableTextColour } from "./colour";

describe("getReadableTextColour", () => {
  test.each(["#b2ffd8", "#fbdfff", "#ffd3d8", "#202020", "#777777"])(
    "returns an accessible foreground for %s",
    (background) => {
      const foreground = getReadableTextColour(background);

      expect(getContrast(background, foreground)).toBeGreaterThanOrEqual(4.5);
    },
  );

  test("falls back safely for invalid colours", () => {
    expect(getReadableTextColour("not-a-colour")).toBe("#303030");
  });
});
