import {
  FileUploadFileChangeDetails,
} from "@ark-ui/react";
import { FileError } from "node_modules/@ark-ui/react/dist/components/file-upload/file-upload";

import { Button } from "@/components/ui/button";
import * as FileUpload from "@/components/ui/file-upload";
import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { DeleteSmallIcon } from "@/components/ui/icons/Delete";
import { VStack, styled } from "@/styled-system/jsx";

type Props = {
  disabled?: boolean;
  buttonLabel?: string;
  onFileChange: (file: File | null) => void;
  onError?: (error: string) => void;
};

export function PluginArchiveUpload({
  disabled,
  buttonLabel = "Select File",
  onFileChange,
  onError,
}: Props) {
  function handleFileChange(details: FileUploadFileChangeDetails) {
    const file = details.acceptedFiles[0] ?? null;
    onFileChange(file);

    const rejectionError = mapPluginArchiveRejectionToError(details);
    if (rejectionError) {
      onError?.(rejectionError);
    }
  }

  return (
    <FileUpload.Root
      maxFiles={1}
      accept=".sdx,.zip"
      maxFileSize={50 * 1024 * 1024}
      onFileChange={handleFileChange}
      disabled={disabled}
    >
      <FileUpload.Dropzone minHeight="44">
        <VStack>
          <VStack gap="1" textAlign="center">
            <styled.p fontWeight="medium">
              Drop your plugin file here or click to browse
            </styled.p>
            <styled.p fontSize="sm" color="fg.muted">
              Plugin files only (.zip or .sdx), max 50MB
            </styled.p>
          </VStack>
          <FileUpload.Trigger asChild>
            <Button variant="outline" disabled={disabled}>
              <AddIcon />
              {buttonLabel}
            </Button>
          </FileUpload.Trigger>
        </VStack>
      </FileUpload.Dropzone>

      <FileUpload.ItemGroup>
        <FileUpload.Context>
          {({ acceptedFiles }) =>
            acceptedFiles.map((file) => (
              <FileUpload.Item key={file.name} file={file} alignItems="center">
                <FileUpload.ItemName />
                <FileUpload.ItemSizeText />
                <FileUpload.ItemDeleteTrigger asChild>
                  <IconButton size="xs" variant="ghost">
                    <DeleteSmallIcon h="4" />
                  </IconButton>
                </FileUpload.ItemDeleteTrigger>
              </FileUpload.Item>
            ))
          }
        </FileUpload.Context>
      </FileUpload.ItemGroup>

      <FileUpload.HiddenInput />
    </FileUpload.Root>
  );
}

function mapPluginArchiveRejectionToError(
  details: FileUploadFileChangeDetails,
): string | undefined {
  const file = details.rejectedFiles[0];
  if (!file) {
    return undefined;
  }

  const messages = file.errors.map(mapFileError).filter(Boolean);
  if (messages.length === 0) {
    return undefined;
  }

  return messages.join(", ");
}

function mapFileError(error: FileError): string {
  switch (error) {
    case "FILE_INVALID":
      return "Invalid file.";
    case "FILE_TOO_LARGE":
      return "Plugin file is too large. Maximum size is 50MB.";
    case "FILE_INVALID_TYPE":
      return "File must be a .zip or .sdx archive.";
    case "FILE_TOO_SMALL":
      return "File is too small.";
    case "TOO_MANY_FILES":
      return "Only one plugin file can be uploaded at a time.";
    case "FILE_EXISTS":
      return "A file with this name has already been selected.";
    default:
      return "An unexpected error occurred while reading the file.";
  }
}
