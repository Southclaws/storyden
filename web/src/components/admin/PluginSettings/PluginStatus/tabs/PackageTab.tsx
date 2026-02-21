import { useState } from "react";
import { useSWRConfig } from "swr";

import { fetcher, handle } from "@/api/client";
import { mutateTransaction } from "@/api/mutate";
import { getPluginGetKey } from "@/api/openapi-client/plugins";
import { Plugin } from "@/api/openapi-schema";
import { PluginArchiveUpload } from "@/components/admin/PluginSettings/PluginArchiveUpload";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Admonition } from "@/components/ui/admonition";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";

type Props = {
  plugin: Plugin;
};

type PackageTabError = {
  overview: string | null;
  details: string;
};

type PackageTabSuccess = {
  fileName: string;
  previousVersion: string | null;
  newVersion: string | null;
};

export function PackageTab({ plugin }: Props) {
  const { mutate } = useSWRConfig();
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<PackageTabError | null>(null);
  const [success, setSuccess] = useState<PackageTabSuccess | null>(null);

  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(async () => {
      if (!selectedFile || isUploading) {
        return;
      }

      const previousVersion = plugin.version ?? null;

      setIsUploading(true);
      setError(null);
      setSuccess(null);

      await handle(
        async () => {
          const result = await mutateTransaction(
            mutate,
            [
              {
                key: getPluginGetKey(plugin.id),
                optimistic: (current) => current,
                commit: (_current, updated) => updated,
              },
            ],
            () => pluginUpdatePackage(plugin.id, selectedFile),
            { revalidate: true },
          );

          setSuccess({
            fileName: selectedFile.name,
            previousVersion,
            newVersion: result.version ?? null,
          });
          setSelectedFile(null);
        },
        {
          errorToast: false,
          onError: async (err) => {
            setError(normaliseError(err));
          },
          cleanup: async () => {
            setIsUploading(false);
          },
        },
      );
    });

  function handleFileChange(file: File | null) {
    setSelectedFile(file);
    setError(null);
    setSuccess(null);

    if (!file) {
      void handleCancelAction();
    }
  }

  return (
    <LStack gap="4">
      <styled.p fontSize="sm" color="fg.muted">
        Upload a replacement plugin package.{" "}
        {plugin.status.active_state === "active"
          ? "This will restart the plugin with the new version."
          : "The new version will be used when the plugin is enabled."}
      </styled.p>

      <Alert.Root>
        <Alert.Icon asChild>
          <WarningIcon />
        </Alert.Icon>
        <Alert.Content>
          <Alert.Title>Upgrade notice</Alert.Title>
          <Alert.Description>
            The uploaded manifest must match this plugin&apos;s manifest ID.
            Only upload trusted plugin packages.
          </Alert.Description>
        </Alert.Content>
      </Alert.Root>

      <PluginArchiveUpload
        disabled={isUploading}
        buttonLabel={isUploading ? "Uploading..." : "Select File"}
        onFileChange={handleFileChange}
        onError={(details) =>
          setError({
            overview: "Upload rejected",
            details,
          })
        }
      />

      <WStack justifyContent="end" alignItems="end">
        {isConfirming ? (
          <HStack gap="2">
            <Button
              size="sm"
              variant="subtle"
              onClick={handleConfirmAction}
              disabled={!selectedFile || isUploading}
              loading={isUploading}
            >
              Confirm upgrade
            </Button>
            <Button
              size="sm"
              variant="outline"
              onClick={handleCancelAction}
              disabled={isUploading}
            >
              Cancel
            </Button>
          </HStack>
        ) : (
          <Button
            size="sm"
            variant="subtle"
            onClick={handleConfirmAction}
            disabled={!selectedFile || isUploading}
          >
            Upgrade package
          </Button>
        )}
      </WStack>

      <Admonition
        value={!!success}
        kind="success"
        title="Package updated"
        onChange={() => setSuccess(null)}
      >
        {success && (
          <LStack gap="1">
            <styled.p fontSize="sm">
              Uploaded <styled.code>{success.fileName}</styled.code>.
            </styled.p>
            {success.previousVersion || success.newVersion ? (
              <styled.p fontSize="sm">
                Version:{" "}
                <styled.code>
                  {success.previousVersion ?? "-"} → {success.newVersion ?? "-"}
                </styled.code>
              </styled.p>
            ) : (
              <styled.p fontSize="sm">
                The plugin package was replaced successfully.
              </styled.p>
            )}
          </LStack>
        )}
      </Admonition>

      <Admonition
        value={!!error}
        kind="failure"
        title="Package Upgrade Error"
        onChange={() => setError(null)}
      >
        {error && (
          <LStack gap="1">
            {error.overview && (
              <styled.p fontSize="sm">{error.overview}</styled.p>
            )}
            <styled.pre fontSize="xs" whiteSpace="pre-wrap">
              {error.details}
            </styled.pre>
          </LStack>
        )}
      </Admonition>
    </LStack>
  );
}

async function pluginUpdatePackage(
  pluginID: string,
  archive: File,
): Promise<Plugin> {
  return fetcher<Plugin>({
    url: `/plugins/${pluginID}/package`,
    method: "PATCH",
    headers: { "Content-Type": "application/zip" },
    data: archive,
  });
}

function normaliseError(err: unknown): PackageTabError {
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

  const value = input.trim();
  return value === "" ? null : value;
}
