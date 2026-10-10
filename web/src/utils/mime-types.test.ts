import assert from "node:assert/strict";
import { test } from "vitest";

import {
  getExtensionsForMimeType,
  getExtensionsForMimeTypes,
} from "./mime-types";

test("getExtensionsForMimeType returns mapped extensions", () => {
  assert.deepStrictEqual(getExtensionsForMimeType("image/jpeg"), [
    "jpg",
    "jpeg",
  ]);
});

test("getExtensionsForMimeType returns empty for unknown types", () => {
  assert.deepStrictEqual(getExtensionsForMimeType("application/unknown"), []);
});

test("getExtensionsForMimeTypes expands wildcard groups", () => {
  const result = getExtensionsForMimeTypes(["image/*"]);
  assert.ok(result.includes("jpg"));
  assert.ok(result.includes("png"));
  assert.ok(result.includes("svg"));
});

test("getExtensionsForMimeTypes deduplicates across explicit and wildcard", () => {
  const result = getExtensionsForMimeTypes(["image/jpeg", "image/*"]);
  const jpgCount = result.filter((ext) => ext === "jpg").length;
  const jpegCount = result.filter((ext) => ext === "jpeg").length;
  assert.strictEqual(jpgCount, 1);
  assert.strictEqual(jpegCount, 1);
});

test("getExtensionsForMimeTypes supports mixed specific mime types", () => {
  assert.deepStrictEqual(
    getExtensionsForMimeTypes(["application/pdf", "text/plain"]),
    ["pdf", "txt"],
  );
});

test("getExtensionsForMimeTypes returns empty for unknown wildcard groups", () => {
  assert.deepStrictEqual(getExtensionsForMimeTypes(["unknown/*"]), []);
});
