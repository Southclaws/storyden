import { usePinch, useWheel } from "@use-gesture/react";
import { ChangeEvent, useEffect, useRef, useState } from "react";
import AvatarEditor from "react-avatar-editor";

import { handle } from "@/api/client";
import { accountSetAvatar } from "@/api/openapi-client/accounts";
import { assetUpload } from "@/api/openapi-client/assets";
import { Asset } from "@/api/openapi-schema";
import { useSession } from "@/auth";

export type Props = {
  initialValue?: File | undefined;
  isAvatar?: boolean;
  showPreviews?: boolean;

  onUpload: (asset: Asset | undefined) => void;
};

export function useIconEditor(props: Props) {
  const session = useSession();
  const ref = useRef<AvatarEditor>(null);

  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [file, setFile] = useState<File | string>(props.initialValue ?? "");
  const [saving, setSaving] = useState(false);
  const [isDirty, setIsDirty] = useState(false);
  const [scale, setScale] = useState(1);

  const pinchCaptureRef = useRef<HTMLDivElement>(null);

  usePinch(
    ({ offset: [s] }) => {
      setScale(s);
    },
    {
      target: pinchCaptureRef,
      scaleBounds() {
        return { min: 1, max: 100 };
      },
    },
  );

  useWheel(
    ({ offset: [_, dy] }) => {
      setScale((prevScale) => {
        const newScale = prevScale - dy / 1000;
        return Math.min(Math.max(newScale, 1), 100);
      });
    },
    {
      target: pinchCaptureRef,
    },
  );

  useEffect(() => {
    const handler = (e: Event) => e.preventDefault();
    document.addEventListener("gesturestart", handler);
    document.addEventListener("gesturechange", handler);
    document.addEventListener("gestureend", handler);
    return () => {
      document.removeEventListener("gesturestart", handler);
      document.removeEventListener("gesturechange", handler);
      document.removeEventListener("gestureend", handler);
    };
  }, []);

  useEffect(() => {
    if (props.initialValue && !isDirty) {
      setFile(props.initialValue);
    }
  }, [props.initialValue, isDirty]);

  function onPositionChange(p: any) {
    // Don't allow position changes while upload in progress.
    if (saving) return;

    setPosition(p);
    setIsDirty(true);
  }

  function onFileChange(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];

    if (!file) {
      throw new Error("Unexpected problem: File is missing from uploader.");
    }

    setFile(file);
    setIsDirty(true);
  }

  function onSave() {
    if (!ref || !ref.current) {
      return;
    }

    setSaving(true);

    const canvasScaled = ref.current.getImageScaledToCanvas();
    canvasScaled.toBlob(async (f: Blob | null) => {
      if (!f) {
        throw new Error("Image scaled to canvas binary data is null.");
      }

      handle(
        async () => {
          if (props.isAvatar && session) {
            await accountSetAvatar(f);
            props.onUpload(undefined);
          } else {
            const asset = await assetUpload(f);
            props.onUpload(asset);
          }

          setIsDirty(false);
        },
        {
          promiseToast: {
            loading: "Uploading image...",
            success: "Upload complete!",
          },
          cleanup: async () => {
            setSaving(false);
          },
        },
      );
    });
  }

  return {
    ref,
    pinchCaptureRef,
    scale,
    position,
    onPositionChange,
    onFileChange,
    onSave,
    saving,
    file,
    isDirty,
  };
}
