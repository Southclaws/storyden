"use client";

import { useMemo, useState } from "react";

import { handle } from "@/api/client";
import { Asset } from "@/api/openapi-schema";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { AdminSettings } from "@/lib/settings/settings";
import { uploadThemeAsset } from "@/lib/theme/theme-client";

export type Props = {
  settings: AdminSettings;
};

type AssetType = "css" | "script";

export function useThemeSettings({ settings }: Props) {
  const { revalidate, updateSettings } = useSettingsMutation();
  const [css, setCSS] = useState(settings.metadata.theme.css ?? []);
  const [scripts, setScripts] = useState(settings.metadata.theme.scripts ?? []);
  const [isUploading, setIsUploading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  const hasChanges = useMemo(() => {
    const initialCSS = settings.metadata.theme.css ?? [];
    const initialScripts = settings.metadata.theme.scripts ?? [];

    return (
      JSON.stringify(initialCSS) !== JSON.stringify(css) ||
      JSON.stringify(initialScripts) !== JSON.stringify(scripts)
    );
  }, [css, scripts, settings.metadata.theme.css, settings.metadata.theme.scripts]);

  function moveItem(type: AssetType, index: number, direction: -1 | 1) {
    const target = type === "css" ? css : scripts;
    const setTarget = type === "css" ? setCSS : setScripts;
    const nextIndex = index + direction;
    if (nextIndex < 0 || nextIndex >= target.length) {
      return;
    }

    const next = [...target];
    const [item] = next.splice(index, 1);
    if (item === undefined) {
      return;
    }
    next.splice(nextIndex, 0, item);
    setTarget(next);
  }

  function removeItem(type: AssetType, index: number) {
    if (type === "css") {
      setCSS((prev) => prev.filter((_, i) => i !== index));
    } else {
      setScripts((prev) => prev.filter((_, i) => i !== index));
    }
  }

  async function onUpload(file: File) {
    setIsUploading(true);
    await handle(
      async () => {
        const asset = await uploadThemeAsset(file);
        appendUploadedAsset(asset);
      },
      {
        promiseToast: {
          loading: "Uploading theme asset...",
          success: "Theme asset uploaded",
        },
      },
    );
    setIsUploading(false);
  }

  function appendUploadedAsset(asset: Asset) {
    const mime = asset.mime_type;
    if (mime === "text/css") {
      setCSS((prev) => [...prev, asset.path]);
      return;
    }

    if (mime === "application/javascript" || mime === "text/javascript") {
      setScripts((prev) => [...prev, asset.path]);
    }
  }

  async function onSave() {
    setIsSaving(true);
    await handle(
      async () => {
        await updateSettings({
          metadata: {
            ...settings.metadata,
            theme: {
              css,
              scripts,
            },
          },
        });
      },
      {
        promiseToast: {
          loading: "Saving theme settings...",
          success: "Theme settings saved",
        },
        cleanup: async () => {
          await revalidate();
        },
      },
    );
    setIsSaving(false);
  }

  return {
    css,
    scripts,
    hasChanges,
    isUploading,
    isSaving,
    onUpload,
    onSave,
    removeItem,
    moveItem,
  };
}
