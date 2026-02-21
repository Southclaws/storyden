import { useEffect, useMemo, useState } from "react";
import { useSWRConfig } from "swr";
import { parse, stringify } from "yaml";

import { handle } from "@/api/client";
import { mutateTransaction } from "@/api/mutate";
import {
  getPluginGetKey,
  usePluginUpdateManifest,
} from "@/api/openapi-client/plugins";
import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { Box, LStack, WStack, styled } from "@/styled-system/jsx";

type Props = {
  pluginID: string;
  manifest: Record<string, unknown>;
  editable: boolean;
};

type ManifestError = {
  overview: string | null;
  details: string;
};

export function ManifestTab({ pluginID, manifest, editable }: Props) {
  const { mutate } = useSWRConfig();
  const { trigger: updateManifest } = usePluginUpdateManifest(pluginID);

  const initialYAML = useMemo(() => toYAML(manifest), [manifest]);
  const [yaml, setYAML] = useState(initialYAML);
  const [dirty, setDirty] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<ManifestError | null>(null);

  useEffect(() => {
    if (!dirty) {
      setYAML(initialYAML);
    }
  }, [initialYAML, dirty]);

  function handleChange(value: string) {
    setYAML(value);
    setDirty(value !== initialYAML);
  }

  async function handleSaveManifest() {
    setIsSaving(true);
    setError(null);

    await handle(
      async () => {
        const parsedManifest = parseManifest(yaml);

        await mutateTransaction(
          mutate,
          [
            {
              key: getPluginGetKey(pluginID),
              optimistic: (current) =>
                current ? { ...current, manifest: parsedManifest } : current,
              commit: (_current, result) => result,
            },
          ],
          () => updateManifest(parsedManifest),
          { revalidate: true },
        );

        setDirty(false);
      },
      {
        errorToast: false,
        onError: async (err) => {
          setError(normaliseError(err));
        },
        cleanup: async () => {
          setIsSaving(false);
        },
      },
    );
  }

  return (
    <LStack gap="2" w="full">
      <styled.p fontSize="sm" color="fg.muted">
        Defines plugin metadata and which features the plugin has access to.
      </styled.p>

      <Box
        w="full"
        borderWidth="thin"
        borderColor="border.default"
        borderRadius="md"
        bgColor={editable ? "bg.default" : "bg.subtle"}
        p="3"
      >
        <styled.textarea
          value={yaml}
          onChange={(event) => handleChange(event.currentTarget.value)}
          disabled={!editable}
          minH="96"
          w="full"
          display="block"
          resize="vertical"
          border="none"
          bgColor="transparent"
          p="0"
          fontFamily="mono"
          fontSize="xs"
          lineHeight="tight"
          spellCheck={false}
        />
      </Box>

      {editable && (
        <WStack justifyContent="space-between">
          <styled.p fontSize="xs" color="fg.muted">
            Updating the manifest will force the plugin to disconnect.
          </styled.p>
          <Button
            size="sm"
            variant="subtle"
            onClick={handleSaveManifest}
            disabled={!dirty || isSaving}
            loading={isSaving}
          >
            Save Manifest
          </Button>
        </WStack>
      )}

      <Admonition
        value={!!error}
        kind="failure"
        title="Manifest Update Error"
        onChange={() => setError(null)}
      >
        {error && (
          <LStack>
            {error.overview && (
              <styled.p fontSize="sm">{error.overview}</styled.p>
            )}

            <styled.pre fontSize="xs" mt="1" whiteSpace="pre-wrap">
              {error.details}
            </styled.pre>
          </LStack>
        )}
      </Admonition>
    </LStack>
  );
}

function normaliseError(err: unknown): ManifestError {
  const overview = toStringOrNull(field(err, "message"));
  const detailsFromField = toStringOrNull(field(err, "error"));

  if (detailsFromField) {
    return {
      overview: overview && overview !== detailsFromField ? overview : null,
      details: detailsFromField,
    };
  }

  if (typeof err === "string") {
    return { overview: null, details: err };
  }
  if (err instanceof Error) {
    return { overview: null, details: err.message };
  }
  return { overview: null, details: String(err) };
}

function field(input: unknown, key: string): unknown {
  if (!input || typeof input !== "object") {
    return null;
  }
  return (input as Record<string, unknown>)[key];
}

function toStringOrNull(input: unknown): string | null {
  if (typeof input !== "string") {
    return null;
  }
  const v = input.trim();
  return v === "" ? null : v;
}

function parseManifest(raw: string): Record<string, unknown> {
  const parsed = parse(raw) as unknown;
  if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
    throw new Error("Manifest payload must be an object");
  }
  return parsed as Record<string, unknown>;
}

function toYAML(manifest: Record<string, unknown>): string {
  try {
    return stringify(manifest);
  } catch {
    return "# Failed to render manifest";
  }
}
