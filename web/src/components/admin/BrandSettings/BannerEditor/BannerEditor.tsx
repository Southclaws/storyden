import {
  FileUploadFileAcceptDetails,
  FileUploadFileRejectDetails,
} from "@ark-ui/react";
import { useRef, useState } from "react";
import {
  FixedCropper,
  FixedCropperRef,
  ImageRestriction,
} from "react-advanced-cropper";
import { toast } from "sonner";

import { handle } from "@/api/client";
import { bannerUpload } from "@/api/openapi-client/misc";
import { Button } from "@/components/ui/button";
import * as FileUpload from "@/components/ui/file-upload";
import { MediaAddIcon } from "@/components/ui/icons/Media";
import { SaveIcon } from "@/components/ui/icons/Save";
import { css } from "@/styled-system/css";
import { Box, HStack, LStack } from "@/styled-system/jsx";
import { getBannerURL } from "@/utils/icon";
import { getExtensionsForMimeTypes } from "@/utils/mime-types";

import "react-advanced-cropper/dist/style.css";

export const CROP_STENCIL_WIDTH = 1200;
export const CROP_STENCIL_HEIGHT = 630;
const ACCEPTED_BANNER_MIMES = ["image/png", "image/jpeg"] as const;

export function BannerEditor() {
  const [bannerURL, setBannerURL] = useState<string | undefined>(
    getBannerURL(),
  );

  const cropperRef = useRef<FixedCropperRef>(null);

  const handleSaveCurrentCrop = async () => {
    if (!cropperRef.current) {
      return;
    }

    const canvas = cropperRef.current.getCanvas();
    if (!canvas) {
      throw new Error("An unexpected error occurred with the image editor.");
    }

    const coordinates = cropperRef.current.getCoordinates();
    if (!coordinates) {
      throw new Error(
        "An unexpected error occurred with the image editor: unable to get crop coordinates.",
      );
    }

    const blob = await new Promise<Blob>((resolve, reject) => {
      canvas.toBlob((blob) => {
        if (blob == null) {
          reject("An unexpected error occurred with the image editor.");
          return;
        }

        resolve(blob);
      });
    });

    await bannerUpload(blob);

    toast.success("Banner saved!");
  };

  async function handleFile({ files }: FileUploadFileAcceptDetails) {
    await handle(async () => {
      // NOTE: For some reason (Zag bug?) this is called for rejected files too.
      const file = files[0];
      if (!file) {
        console.error("handleFile: no file was provided", files);
        return;
      }

      const base64 = await new Promise<string>((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => {
          const result = reader.result;
          if (typeof result !== "string") {
            reject("An unexpected error occurred while reading the file.");
            return;
          }

          resolve(result);
        };
        reader.onerror = () => {
          reject("An unexpected error occurred while reading the file.");
        };

        reader.readAsDataURL(file);
      });

      setBannerURL(base64);
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

    const accepted = getExtensionsForMimeTypes([...ACCEPTED_BANNER_MIMES]);

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
    <LStack gap="1">
      <HStack w="full">
        <FileUpload.Root
          w="min"
          maxFiles={1}
          accept={[...ACCEPTED_BANNER_MIMES]}
          onFileAccept={handleFile}
          onFileReject={handleFileReject}
        >
          <FileUpload.Trigger w="min" asChild>
            <Button type="button" size="xs" variant="outline">
              <MediaAddIcon /> Upload banner
            </Button>
          </FileUpload.Trigger>
          <FileUpload.HiddenInput data-testid="input" />
        </FileUpload.Root>

        <Button
          type="button"
          size="xs"
          variant="solid"
          onClick={handleSaveCurrentCrop}
        >
          <SaveIcon /> Save banner
        </Button>
      </HStack>

      <Box width="full" height="64">
        <FixedCropper
          ref={cropperRef}
          className={css({
            maxWidth: "full",
            maxHeight: "64",
            borderRadius: "md",
            // TODO: Remove black background when empty
            backgroundColor: "bg.default",
          })}
          onTransformImageEnd={handleSaveCurrentCrop}
          // defaultPosition={
          //   initialCoverCoordinates && {
          //     top: initialCoverCoordinates.top,
          //     left: initialCoverCoordinates.left,
          //   }
          // }
          backgroundWrapperProps={{
            scaleImage: false,
          }}
          stencilProps={{
            handlers: false,
            lines: false,
            movable: false,
            resizable: false,
          }}
          stencilSize={{
            width: CROP_STENCIL_WIDTH,
            height: CROP_STENCIL_HEIGHT,
          }}
          imageRestriction={ImageRestriction.stencil}
          src={bannerURL}
        />
      </Box>
    </LStack>
  );
}
