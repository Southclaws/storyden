import { Portal } from "@ark-ui/react";
import type { FileUploadFileAcceptDetails } from "@ark-ui/react";
import { motion, useReducedMotion } from "framer-motion";
import { type ReactNode, useEffect, useId, useRef, useState } from "react";
import { toast } from "sonner";

import { imageCropper } from "@/styled-system/recipes";

import { Button } from "../button";
import * as FileUpload from "../file-upload";
import { IconButton } from "../icon-button";
import { CheckIcon } from "../icons/Check";
import { DeleteIcon } from "../icons/Delete";
import { EditIcon } from "../icons/Edit";
import { MediaAddIcon } from "../icons/Media";
import * as Menu from "../menu";
import { Menu as OptionsMenu } from "../menu/Menu";
import { Spinner } from "../spinner";

const cropShapes = [
  ["free", "Free crop"],
  ["original", "Original proportions"],
  ["1", "Square · 1:1"],
  ["1.3333333333333333", "Landscape · 4:3"],
  ["1.5", "Photograph · 3:2"],
  ["0.6666666666666666", "Portrait · 2:3"],
  ["1.7777777777777777", "Widescreen · 16:9"],
  ["3", "Wide banner · 3:1"],
  ["4", "Banner · 4:1"],
];

type Props = {
  hasImage: boolean;
  busy: boolean;
  saved: boolean;
  onSavedEnd: () => void;
  disabled: boolean;
  toolbar: "collapsible" | "expanded";
  accept: string[];
  menuLabel: string;
  children?: ReactNode;
  onFile: (details: FileUploadFileAcceptDetails) => Promise<void>;
  onRemove?: () => Promise<void>;
  onCropShapeChange?: (shape: string) => void;
};

// Owns the floating toolbar, including its collapsed index-page interaction.
// Menu items extend these controls without placing UI around the crop viewport.
export function ImageCropperTools({
  hasImage,
  busy,
  saved,
  onSavedEnd,
  disabled,
  toolbar,
  accept,
  menuLabel,
  children,
  onFile,
  onRemove,
  onCropShapeChange,
}: Props) {
  const styles = imageCropper();
  const [expanded, setExpanded] = useState(false);
  const closeTimeoutRef = useRef<ReturnType<typeof setTimeout>>(undefined);
  const reducedMotion = useReducedMotion();
  const actionsID = useId();
  const open = toolbar === "expanded" || expanded;
  const hasMenu = Boolean(children || onCropShapeChange);
  useEffect(() => () => clearTimeout(closeTimeoutRef.current), []);

  return (
    <div className={styles.controls} data-has-image={hasImage}>
      <div
        className={styles.toolbar}
        data-expanded={open || busy || saved}
        onPointerEnter={(event) => {
          if (event.pointerType === "touch") return;
          clearTimeout(closeTimeoutRef.current);
          setExpanded(true);
        }}
        onPointerLeave={(event) => {
          if (event.pointerType === "touch") return;
          closeTimeoutRef.current = setTimeout(() => setExpanded(false), 650);
        }}
        onBlur={(event) => {
          if (!event.currentTarget.contains(event.relatedTarget))
            setExpanded(false);
        }}
      >
        {hasImage &&
          (hasMenu ? (
            <OptionsMenu
              trigger={
                <IconButton
                  type="button"
                  variant="ghost"
                  className={styles.editTrigger}
                  aria-label={menuLabel}
                  disabled={disabled}
                >
                  <EditIcon />
                </IconButton>
              }
            >
              {children}
              {onCropShapeChange && (
                <Menu.Root onSelect={({ value }) => onCropShapeChange(value)}>
                  <Menu.TriggerItem>Crop ratio</Menu.TriggerItem>
                  <Portal>
                    <Menu.Positioner>
                      <Menu.Content>
                        {cropShapes.map(([value, label]) => (
                          <Menu.Item key={value} value={value!}>
                            {label}
                          </Menu.Item>
                        ))}
                      </Menu.Content>
                    </Menu.Positioner>
                  </Portal>
                </Menu.Root>
              )}
            </OptionsMenu>
          ) : (
            <IconButton
              type="button"
              variant="ghost"
              className={styles.editTrigger}
              aria-label="Image tools"
              aria-expanded={open}
              aria-controls={actionsID}
              onClick={() => setExpanded((value) => !value)}
            >
              <EditIcon />
            </IconButton>
          ))}
        <motion.div
          className={styles.reveal}
          animate={{
            width: open || !hasImage ? "auto" : busy || saved ? 24 : 0,
          }}
          transition={{ duration: reducedMotion ? 0 : 0.2, ease: "easeInOut" }}
        >
          {(busy || saved) && (
            <div
              className={styles.status}
              data-state={busy ? "saving" : "saved"}
              aria-hidden="true"
              onAnimationEnd={(event) => {
                if (event.animationName === "imageCropperSaved") onSavedEnd();
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
            hidden={!open && hasImage}
          >
            <FileUpload.Root
              className={styles.fileUpload}
              maxFiles={1}
              accept={accept}
              disabled={disabled}
              onFileAccept={onFile}
              onFileReject={() =>
                toast.error(
                  accept.includes("image/gif")
                    ? "Choose a PNG, JPEG, WebP or GIF image."
                    : "Choose a PNG, JPEG or WebP image.",
                )
              }
            >
              <FileUpload.Trigger asChild>
                {hasImage ? (
                  <IconButton
                    type="button"
                    variant="ghost"
                    aria-label="Replace image"
                    disabled={disabled}
                  >
                    <MediaAddIcon />
                  </IconButton>
                ) : (
                  <Button type="button" variant="ghost" disabled={disabled}>
                    <MediaAddIcon />
                    Upload image
                  </Button>
                )}
              </FileUpload.Trigger>
              <FileUpload.HiddenInput />
            </FileUpload.Root>
            {hasImage && onRemove && (
              <IconButton
                type="button"
                variant="ghost"
                aria-label="Remove image"
                disabled={disabled}
                onClick={async () => {
                  await onRemove();
                  setExpanded(false);
                }}
              >
                <DeleteIcon />
              </IconButton>
            )}
          </div>
        </motion.div>
      </div>
    </div>
  );
}
