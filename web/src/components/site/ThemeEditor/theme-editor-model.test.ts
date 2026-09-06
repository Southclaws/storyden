import { suite } from "uvu";
import * as assert from "uvu/assert";

import {
  type ThemeEditorDocument,
  normaliseThemeDocumentOrder,
  themeDocumentBytes,
  themeDocumentsSignature,
  themeScriptSignature,
} from "./theme-editor-model";

const test = suite("theme editor model");

const css = document("css", "stylesheet", "body{}", "css-id");
const script = document("js", "script", "window.ready=true", "js-id");

test("normalises stylesheets before scripts without changing peer order", () => {
  assert.equal(
    normaliseThemeDocumentOrder([script, { ...css, key: "css-2" }, css]).map(
      ({ key }) => key,
    ),
    ["css-2", "css", "js"],
  );
});

test("signatures detect source, identity, and script ordering changes", () => {
  assert.not.equal(
    themeDocumentsSignature([css]),
    themeDocumentsSignature([{ ...css, source: "body{color:red}" }]),
  );
  assert.ok(themeScriptSignature([css, script]).endsWith("window.ready=true"));
  assert.not.equal(
    themeScriptSignature([script, { ...script, key: "js-2", source: "b()" }]),
    themeScriptSignature([{ ...script, key: "js-2", source: "b()" }, script]),
  );
});

test("counts encoded UTF-8 bytes rather than JavaScript characters", () => {
  assert.is(themeDocumentBytes("a😀"), 5);
});

test.run();

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
