"use client";

import { useEffect, useMemo, useState } from "react";
import { toast } from "sonner";

import { useAccountGet } from "@/api/openapi-client/accounts";
import {
  adminThemeAssetDelete,
  adminThemeAssetUpload,
  adminThemeUpdate,
  useAdminThemeGet,
} from "@/api/openapi-client/admin";
import { Permission, type ThemeAsset } from "@/api/openapi-schema";
import { getAPIAddress } from "@/config";
import { resolveThemeAssetHref } from "@/lib/theme/manifest";
import { revalidateTheme } from "@/lib/theme/revalidate-theme";
import {
  setThemeEditingEnabled,
  useThemeEditingEnabled,
} from "@/lib/theme/theme-editing";
import { hasPermission } from "@/utils/permissions";

import { ThemeEditorPanel } from "./ThemeEditorPanel";
import {
  type ThemeEditorDocument,
  type ThemeEditorKind,
  normaliseThemeDocumentOrder,
  themeDocumentBytes,
  themeDocumentsSignature,
  themeScriptSignature,
} from "./theme-editor-model";

const MAX_DOCUMENTS = 32;
const MAX_DOCUMENT_BYTES = 1024 * 1024;
const MAX_THEME_BYTES = 5 * 1024 * 1024;

export function ThemeEditorHost() {
  const editing = useThemeEditingEnabled();
  const account = useAccountGet();

  if (!editing || !hasPermission(account.data, Permission.ADMINISTRATOR)) {
    return null;
  }

  return <ThemeEditorController />;
}

function ThemeEditorController() {
  const query = useAdminThemeGet();
  const [documents, setDocuments] = useState<ThemeEditorDocument[] | null>(
    null,
  );
  const [baseline, setBaseline] = useState("");
  const [baselineScripts, setBaselineScripts] = useState("");
  const [selectedKey, setSelectedKey] = useState<string>();
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string>();
  const panelRect = useMemo(getPanelRect, []);

  useEffect(() => {
    if (!query.data || documents !== null) return;

    let cancelled = false;
    void loadDocuments(query.data.stylesheets, query.data.scripts)
      .then((loaded) => {
        if (cancelled) return;
        setDocuments(loaded);
        setBaseline(themeDocumentsSignature(loaded));
        setBaselineScripts(themeScriptSignature(loaded));
        setSelectedKey(loaded[0]?.key);
      })
      .catch((cause) => {
        if (!cancelled) setError(errorMessage(cause));
      });
    return () => {
      cancelled = true;
    };
  }, [documents, query.data]);

  const dirty =
    documents !== null && themeDocumentsSignature(documents) !== baseline;

  function addDocument(kind: ThemeEditorKind) {
    const next: ThemeEditorDocument = {
      key: `new-${crypto.randomUUID()}`,
      kind,
      label: nextDocumentLabel(documents ?? [], kind),
      source:
        kind === "stylesheet"
          ? "/* Add theme CSS here */\n"
          : "// Add theme JavaScript here\n",
      savedSource: "",
    };
    setDocuments((current) =>
      normaliseThemeDocumentOrder([...(current ?? []), next]),
    );
    setSelectedKey(next.key);
  }

  function updateDocument(key: string, source: string) {
    setDocuments(
      (current) =>
        current?.map((item) =>
          item.key === key ? { ...item, source } : item,
        ) ?? null,
    );
  }

  function deleteDocument(key: string) {
    setDocuments((current) => {
      if (!current) return current;
      const index = current.findIndex((item) => item.key === key);
      const next = current.filter((item) => item.key !== key);
      if (selectedKey === key) {
        setSelectedKey(next[Math.min(index, next.length - 1)]?.key);
      }
      return next;
    });
  }

  function moveDocument(key: string, offset: -1 | 1) {
    setDocuments((current) => {
      if (!current) return current;
      const index = current.findIndex((item) => item.key === key);
      const item = current[index];
      if (!item) return current;
      const peerIndices = current.flatMap((candidate, candidateIndex) =>
        candidate.kind === item.kind ? [candidateIndex] : [],
      );
      const peerIndex = peerIndices.indexOf(index);
      const target = peerIndices[peerIndex + offset];
      if (target === undefined) return current;
      const next = [...current];
      [next[index], next[target]] = [next[target]!, next[index]!];
      return next;
    });
  }

  async function save() {
    if (!documents || !dirty) return;

    setError(undefined);
    const validationError = validateDocuments(documents);
    if (validationError) {
      setError(validationError);
      return;
    }

    const scriptsChanged = themeScriptSignature(documents) !== baselineScripts;
    if (
      scriptsChanged &&
      !window.confirm(
        "Save trusted JavaScript? Theme scripts execute with Storyden's browser privileges for every visitor, including administrators.",
      )
    ) {
      return;
    }

    setBusy(true);
    const uploaded: ThemeAsset[] = [];
    let published = false;
    try {
      const publishDocuments: Array<{
        document: ThemeEditorDocument;
        asset: ThemeAsset;
      }> = [];

      for (const [index, document] of documents.entries()) {
        if (document.asset && document.source === document.savedSource) {
          publishDocuments.push({ document, asset: document.asset });
          continue;
        }

        const extension = document.kind === "stylesheet" ? "css" : "js";
        const mimeType =
          document.kind === "stylesheet"
            ? "text/css"
            : "application/javascript";
        const filename = `theme-editor-${Date.now()}-${index + 1}.${extension}`;
        const file = new File([document.source], filename, { type: mimeType });
        const asset = await adminThemeAssetUpload(file, { filename });
        uploaded.push(asset);
        publishDocuments.push({ document, asset });
      }

      const result = await adminThemeUpdate({
        stylesheets: publishDocuments
          .filter(({ document }) => document.kind === "stylesheet")
          .map(({ asset }) => asset.id),
        scripts: publishDocuments
          .filter(({ document }) => document.kind === "script")
          .map(({ asset }) => asset.id),
      });
      published = true;

      const nextDocuments = remapPublishedDocuments(documents, result);
      setDocuments(nextDocuments);
      setBaseline(themeDocumentsSignature(nextDocuments));
      setBaselineScripts(themeScriptSignature(nextDocuments));
      await Promise.all([
        query.mutate(result, { revalidate: false }),
        revalidateTheme(),
      ]);

      if (scriptsChanged) {
        window.location.reload();
        return;
      }

      applyLiveStyles(nextDocuments, result.revision);
      toast.success("Live theme saved");
    } catch (cause) {
      if (!published) {
        await Promise.allSettled(
          uploaded.map(({ filename }) => adminThemeAssetDelete(filename)),
        );
      }
      setError(errorMessage(cause));
    } finally {
      setBusy(false);
    }
  }

  return (
    <ThemeEditorPanel
      documents={documents}
      selectedKey={selectedKey}
      dirty={dirty}
      busy={busy}
      runtimeDisabled={query.data?.runtime_disabled ?? false}
      error={error ?? (query.error ? errorMessage(query.error) : undefined)}
      panelSize={panelRect.size}
      panelPosition={panelRect.position}
      onSelect={setSelectedKey}
      onAdd={addDocument}
      onChange={updateDocument}
      onMove={moveDocument}
      onDelete={deleteDocument}
      onSave={save}
      onExit={() => setThemeEditingEnabled(false)}
    />
  );
}

async function loadDocuments(stylesheets: ThemeAsset[], scripts: ThemeAsset[]) {
  const load = async (
    asset: ThemeAsset,
    kind: ThemeEditorKind,
    index: number,
  ): Promise<ThemeEditorDocument> => {
    const href = resolveThemeAssetHref(asset.path, getAPIAddress());
    if (!href) {
      throw new Error(`Could not load ${asset.filename}: invalid asset URL.`);
    }
    const response = await fetch(href, { credentials: "omit" });
    if (!response.ok) {
      throw new Error(`Could not load ${asset.filename} (${response.status}).`);
    }
    const source = await response.text();
    return {
      key: asset.id,
      kind,
      label: `${kind === "stylesheet" ? "CSS" : "JS"} ${index + 1}`,
      source,
      savedSource: source,
      asset,
    };
  };

  return Promise.all([
    ...stylesheets.map((asset, index) => load(asset, "stylesheet", index)),
    ...scripts.map((asset, index) => load(asset, "script", index)),
  ]);
}

function validateDocuments(documents: ThemeEditorDocument[]) {
  if (documents.length > MAX_DOCUMENTS) {
    return `A theme may contain at most ${MAX_DOCUMENTS} CSS and JavaScript tabs.`;
  }
  let total = 0;
  for (const document of documents) {
    const size = themeDocumentBytes(document.source);
    if (size === 0) {
      return `${document.label} is empty. Delete the tab instead of saving an empty asset.`;
    }
    if (size > MAX_DOCUMENT_BYTES) {
      return `${document.label} exceeds the 1 MiB asset limit.`;
    }
    total += size;
  }
  if (total > MAX_THEME_BYTES) {
    return "The active theme exceeds the 5 MiB total limit.";
  }
  return undefined;
}

function remapPublishedDocuments(
  documents: ThemeEditorDocument[],
  result: {
    stylesheets: ThemeAsset[];
    scripts: ThemeAsset[];
  },
) {
  let stylesheet = 0;
  let script = 0;
  return documents.map((document) => {
    const asset =
      document.kind === "stylesheet"
        ? result.stylesheets[stylesheet++]
        : result.scripts[script++];
    if (!asset) throw new Error("Published theme response was incomplete.");
    return {
      ...document,
      key: asset.id,
      source: document.source,
      savedSource: document.source,
      asset,
    };
  });
}

function applyLiveStyles(documents: ThemeEditorDocument[], revision: string) {
  document
    .querySelectorAll(
      'link[data-sd-theme-asset="stylesheet"], style[data-sd-theme-editor-live]',
    )
    .forEach((element) => element.remove());

  const stylesheets = documents.filter(({ kind }) => kind === "stylesheet");
  if (stylesheets.length > 0) {
    const style = document.createElement("style");
    style.dataset["sdThemeEditorLive"] = "";
    style.textContent = stylesheets.map(({ source }) => source).join("\n");
    document.head.appendChild(style);
  }
  document.documentElement.dataset["sdThemeRevision"] = revision;
}

function nextDocumentLabel(
  documents: ThemeEditorDocument[],
  kind: ThemeEditorKind,
) {
  const prefix = kind === "stylesheet" ? "CSS" : "JS";
  const used = new Set(
    documents.filter((item) => item.kind === kind).map(({ label }) => label),
  );
  let index = 1;
  while (used.has(`${prefix} ${index}`)) index += 1;
  return `${prefix} ${index}`;
}

function getPanelRect() {
  if (typeof window === "undefined") {
    return {
      size: { width: 540, height: 560 },
      position: { x: 16, y: 16 },
    };
  }
  const size = {
    width: Math.max(320, Math.min(560, window.innerWidth - 24)),
    height: Math.max(300, Math.min(620, window.innerHeight - 32)),
  };
  return {
    size,
    position: {
      x: Math.max(12, window.innerWidth - size.width - 16),
      y: 16,
    },
  };
}

function errorMessage(cause: unknown) {
  return cause instanceof Error ? cause.message : String(cause);
}
