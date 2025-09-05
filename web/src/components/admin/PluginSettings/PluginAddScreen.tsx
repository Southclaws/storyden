import {
  FileUploadFileAcceptDetails,
  FileUploadFileRejectDetails,
} from "@ark-ui/react";
import { join } from "lodash";
import { flow, map } from "lodash/fp";
import { FileError } from "node_modules/@ark-ui/react/dist/components/file-upload/file-upload";
import { useState } from "react";
import { mutate } from "swr";

import { handle } from "@/api/client";
import {
  getPluginListKey,
  pluginAdd,
  usePluginAdd,
} from "@/api/openapi-client/plugins";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import * as FileUpload from "@/components/ui/file-upload";
import { AddIcon } from "@/components/ui/icons/Add";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { LStack, WStack, styled } from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";
import { UseDisclosureProps } from "@/utils/useDisclosure";

type Props = UseDisclosureProps;

export function PluginAddScreen({ isOpen, onClose }: Props) {
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleFileAccept = async ({ files }: FileUploadFileAcceptDetails) => {
    const file = files[0];
    if (!file) {
      console.error("handleFileAccept: no file was provided", files);
      return;
    }

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
  };

  const handleFileReject = ({ files }: FileUploadFileRejectDetails) => {
    const errorMessage = mapRejectionToError({ files });
    if (errorMessage) {
      setError(errorMessage);
    }
  };

  const handleClose = () => {
    if (!isUploading) {
      setError(null);
      onClose?.();
    }
  };

  return (
    <ModalDrawer
      isOpen={isOpen}
      onClose={handleClose}
      title="Add Plugin"
      dismissable={!isUploading}
    >
      <LStack gap="4">
        <styled.p color="fg.muted">
          Upload a WebAssembly (.wasm) plugin file to extend Storyden's
          functionality.
        </styled.p>

        <Alert.Root>
          <Alert.Icon asChild>
            <WarningIcon />
          </Alert.Icon>
          <Alert.Content>
            <Alert.Title>Security notice</Alert.Title>
            <Alert.Description>
              Only upload plugins from trusted sources. Malicious plugins can
              compromise the security of your data and system.
            </Alert.Description>
          </Alert.Content>
        </Alert.Root>

        <FileUpload.Root
          maxFiles={1}
          accept=".wasm"
          maxFileSize={10 * 1024 * 1024} // 10MB
          onFileAccept={handleFileAccept}
          onFileReject={handleFileReject}
          disabled={isUploading}
        >
          <FileUpload.Dropzone>
            <LStack gap="3" alignItems="center" p="6">
              <LStack gap="1" textAlign="center">
                <styled.p fontWeight="medium">
                  Drop your .wasm file here or click to browse
                </styled.p>
                <styled.p fontSize="sm" color="fg.muted">
                  WebAssembly files only (.wasm), max 10MB
                </styled.p>
              </LStack>
              <FileUpload.Trigger asChild>
                <Button variant="outline" disabled={isUploading}>
                  <AddIcon />
                  {isUploading ? "Uploading..." : "Select File"}
                </Button>
              </FileUpload.Trigger>
            </LStack>
          </FileUpload.Dropzone>
          <FileUpload.HiddenInput />
        </FileUpload.Root>

        {error && (
          <Alert.Root colorPalette="red">
            <Alert.Icon asChild>
              <WarningIcon />
            </Alert.Icon>
            <Alert.Content>
              <Alert.Title>Upload Error</Alert.Title>
              <Alert.Description>{error}</Alert.Description>
            </Alert.Content>
          </Alert.Root>
        )}

        <WStack justifyContent="end" gap="2">
          <Button
            variant="outline"
            onClick={handleClose}
            disabled={isUploading}
          >
            {isUploading ? "Uploading..." : "Cancel"}
          </Button>
        </WStack>
      </LStack>
    </ModalDrawer>
  );
}

const mapFileError = map((error: FileError) => {
  switch (error) {
    case "FILE_INVALID":
      return "Invalid file.";
    case "FILE_TOO_LARGE":
      return "Plugin file is too large. Maximum size is 10MB.";
    case "FILE_INVALID_TYPE":
      return "File must be a .wasm file";
    case "FILE_TOO_SMALL":
      return "File is too small.";
    case "TOO_MANY_FILES":
      return "Only one plugin file can be uploaded at a time.";
    case "FILE_EXISTS":
      return "A file with this name has already been selected.";
    default:
      return "An unexpected error occurred while reading the file.";
  }
});

const mapFileErrors = flow(mapFileError, join);

const mapRejectionToError = ({ files }: FileUploadFileRejectDetails) => {
  if (files.length === 0) {
    return;
  }

  const file = files[0];
  if (!file) {
    console.error(
      "handleFileReject: files list non-empty but first file is falsy",
    );
    return "An unexpected error occurred while reading the file.";
  }

  const errorMessage = mapFileErrors(file.errors);

  if (errorMessage) {
    return errorMessage;
  }

  return undefined;
};
