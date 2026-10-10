import assert from "node:assert/strict";
import { describe, test } from "vitest";

import {
  type ThemeEditorDocument,
  normaliseThemeDocumentOrder,
  themeDocumentBytes,
  themeDocumentsSignature,
  themeScriptSignature,
} from "./theme-editor-model";

describe("theme editor model", () => {
  const css = document("css", "stylesheet", "body{}", "css-id");
  const script = document("js", "script", "window.ready=true", "js-id");

  test("normalises stylesheets before scripts without changing peer order", () => {
    assert.deepStrictEqual(
      normaliseThemeDocumentOrder([script, { ...css, key: "css-2" }, css]).map(
        ({ key }) => key,
      ),
      ["css-2", "css", "js"],
    );
  });

  test("signatures detect source, identity, and script ordering changes", () => {
    assert.notDeepStrictEqual(
      themeDocumentsSignature([css]),
      themeDocumentsSignature([{ ...css, source: "body{color:red}" }]),
    );
    assert.ok(
      themeScriptSignature([css, script]).endsWith("window.ready=true"),
    );
    assert.notDeepStrictEqual(
      themeScriptSignature([script, { ...script, key: "js-2", source: "b()" }]),
      themeScriptSignature([{ ...script, key: "js-2", source: "b()" }, script]),
    );
  });

  test("counts encoded UTF-8 bytes rather than JavaScript characters", () => {
    assert.strictEqual(themeDocumentBytes("a😀"), 5);
  });
});

function document(
  key: string,
  kind: ThemeEditorDocument["kind"],
  source: string,
  id: string,
): ThemeEditorDocument {
  return {
    key,
    kind,
    label: key,
    source,
    savedSource: source,
    asset: {
      id: id.padEnd(20, "0"),
      filename: key,
      integrity: "sha256-dGVzdA==",
      mime_type: kind === "stylesheet" ? "text/css" : "application/javascript",
      path: `/api/info/theme/assets/${key}`,
      size: source.length,
    },
  };
}
