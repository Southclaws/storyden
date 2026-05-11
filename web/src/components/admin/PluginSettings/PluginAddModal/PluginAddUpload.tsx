import { useState } from "react";
import { mutate } from "swr";

import { handle } from "@/api/client";
import { getPluginListKey, pluginAdd } from "@/api/openapi-client/plugins";
import { PluginArchiveUpload } from "@/components/admin/PluginSettings/PluginArchiveUpload";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { useI18n } from "@/i18n/provider";
import { LStack, WStack, styled } from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";
import { UseDisclosureProps } from "@/utils/useDisclosure";

export function PluginAddUpload({ onClose }: UseDisclosureProps) {
  const { t } = useI18n();
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function uploadArchive(file: File) {
    setIsUploading(true);
    setError(null);

    await handle(
      async () => {
        await pluginAdd(file);
        await mutate(getPluginListKey());
        onClose?.();
      },
      {
        async onError(err) {
          const error = deriveError(err);
          setError(error);
        },
        async cleanup() {
          setIsUploading(false);
        },
      },
    );
  }

  function handleFileChange(file: File | null) {
    if (!file) {
      return;
    }

    void uploadArchive(file);
  }

  const handleClose = () => {
    if (!isUploading) {
      setError(null);
      onClose?.();
    }
  };

  return (
    <LStack gap="4">
      <styled.p color="fg.muted">
        {t(
          "Upload a Storyden Plugin (.sdx or .zip) file to extend Storyden's functionality.",
        )}
      </styled.p>

      <Alert.Root>
        <Alert.Icon asChild>
          <WarningIcon />
        </Alert.Icon>
        <Alert.Content>
          <Alert.Title>{t("Security notice")}</Alert.Title>
          <Alert.Description>
            {t(
              "Only upload plugins from trusted sources. Malicious plugins can compromise the security of your data and system.",
            )}
          </Alert.Description>
        </Alert.Content>
      </Alert.Root>

      <PluginArchiveUpload
        disabled={isUploading}
        buttonLabel={isUploading ? t("Uploading...") : t("Select File")}
        onFileChange={handleFileChange}
        onError={setError}
      />

      {error && (
        <Alert.Root colorPalette="red">
          <Alert.Icon asChild>
            <WarningIcon />
          </Alert.Icon>
          <Alert.Content>
            <Alert.Title>{t("Upload Error")}</Alert.Title>
            <Alert.Description>{error}</Alert.Description>
          </Alert.Content>
        </Alert.Root>
      )}

      <WStack justifyContent="end" gap="2">
        <Button variant="outline" onClick={handleClose} disabled={isUploading}>
          {isUploading ? t("Uploading...") : t("Cancel")}
        </Button>
      </WStack>
    </LStack>
  );
}
