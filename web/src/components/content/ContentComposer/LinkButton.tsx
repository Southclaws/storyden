"use client";

import { Editor } from "@tiptap/react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { LinkIcon } from "@/components/ui/icons/Typography";
import { Input } from "@/components/ui/input";
import * as Popover from "@/components/ui/popover";
import { HStack } from "@/styled-system/jsx";

type LinkButtonProps = {
  editor: Editor;
};

export function LinkButton({ editor }: LinkButtonProps) {
  const [url, setUrl] = useState("");
  const [open, setOpen] = useState(false);

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

  const handleSetLink = () => {
    if (url === "") {
      if (isActive) {
        editor.chain().focus().extendMarkRange("link").unsetLink().run();
      }
    } else {
      if (isActive) {
        editor
          .chain()
          .focus()
          .extendMarkRange("link")
          .setLink({ href: url })
          .run();
      } else {
        editor.chain().focus().setLink({ href: url }).run();
      }
    }
    setOpen(false);
    setUrl("");
  };

  const handleRemoveLink = () => {
    editor.chain().focus().extendMarkRange("link").unsetLink().run();
    setOpen(false);
  };

  return (
    <Popover.Root open={open} onOpenChange={(details) => setOpen(details.open)}>
      <Popover.Trigger asChild>
        <Button
          type="button"
          size="xs"
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
              size="xs"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
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
            />
            <HStack gap="2" justifyContent="flex-end">
              {isActive && (
                <Button
                  type="button"
                  size="xs"
                  variant="ghost"
                  onClick={handleRemoveLink}
                  title="Remove link"
                >
                  <DeleteIcon />
                </Button>
              )}
              <Button type="button" size="xs" onClick={handleSetLink}>
                {isActive ? "Update" : "Add"}
              </Button>
            </HStack>
          </HStack>
        </Popover.Content>
      </Popover.Positioner>
    </Popover.Root>
  );
}
