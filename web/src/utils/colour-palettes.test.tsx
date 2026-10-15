import { describe, expect, test } from "vitest";

import { badgeColours as roleColours } from "@/components/role/colours";
import { badgeColours } from "@/components/ui/badge/Badge.internal";
import { categoryColours } from "@/lib/category/colours";

import { deriveColour, getColourAsHex } from "./colour";

describe("colour palettes", () => {
  // Fixed outputs from the prior Chroma entry point, including alpha and contrast.
  test.each([
    {
      input: "#27b981",
      category: {
        background: "rgb(148 229 189)",
        border: "rgb(99 179 141 / 0.8)",
        foreground: "#303030",
      },
      role: {
        background: "rgb(236 236 236)",
        border: "rgb(70 130 97 / 0.2)",
        foreground: "#303030",
      },
      badge: {
        background: "rgb(178 255 216)",
        border: "rgb(0 244 160)",
        foreground: "#303030",
      },
    },
    {
      input: "#123456",
      category: {
        background: "rgb(87 92 102)",
        border: "rgb(46 51 59 / 0.8)",
        foreground: "#E8ECEA",
      },
      role: {
        background: "rgb(137 137 137)",
        border: "rgb(14 14 23 / 0.2)",
        foreground: "#000",
      },
      badge: {
        background: "rgb(215 235 255)",
        border: "rgb(163 193 255)",
        foreground: "#303030",
      },
    },
    {
      input: "rgb(255 0 128 / 0.5)",
      category: {
        background: "rgb(231 115 170 / 0.5)",
        border: "rgb(232 68 129 / 0.8)",
        foreground: "#303030",
      },
      role: {
        background: "rgb(185 185 185 / 0.5)",
        border: "rgb(175 55 88 / 0.2)",
        foreground: "#303030",
      },
      badge: {
        background: "rgb(255 211 232)",
        border: "rgb(255 142 194)",
        foreground: "#303030",
      },
    },
  ])(
    "preserves category, role and badge colours for $input",
    ({ input, category, role, badge }) => {
      expect(categoryColours(input)).toEqual(category);
      expect(roleColours(input)).toEqual(role);
      expect(badgeColours(input)).toEqual(badge);
    },
  );

  test.each([
    ["Storyden", "#b194e5"],
    ["🌱", "#4bb1e5"],
    ["alice", "#c58dd6"],
  ])("keeps deterministic colour for %s", (input, expected) => {
    expect(deriveColour(input)).toBe(expected);
  });

  test("preserves invalid-colour fallbacks", () => {
    expect(getColourAsHex("invalid")).toBe("#27b981");
    expect(roleColours("invalid")).toEqual(roleColours("#22c55e"));
  });
});
