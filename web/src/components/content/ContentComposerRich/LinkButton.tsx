"use client";

import { Editor } from "@tiptap/react";
import { useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { LinkIcon } from "@/components/ui/icons/Typography";
import { Input } from "@/components/ui/input";
import * as Popover from "@/components/ui/popover";
import { normalizeLink } from "@/lib/link/validation";
import { HStack } from "@/styled-system/jsx";

type LinkButtonProps = {
  editor: Editor;
};

export function LinkButton({ editor }: LinkButtonProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [defaultUrl, setDefaultUrl] = useState("");
  const [open, setOpen] = useState(false);
  const [isInvalid, setIsInvalid] = useState(false);

  const isActive = editor.isActive("link");
  const currentUrl = editor.getAttributes("link")["href"] || "";

  useEffect(() => {
    if (open && inputRef.current) {
      inputRef.current.value = defaultUrl;
    }
  }, [defaultUrl, open]);

  const handleOpen = () => {
    if (isActive) {
      setDefaultUrl(currentUrl);
    } else {
      setDefaultUrl("");
    }
    setOpen(true);
  };

  const handleSetLink = () => {
    const trimmedUrl = inputRef.current?.value.trim() ?? "";

    if (!trimmedUrl) {
      if (isActive) {
        editor.chain().focus().extendMarkRange("link").unsetLink().run();
      }
      setOpen(false);
      setDefaultUrl("");
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
    setDefaultUrl("");
  };

  const handleRemoveLink = () => {
    editor.chain().focus().extendMarkRange("link").unsetLink().run();
    setOpen(false);
    setDefaultUrl("");
  };

  return (
    <Popover.Root open={open} onOpenChange={(details) => setOpen(details.open)}>
      <Popover.Trigger asChild>
        <Button
          type="button"
          variant={isActive ? "subtle" : "ghost"}
          title={isActive ? "Edit link" : "Add link"}
          onClick={handleOpen}
        >
          <LinkIcon />
        </Button>
      </Popover.Trigger>

      <Popover.Positioner>
        <Popover.Content>
          <HStack gap="1" alignItems="stretch">
            <Input
              key={defaultUrl}
              ref={inputRef}
              borderColor={isInvalid ? "status.danger.border" : undefined}
              defaultValue={defaultUrl}
              placeholder="Enter or paste URL"
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
              aria-label="Link URL"
            />
            <HStack gap="2" justifyContent="flex-end">
              {isActive && (
                <Button
                  type="button"
                  variant="ghost"
                  onClick={handleRemoveLink}
                  title="Remove link"
                >
                  <DeleteIcon />
                </Button>
              )}
              <Button type="button" onClick={handleSetLink}>
                {isActive ? "Update" : "Add"}
              </Button>
            </HStack>
          </HStack>
        </Popover.Content>
      </Popover.Positioner>
    </Popover.Root>
  );
}
