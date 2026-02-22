import { test } from "uvu";
import * as assert from "uvu/assert";

import {
  DefaultRoleMetadata,
  parseRoleMetadata,
  writeRoleMetadata,
} from "./metadata";

test("parseRoleMetadata keeps valid metadata", () => {
  const parsed = parseRoleMetadata({
    bold: true,
    italic: false,
    coloured: false,
  });

  assert.equal(parsed, {
    bold: true,
    italic: false,
    coloured: false,
  });
});

test("parseRoleMetadata falls back to defaults for invalid metadata", () => {
  const originalWarn = console.warn;
  let warned = false;

  console.warn = () => {
    warned = true;
  };

  let parsed;
  try {
    parsed = parseRoleMetadata({
      bold: "yes",
      italic: 123,
      coloured: "no",
    } as unknown as Record<string, unknown>);
  } finally {
    console.warn = originalWarn;
  }

  assert.ok(warned);
  assert.equal(parsed, DefaultRoleMetadata);
});

test("writeRoleMetadata preserves unknown keys", () => {
  const written = writeRoleMetadata(
    {
      rainbow: true,
    },
    {
      bold: true,
      italic: false,
      coloured: false,
    },
  );

  assert.equal(written, {
    rainbow: true,
    bold: true,
    italic: false,
    coloured: false,
  });
});

test.run();
