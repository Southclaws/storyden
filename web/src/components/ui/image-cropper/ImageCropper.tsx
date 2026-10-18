"use client";

import type { FileUploadFileAcceptDetails } from "@ark-ui/react";
import { motion, useReducedMotion } from "framer-motion";
import {
  type CSSProperties,
  type KeyboardEvent,
  type Ref,
  useEffect,
  useEffectEvent,
  useId,
  useImperativeHandle,
  useRef,
  useState,
} from "react";
import {
  FixedCropper,
  type FixedCropperRef,
  ImageRestriction,
} from "react-advanced-cropper";
import { toast } from "sonner";

import { handle } from "@/api/client";
import { imageCropper } from "@/styled-system/recipes";

import { Button } from "../button";
import * as FileUpload from "../file-upload";
import { IconButton } from "../icon-button";
import { CheckIcon } from "../icons/Check";
import { DeleteIcon } from "../icons/Delete";
import { EditIcon } from "../icons/Edit";
import { MediaAddIcon } from "../icons/Media";
import { Spinner } from "../spinner";
import { Text } from "../text";

import "react-advanced-cropper/dist/style.css";

export type ImageCropperCoordinates = NonNullable<
  ReturnType<FixedCropperRef["getCoordinates"]>
>;

export type ImageCropperRef = {
  getCroppedImage: () => Promise<Blob | null>;
  getCoordinates: () => ImageCropperCoordinates | null;
};

export type ImageCropperProps = {
  ref?: Ref<ImageCropperRef>;
  "aria-label"?: string;
  src?: string;
  width: number;
  height: number;
  disabled?: boolean;
  saving?: boolean;
  saved?: boolean;
  defaultCoordinates?: ImageCropperCoordinates;
  onFileChange?: (file: File) => Promise<void>;
  onRemove?: () => Promise<void>;
  onChange?: () => void;
};

export function ImageCropper({
  ref,
  "aria-label": ariaLabel = "Image cropper",
  src,
  width,
  height,
  disabled = false,
  saving = false,
  saved = false,
  defaultCoordinates,
  onFileChange,
  onRemove,
  onChange,
}: ImageCropperProps) {
  const cropperRef = useRef<FixedCropperRef>(null);
  const [selectedURL, setSelectedURL] = useState<string>();
  const [ready, setReady] = useState(false);
  const [loading, setLoading] = useState(Boolean(src));
  const [loadError, setLoadError] = useState(false);
  const [expanded, setExpanded] = useState(false);
  const [showSaved, setShowSaved] = useState(saved);
  const closeTimeoutRef = useRef<ReturnType<typeof setTimeout>>(undefined);
  const reducedMotion = useReducedMotion();
  const imageURL = selectedURL ?? src;
  const busy = loading || saving;
  const controlsDisabled = disabled || saving;
  const cropDisabled = controlsDisabled || !ready;
  const descriptionID = useId();
  const actionsID = useId();
  const styles = imageCropper();
  const notifyChange = useEffectEvent(() => onChange?.());
  let statusMessage = "";

  if (loading) {
    statusMessage = "Loading image";
  } else if (saving) {
    statusMessage = "Saving image";
  } else if (loadError) {
    statusMessage = "The image could not be opened. Please choose another.";
  } else if (showSaved) {
    statusMessage = "Image saved";
  }

  useEffect(() => {
    setShowSaved(saved && !busy);
  }, [saved, busy]);

  useEffect(() => {
    return () => clearTimeout(closeTimeoutRef.current);
  }, []);

  // The cropper's ready callback runs before its canvas is committed.
  useEffect(() => {
    if (ready && selectedURL) {
      notifyChange();
    }
  }, [ready, selectedURL]);

  useEffect(() => {
    return () => {
      if (selectedURL) {
        URL.revokeObjectURL(selectedURL);
      }
    };
  }, [selectedURL]);

  useImperativeHandle(ref, () => ({
    getCoordinates() {
      return cropperRef.current?.getCoordinates() ?? null;
    },
    async getCroppedImage() {
      const canvas = cropperRef.current?.getCanvas({ width, height });
      if (!canvas) {
        return null;
      }

      return new Promise<Blob>((resolve, reject) => {
        canvas.toBlob((blob) => {
          if (!blob) {
            reject(new Error("The cropped image could not be saved."));
            return;
          }

          resolve(blob);
        }, "image/png");
      });
    },
  }));

  async function handleFile({ files }: FileUploadFileAcceptDetails) {
    const file = files[0];
    if (!file) {
      return;
    }

    await handle(async () => {
      await onFileChange?.(file);

      setReady(false);
      setSelectedURL(URL.createObjectURL(file));
    });
  }

  async function handleRemove() {
    await handle(async () => {
      await onRemove?.();

      setReady(false);
      setLoading(false);
      setLoadError(false);
      setSelectedURL(undefined);
      setExpanded(false);
    });
  }

  function handleKeyDown(event: KeyboardEvent<HTMLDivElement>) {
    const cropper = cropperRef.current;
    if (cropDisabled || !cropper) {
      return;
    }

    const step = event.shiftKey ? 50 : 10;
    let left = 0;
    let top = 0;

    switch (event.key) {
      case "ArrowLeft":
        left = -step;
        break;
      case "ArrowRight":
        left = step;
        break;
      case "ArrowUp":
        top = -step;
        break;
      case "ArrowDown":
        top = step;
        break;
      default:
        return;
    }

    event.preventDefault();
    // The cropper moves the viewport, so invert directions to move the image.
    cropper.moveImage(-left, -top, { interaction: true });
    cropper.transformImageEnd();
  }

  return (
    <div
      className={styles.root}
      role="group"
      aria-label={ariaLabel}
      style={
        { "--image-cropper-aspect-ratio": width / height } as CSSProperties
      }
    >
      <div
        className={styles.viewport}
        role="group"
        aria-label="Crop area"
        aria-describedby={descriptionID}
        aria-busy={busy}
        aria-disabled={cropDisabled}
        tabIndex={cropDisabled ? undefined : 0}
        data-saving={saving}
        onKeyDown={handleKeyDown}
      >
        <FixedCropper
          ref={cropperRef}
          className={styles.cropper}
          disabled={cropDisabled}
          src={imageURL}
          defaultCoordinates={selectedURL ? undefined : defaultCoordinates}
          onUpdate={(cropper) => {
            if (cropper.isLoading()) {
              setReady(false);
              setLoading(true);
              setLoadError(false);
            }
          }}
          onReady={() => {
            setReady(true);
            setLoading(false);
          }}
          onError={() => {
            setReady(false);
            setLoading(false);
            setLoadError(true);
            toast.error(
              "The image could not be opened. Please choose another.",
            );
          }}
          onInteractionEnd={cropDisabled ? undefined : onChange}
          backgroundWrapperProps={{ scaleImage: false }}
          stencilProps={{
            handlers: false,
            lines: false,
            movable: false,
            resizable: false,
          }}
          stencilSize={{ width, height }}
          imageRestriction={ImageRestriction.stencil}
        />
      </div>

      <Text as="span" id={descriptionID} className={styles.description}>
        Upload a PNG, JPEG or WebP image, then drag or use the arrow keys to
        reposition it. Hold Shift for larger adjustments. Cropping is
        unavailable while loading or saving.
      </Text>
      <Text
        as="span"
        className={styles.description}
        role="status"
        aria-label={statusMessage || undefined}
        aria-live="polite"
        aria-atomic="true"
      >
        {statusMessage}
      </Text>

      <div className={styles.controls} data-has-image={Boolean(imageURL)}>
        <div
          className={styles.toolbar}
          data-expanded={expanded || busy || showSaved}
          onPointerEnter={(event) => {
            if (event.pointerType === "touch") {
              return;
            }

            clearTimeout(closeTimeoutRef.current);
            setExpanded(true);
          }}
          onPointerLeave={(event) => {
            if (event.pointerType === "touch") {
              return;
            }

            closeTimeoutRef.current = setTimeout(() => setExpanded(false), 650);
          }}
          onBlur={(event) => {
            if (!event.currentTarget.contains(event.relatedTarget)) {
              setExpanded(false);
            }
          }}
        >
          {imageURL && (
            <IconButton
              type="button"
              variant="ghost"
              className={styles.editTrigger}
              aria-label="Image tools"
              aria-expanded={expanded}
              aria-controls={actionsID}
              onClick={() => setExpanded((value) => !value)}
            >
              <EditIcon />
            </IconButton>
          )}

          <motion.div
            className={styles.reveal}
            animate={{
              width:
                expanded || !imageURL ? "auto" : busy || showSaved ? 24 : 0,
            }}
            transition={{
              duration: reducedMotion ? 0 : 0.2,
              ease: "easeInOut",
            }}
          >
            {(busy || showSaved) && (
              <div
                className={styles.status}
                data-state={busy ? "saving" : "saved"}
                aria-hidden="true"
                onAnimationEnd={(event) => {
                  if (event.animationName === "imageCropperSaved") {
                    setShowSaved(false);
                  }
                }}
              >
                {busy ? (
                  <Spinner size="sm" aria-hidden="true" />
                ) : (
                  <CheckIcon className={styles.checkmark} aria-hidden="true" />
                )}
              </div>
            )}
            <div
              id={actionsID}
              className={styles.actions}
              role="group"
              aria-label="Image actions"
              hidden={!expanded && Boolean(imageURL)}
            >
              <FileUpload.Root
                className={styles.fileUpload}
                maxFiles={1}
                accept={["image/png", "image/jpeg", "image/webp"]}
                disabled={controlsDisabled}
                onFileAccept={handleFile}
                onFileReject={() =>
                  toast.error("Choose a PNG, JPEG or WebP image.")
                }
              >
                <FileUpload.Trigger asChild>
                  {imageURL ? (
                    <IconButton
                      type="button"
                      variant="ghost"
                      aria-label="Replace image"
                      disabled={controlsDisabled}
                    >
                      <MediaAddIcon />
                    </IconButton>
                  ) : (
                    <Button
                      type="button"
                      variant="ghost"
                      disabled={controlsDisabled}
                    >
                      <MediaAddIcon />
                      Upload image
                    </Button>
                  )}
                </FileUpload.Trigger>
                <FileUpload.HiddenInput />
              </FileUpload.Root>
              {imageURL && onRemove && (
                <IconButton
                  type="button"
                  variant="ghost"
                  aria-label="Remove image"
                  disabled={controlsDisabled}
                  onClick={handleRemove}
                >
                  <DeleteIcon />
                </IconButton>
              )}
            </div>
          </motion.div>
        </div>
      </div>

      {ready && !controlsDisabled && (
        <Text className={styles.hint} variant="metadata">
          Drag to reposition
        </Text>
      )}
    </div>
  );
}
