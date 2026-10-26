"use client";

import type { FileUploadFileAcceptDetails } from "@ark-ui/react";
import {
  type CSSProperties,
  type ComponentPropsWithoutRef,
  type KeyboardEvent,
  type ReactNode,
  type Ref,
  useCallback,
  useEffect,
  useEffectEvent,
  useId,
  useImperativeHandle,
  useRef,
  useState,
} from "react";
import {
  Cropper,
  type CropperRef,
  FixedCropper,
  type FixedCropperRef,
  ImageRestriction,
} from "react-advanced-cropper";
import { toast } from "sonner";

import { handle } from "@/api/client";
import { cx } from "@/styled-system/css";
import { imageCropper } from "@/styled-system/recipes";

import { Text } from "../text";

import "react-advanced-cropper/dist/style.css";

import { ImageCropperTools } from "./ImageCropperTools";

export type ImageCropperCoordinates = {
  left: number;
  top: number;
  width: number;
  height: number;
};
export type ImageCropperSize = { width: number; height: number };
/** A captured adjustment, exportable even after the editor unmounts. */
export type ImageCropperSnapshot = {
  coordinates: ImageCropperCoordinates;
  imageSize: ImageCropperSize;
  width: number;
  height: number;
  getBlob: () => Promise<Blob>;
};
export type ImageCropperRef = {
  getCrop: () => ImageCropperSnapshot | null;
  getCroppedImage: () => Promise<Blob | null>;
  getCoordinates: () => ImageCropperCoordinates | null;
};
export type ImageCropperProps = Omit<
  ComponentPropsWithoutRef<"div">,
  "onChange" | "children"
> & {
  ref?: Ref<ImageCropperRef>;
  "aria-label"?: string;
  src?: string;
  width: number;
  height: number;
  disabled?: boolean;
  saving?: boolean;
  saved?: boolean;
  defaultCoordinates?:
    | ImageCropperCoordinates
    | ((imageSize: ImageCropperSize) => ImageCropperCoordinates);
  defaultPosition?: { left: number; top: number };
  /** Fixed fills the frame; resizable exposes a selection and crop-ratio menu. */
  cropMode?: "fixed" | "resizable";
  /** Absent: export exactly width × height. Present: retain crop proportions. */
  maxOutputSize?: ImageCropperSize;
  /** False displays the caller's uploaded src without immediately cropping it. */
  cropOnUpload?: boolean;
  toolbar?: "collapsible" | "expanded";
  accept?: string[];
  showHint?: boolean;
  menuLabel?: string;
  children?: ReactNode;
  className?: string;
  style?: CSSProperties;
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
  defaultPosition,
  cropMode = "fixed",
  maxOutputSize,
  cropOnUpload = true,
  toolbar = "collapsible",
  accept = ["image/png", "image/jpeg", "image/webp"],
  showHint = true,
  menuLabel = "Image options",
  children,
  className,
  style,
  onFileChange,
  onRemove,
  onChange,
  ...rootProps
}: ImageCropperProps) {
  const cropperRef = useRef<CropperRef | FixedCropperRef>(null);
  const attachCropper = useCallback(
    (instance: CropperRef | FixedCropperRef | null) => {
      cropperRef.current = instance;
    },
    [],
  );
  const [selectedURL, setSelectedURL] = useState<string>();
  const [ready, setReady] = useState(false);
  const [loading, setLoading] = useState(Boolean(src));
  const [showSaved, setShowSaved] = useState(saved);
  const [loadError, setLoadError] = useState(false);
  const [originalRatio, setOriginalRatio] = useState(width / height);
  const [selection, setSelection] = useState({ shape: "free", revision: 0 });
  const imageURL = selectedURL ?? src;
  const busy = loading || saving;
  const controlsDisabled = disabled || saving;
  const cropDisabled = controlsDisabled || !ready;
  const descriptionID = useId();
  const styles = imageCropper();
  const notifyChange = useEffectEvent(() => onChange?.());
  let statusMessage = "";
  if (loading) statusMessage = "Loading image";
  else if (saving) statusMessage = "Saving image";
  else if (loadError)
    statusMessage = "The image could not be opened. Please choose another.";
  else if (showSaved) statusMessage = "Image saved";

  useEffect(() => setShowSaved(saved && !busy), [saved, busy]);

  // Ready precedes the canvas commit. Menu choices and default local uploads
  // notify after that commit, so consumers always capture the finished crop.
  useEffect(() => {
    if (ready && (selectedURL || selection.revision > 0)) notifyChange();
  }, [ready, selectedURL, selection.revision]);
  useEffect(
    () => () => {
      if (selectedURL) URL.revokeObjectURL(selectedURL);
    },
    [selectedURL],
  );

  function captureCrop(): ImageCropperSnapshot | null {
    const cropper = cropperRef.current;
    // Preserve fractional source coordinates when zooming to avoid rounding
    // the fixed ratio or shifting the crop when it is reopened.
    const coordinates = cropper?.getCoordinates({ round: false });
    const image = cropper?.getImage();
    const canvas = cropper?.getCanvas(
      maxOutputSize
        ? {
            minWidth: 1,
            minHeight: 1,
            maxWidth: maxOutputSize.width,
            maxHeight: maxOutputSize.height,
          }
        : { width, height },
    );
    if (!canvas || !coordinates || !image) return null;
    return {
      coordinates,
      imageSize: { width: image.width, height: image.height },
      width: canvas.width,
      height: canvas.height,
      getBlob: () =>
        new Promise<Blob>((resolve, reject) => {
          canvas.toBlob((blob) => {
            if (blob) resolve(blob);
            else reject(new Error("The cropped image could not be saved."));
          }, "image/png");
        }),
    };
  }

  useImperativeHandle(ref, () => ({
    getCrop: captureCrop,
    getCoordinates: () => cropperRef.current?.getCoordinates() ?? null,
    getCroppedImage: async () => captureCrop()?.getBlob() ?? null,
  }));

  async function handleFile({ files }: FileUploadFileAcceptDetails) {
    const file = files[0];
    if (!file) {
      return;
    }

    await handle(async () => {
      await onFileChange?.(file);

      // Original uploads are rendered through the caller's src; only the
      // default upload-and-crop flow needs a temporary local preview.
      if (cropOnUpload) {
        setReady(false);
        setSelectedURL(URL.createObjectURL(file));
      }
    });
  }

  async function handleRemove() {
    await handle(async () => {
      await onRemove?.();

      setReady(false);
      setLoading(false);
      setLoadError(false);
      setSelectedURL(undefined);
    });
  }

  const cropperProps = {
    ref: attachCropper,
    className: styles.cropper,
    disabled: cropDisabled,
    src: imageURL,
    defaultCoordinates: getInitialCoordinates(
      selection,
      selectedURL ? undefined : defaultCoordinates,
      cropMode === "fixed" ? width / height : undefined,
    ),
    defaultPosition: selection.revision === 0 ? defaultPosition : undefined,
    onUpdate: (cropper: CropperRef | FixedCropperRef) => {
      if (cropper.isLoading()) {
        setReady(false);
        setLoading(true);
        setLoadError(false);
      }
    },
    onReady: (cropper: CropperRef | FixedCropperRef) => {
      const image = cropper.getImage();
      if (image) setOriginalRatio(image.width / image.height);
      setReady(true);
      setLoading(false);
    },
    onError: () => {
      setReady(false);
      setLoading(false);
      setLoadError(true);
      toast.error("The image could not be opened. Please choose another.");
    },
    onInteractionEnd: cropDisabled ? undefined : onChange,
    backgroundWrapperProps: {
      scaleImage:
        cropMode === "fixed" ? { touch: true, wheel: { ratio: 0.025 } } : false,
    },
    imageRestriction: ImageRestriction.stencil,
  };

  return (
    <div
      {...rootProps}
      className={cx(styles.root, className)}
      role="group"
      aria-label={ariaLabel}
      style={
        {
          "--image-cropper-aspect-ratio": width / height,
          ...style,
        } as CSSProperties
      }
    >
      {/* The canvas is pointer-editable; this focusable group adds keyboard panning. */}
      {/* eslint-disable-next-line jsx-a11y/no-noninteractive-element-interactions */}
      <div
        className={styles.viewport}
        role="group"
        aria-label="Crop area"
        aria-describedby={descriptionID}
        aria-busy={busy}
        aria-disabled={cropDisabled}
        tabIndex={cropDisabled ? undefined : 0}
        data-saving={saving}
        onKeyDown={(event) => {
          if (!cropDisabled)
            handleCropperKeyDown(event, cropperRef.current, cropMode);
        }}
      >
        {cropMode === "fixed" ? (
          <FixedCropper
            key={`fixed:${selection.revision}`}
            {...cropperProps}
            stencilProps={{
              className: "image-cropper__selection",
              handlers: false,
              lines: false,
              movable: false,
              resizable: false,
            }}
            stencilSize={{ width, height }}
          />
        ) : (
          <Cropper
            key={`resizable:${selection.revision}`}
            {...cropperProps}
            stencilProps={{
              className: "image-cropper__selection",
              aspectRatio:
                selection.shape === "original"
                  ? originalRatio
                  : Number(selection.shape) || undefined,
            }}
          />
        )}
      </div>

      <Text as="span" id={descriptionID} className={styles.description}>
        Upload an image, then drag or use the arrow keys to reposition it. Hold
        Shift for larger adjustments. Cropping is unavailable while loading or
        saving.{" "}
        {cropMode === "fixed" &&
          "Pinch or scroll to zoom, or use the plus and minus keys."}{" "}
        {cropMode === "resizable" &&
          `Drag the crop handles to resize the selection. Choose a crop ratio from ${menuLabel}.`}
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

      <ImageCropperTools
        hasImage={Boolean(imageURL)}
        busy={busy}
        saved={showSaved}
        onSavedEnd={() => setShowSaved(false)}
        disabled={controlsDisabled}
        toolbar={toolbar}
        accept={accept}
        menuLabel={menuLabel}
        onFile={handleFile}
        onRemove={onRemove ? handleRemove : undefined}
        onCropShapeChange={
          cropMode === "resizable"
            ? (shape) => {
                setReady(false);
                setSelection((value) => ({
                  shape,
                  revision: value.revision + 1,
                }));
              }
            : undefined
        }
      >
        {children}
      </ImageCropperTools>

      {showHint && ready && !controlsDisabled && (
        <Text className={styles.hint} variant="metadata">
          {cropMode === "fixed"
            ? "Drag to move · Pinch or scroll to zoom"
            : "Drag to reposition"}
        </Text>
      )}
    </div>
  );
}

function getInitialCoordinates(
  selection: { shape: string; revision: number },
  defaults: ImageCropperProps["defaultCoordinates"],
  fixedRatio?: number,
) {
  if (selection.shape === "original")
    return ({ imageSize }: { imageSize: ImageCropperSize }) => ({
      left: 0,
      top: 0,
      ...imageSize,
    });
  if (selection.revision > 0 || !defaults) return undefined;
  return ({ imageSize }: { imageSize: ImageCropperSize }) => {
    const coordinates =
      typeof defaults === "function" ? defaults(imageSize) : defaults;
    if (!fixedRatio) return coordinates;
    // Match the centered cover presentation without expanding a saved free crop.
    const width = Math.min(coordinates.width, coordinates.height * fixedRatio);
    const height = width / fixedRatio;
    return {
      left: coordinates.left + (coordinates.width - width) / 2,
      top: coordinates.top + (coordinates.height - height) / 2,
      width,
      height,
    };
  };
}

function handleCropperKeyDown(
  event: KeyboardEvent<HTMLDivElement>,
  cropper: CropperRef | FixedCropperRef | null,
  cropMode: NonNullable<ImageCropperProps["cropMode"]>,
) {
  if (!cropper) {
    return;
  }

  if (
    cropMode === "fixed" &&
    !event.ctrlKey &&
    !event.metaKey &&
    !event.altKey &&
    ["+", "=", "-", "_"].includes(event.key)
  ) {
    event.preventDefault();
    cropper.zoomImage(event.key === "+" || event.key === "=" ? 1.1 : 1 / 1.1, {
      interaction: true,
    });
    cropper.transformImageEnd();
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
