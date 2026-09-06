import type { ThemeAsset } from "@/api/openapi-schema";

export type ThemeEditorKind = "stylesheet" | "script";

export type ThemeEditorDocument = {
  key: string;
  kind: ThemeEditorKind;
  label: string;
  source: string;
  savedSource: string;
  asset?: ThemeAsset;
};

export function themeDocumentsSignature(documents: ThemeEditorDocument[]) {
  return documents
    .map(
      ({ kind, source, asset }) =>
        `${kind}\u0000${asset?.id ?? "new"}\u0000${source}`,
    )
    .join("\u0001");
}

export function themeScriptSignature(documents: ThemeEditorDocument[]) {
  return documents
    .filter(({ kind }) => kind === "script")
    .map(({ key, source, asset }) => `${asset?.id ?? key}\u0000${source}`)
    .join("\u0000");
}

export function normaliseThemeDocumentOrder(documents: ThemeEditorDocument[]) {
  return [
    ...documents.filter(({ kind }) => kind === "stylesheet"),
    ...documents.filter(({ kind }) => kind === "script"),
  ];
}

export function themeDocumentBytes(source: string) {
  return new TextEncoder().encode(source).byteLength;
}
