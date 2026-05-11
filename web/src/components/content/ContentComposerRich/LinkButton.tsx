"use client";

import { Editor } from "@tiptap/react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { LinkIcon } from "@/components/ui/icons/Typography";
import { Input } from "@/components/ui/input";
import * as Popover from "@/components/ui/popover";
import { useI18n } from "@/i18n/provider";
import { isValidLinkLike, normalizeLink } from "@/lib/link/validation";
import { HStack } from "@/styled-system/jsx";

type LinkButtonProps = {
  editor: Editor;
};

export function LinkButton({ editor }: LinkButtonProps) {
  const { t } = useI18n();
  const [url, setUrl] = useState("");
  const [open, setOpen] = useState(false);
  const [isInvalid, setIsInvalid] = useState(false);

  const isActive = editor.isActive("link");
  const currentUrl = editor.getAttributes("link")["href"] || "";

  const handleOpen = () => {
    if (isActive) {
      setUrl(currentUrl);
    } else {
      setUrl("");
    }
    setOpen(true);
  };

  const handleChangeURL = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;

    if (value === "") {
      setIsInvalid(false);
    } else {
      setIsInvalid(isValidLinkLike(value) === false);
    }

    setUrl(value);
  };

  const handleSetLink = () => {
    const trimmedUrl = url.trim();

    if (!trimmedUrl) {
      if (isActive) {
        editor.chain().focus().extendMarkRange("link").unsetLink().run();
      }
      setOpen(false);
      setUrl("");
      return;
    }

    const normalizedUrl = normalizeLink(trimmedUrl);

    if (!normalizedUrl) {
      // Keep popover open.
      setIsInvalid(true);
      return;
    }

    if (isActive) {
      editor
        .chain()
        .focus()
        .extendMarkRange("link")
        .setLink({ href: normalizedUrl })
        .run();
    } else {
      editor.chain().focus().setLink({ href: normalizedUrl }).run();
    }

    setOpen(false);
    setUrl("");
  };

  const handleRemoveLink = () => {
    editor.chain().focus().extendMarkRange("link").unsetLink().run();
    setOpen(false);
    setUrl("");
  };

  return (
    <Popover.Root open={open} onOpenChange={(details) => setOpen(details.open)}>
      <Popover.Trigger asChild>
        <Button
          type="button"
          size="xs"
          variant={isActive ? "subtle" : "ghost"}
          title={isActive ? t("Edit link") : t("Add link")}
          onClick={handleOpen}
        >
          <LinkIcon />
        </Button>
      </Popover.Trigger>

      <Popover.Positioner>
        <Popover.Content>
          <HStack gap="1" alignItems="stretch">
            <Input
              borderColor={isInvalid ? "border.error" : undefined}
              size="xs"
              value={url}
              onChange={handleChangeURL}
              placeholder={t("Enter or paste URL")}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                  handleSetLink();
                }
                if (e.key === "Escape") {
                  setOpen(false);
                }
              }}
              autoFocus
              aria-label={t("Link URL")}
            />
            <HStack gap="2" justifyContent="flex-end">
              {isActive && (
                <Button
                  type="button"
                  size="xs"
                  variant="ghost"
                  onClick={handleRemoveLink}
                  title={t("Remove link")}
                >
                  <DeleteIcon />
                </Button>
              )}
              <Button type="button" size="xs" onClick={handleSetLink}>
                {isActive ? t("Update") : t("Add")}
              </Button>
            </HStack>
          </HStack>
        </Popover.Content>
      </Popover.Positioner>
    </Popover.Root>
  );
}
