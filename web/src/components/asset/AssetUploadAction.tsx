import {
  FileUploadFileAcceptDetails,
  FileUploadFileRejectDetails,
} from "@ark-ui/react";
import { ImageIcon, ImagePlusIcon } from "lucide-react";
import mime from "mime-db";
import { PropsWithChildren } from "react";
import { toast } from "sonner";

import { handle } from "@/api/client";
import { assetUpload } from "@/api/openapi-client/assets";
import { Asset, AssetID } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import * as FileUpload from "@/components/ui/file-upload";
import { ButtonVariantProps, button } from "@/styled-system/recipes";

type AssetUploadActionProps = {
  parentAssetID?: AssetID;
  operation: "add" | "update";
  onFinish: (a: Asset) => Promise<void>;
  hideLabel?: boolean;
};

type Props = AssetUploadActionProps & ButtonVariantProps & FileUpload.RootProps;

export function AssetUploadAction({
  children,
  ...props
}: PropsWithChildren<Props>) {
  const [buttonVariantProps, rest] = button.splitVariantProps(props);

  const { onFinish, ...fileUploadProps } = rest;

  const acceptedMIMEs = getMIMEs(props.accept);

  async function handleFile({ files }: FileUploadFileAcceptDetails) {
    await handle(async () => {
      // NOTE: For some reason (Zag bug?) this is called for rejected files too.
      const file = files[0];
      if (!file) {
        console.error("handleFile: no file was provided", files);
        return;
      }

      const asset = await assetUpload(file, {
        filename: file.name,
        parent_asset_id: props.parentAssetID,
      });

      props.onFinish(asset);
    });
  }

  async function handleFileReject({ files }: FileUploadFileRejectDetails) {
    if (files.length === 0) {
      return;
    }

    const file = files[0];
    if (!file) {
      console.error(
        "handleFileReject: files list non-empty but first file is falsy",
      );
      return;
    }

    const accepted = acceptedMIMEs.reduce((prev: string[], curr: string) => {
      const extensions = mime[curr]?.extensions;
      if (!extensions) {
        return prev;
      }

      return [...prev, ...extensions];
    }, []);

    const acceptedList = accepted.map((e) => `.${e}`).join(", ");

    // Vast majority of the time, there will only be one error, but join anyway.
    const errorMessage = file.errors
      .map((error) => {
        switch (error) {
          case "FILE_INVALID":
            return "Invalid file.";
          case "FILE_TOO_LARGE":
            return "File is too large.";
          case "FILE_INVALID_TYPE":
            return `File must be of type ${acceptedList}`;
          case "FILE_TOO_SMALL":
            return "File is too small.";
          case "TOO_MANY_FILES":
            return "Too many files.";
          default:
            return "An unexpected error occurred while reading the file.";
        }
      })
      .join(", ");

    toast.error(errorMessage);
  }

  return (
    <FileUpload.Root
      w="min"
      maxFiles={1}
      onFileAccept={handleFile}
      onFileReject={handleFileReject}
      {...fileUploadProps}
    >
      <FileUpload.Trigger w="min" asChild>
        {children || (
          <Button
            type="button"
            size="xs"
            variant="outline"
            {...buttonVariantProps}
          >
            {props.operation === "add" ? (
              <>
                <ImagePlusIcon />
                {props.hideLabel ? "" : "add cover"}
              </>
            ) : (
              <>
                <ImageIcon /> {props.hideLabel ? "" : "replace cover"}
              </>
            )}
          </Button>
        )}
      </FileUpload.Trigger>
      <FileUpload.HiddenInput data-testid="input" />
    </FileUpload.Root>
  );
}

// NOTE: For some reason, Ark UI's prop type for "accept" also includes a record
// type (not sure what the use-case is) so, we need to convert it into an array.
function getMIMEs(
  accept: Record<string, string[]> | string | string[] | undefined,
): string[] {
  if (!accept) {
    return [];
  }

  if (typeof accept === "string") {
    return [accept];
  }

  if (Array.isArray(accept)) {
    return accept;
  }

  const mimes = Object.keys(accept);

  return mimes;
}
