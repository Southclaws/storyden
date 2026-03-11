"use client";

import { ChangeEvent, useRef } from "react";

import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { Heading } from "@/components/ui/heading";
import { CardBox, HStack, VStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { Props, useThemeSettings } from "./useThemeSettings";

export function ThemeSettingsForm(props: Props) {
  const {
    css,
    scripts,
    hasChanges,
    isUploading,
    isSaving,
    onUpload,
    onSave,
    removeItem,
    moveItem,
  } = useThemeSettings(props);

  const fileRef = useRef<HTMLInputElement | null>(null);

  async function handleFileChange(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }

    await onUpload(file);

    event.target.value = "";
  }

  return (
    <CardBox className={lstack()} gap="4">
      <WStack>
        <Heading size="md">Theme settings</Heading>
        <Button type="button" onClick={onSave} loading={isSaving}>
          Save
        </Button>
      </WStack>

      <FormControl>
        <FormLabel>Upload theme asset</FormLabel>
        <HStack>
          <Button
            type="button"
            variant="outline"
            loading={isUploading}
            onClick={() => fileRef.current?.click()}
          >
            Upload CSS or JS
          </Button>
          <styled.input
            ref={fileRef}
            display="none"
            type="file"
            accept=".css,.js,text/css,application/javascript,text/javascript"
            onChange={handleFileChange}
          />
        </HStack>
        <FormHelperText>
          Upload `.css` and `.js` files to `/api/info/theme/assets/*`.
        </FormHelperText>
      </FormControl>

      <ThemeAssetList
        title="Stylesheets"
        items={css}
        type="css"
        onMoveUp={(idx) => moveItem("css", idx, -1)}
        onMoveDown={(idx) => moveItem("css", idx, 1)}
        onRemove={(idx) => removeItem("css", idx)}
      />

      <ThemeAssetList
        title="Scripts"
        items={scripts}
        type="script"
        onMoveUp={(idx) => moveItem("script", idx, -1)}
        onMoveDown={(idx) => moveItem("script", idx, 1)}
        onRemove={(idx) => removeItem("script", idx)}
      />

      <WStack justifyContent="end">
        <Button type="button" onClick={onSave} loading={isSaving}>
          Save
        </Button>
      </WStack>

      {!hasChanges && (
        <FormHelperText>No unsaved theme manifest changes.</FormHelperText>
      )}
    </CardBox>
  );
}

type ThemeAssetListProps = {
  title: string;
  items: string[];
  type: "css" | "script";
  onMoveUp: (index: number) => void;
  onMoveDown: (index: number) => void;
  onRemove: (index: number) => void;
};

function ThemeAssetList({
  title,
  items,
  type,
  onMoveUp,
  onMoveDown,
  onRemove,
}: ThemeAssetListProps) {
  return (
    <FormControl>
      <FormLabel>{title}</FormLabel>
      {items.length === 0 ? (
        <FormHelperText>No {type} assets configured.</FormHelperText>
      ) : (
        <VStack alignItems="stretch" gap="2">
          {items.map((item, index) => (
            <CardBox key={`${type}-${index}-${item}`} p="2">
              <WStack justifyContent="space-between" alignItems="start">
                <styled.code
                  fontSize="xs"
                  overflowWrap="anywhere"
                  whiteSpace="pre-wrap"
                >
                  {item}
                </styled.code>
                <HStack>
                  <Button
                    size="xs"
                    variant="outline"
                    type="button"
                    onClick={() => onMoveUp(index)}
                  >
                    Up
                  </Button>
                  <Button
                    size="xs"
                    variant="outline"
                    type="button"
                    onClick={() => onMoveDown(index)}
                  >
                    Down
                  </Button>
                  <Button
                    size="xs"
                    variant="ghost"
                    type="button"
                    onClick={() => onRemove(index)}
                  >
                    Remove
                  </Button>
                </HStack>
              </WStack>
            </CardBox>
          ))}
        </VStack>
      )}
    </FormControl>
  );
}
