import { usePinch, useWheel } from "@use-gesture/react";
import { ChangeEvent, useEffect, useRef, useState } from "react";
import AvatarEditor from "react-avatar-editor";

import { handle } from "@/api/client";
import { assetUpload } from "@/api/openapi-client/assets";
import { Asset } from "@/api/openapi-schema";
import { getAssetURL } from "@/utils/asset";

export type Props = {
  value?: Asset;
  onUpload: (asset: Asset) => void;
};

export function useAssetUploadEditor(props: Props) {
  const ref = useRef<AvatarEditor>(null);

  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [tempFile, setTempFile] = useState<File | null>(null);
  const [saving, setSaving] = useState(false);
  const [scale, setScale] = useState(1);

  const pinchCaptureRef = useRef<HTMLDivElement>(null);

  const file =
    tempFile || (props.value ? getAssetURL(props.value.path) : "") || "";

  usePinch(
    ({ offset: [s], event }) => {
      setScale(Math.min(Math.max(s, 1), 100));
    },
    {
      target: pinchCaptureRef,
      scaleBounds() {
        return { min: 1, max: 100 };
      },
      preventDefault: true,
    },
  );

  useWheel(
    ({ delta: [_, dy], event }) => {
      setScale((prevScale) => {
        const newScale = prevScale - dy / 100;
        return Math.min(Math.max(newScale, 1), 100);
      });
    },
    {
      target: pinchCaptureRef,
      preventDefault: true,
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
    setScale(1);
    setPosition({ x: 0, y: 0 });
    setTempFile(null);
  }, [props.value]);

  function onPositionChange(p: any) {
    if (saving) return;

    setPosition(p);
  }

  function onFileChange(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) {
      return;
    }

    setTempFile(file);
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

      await handle(
        async () => {
          const asset = await assetUpload(f, {
            filename: "cropped-image",
          });

          props.onUpload(asset);
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
  };
}
