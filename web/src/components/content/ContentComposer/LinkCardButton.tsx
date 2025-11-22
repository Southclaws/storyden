"use client";

import { Editor } from "@tiptap/react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import * as Popover from "@/components/ui/popover";
import { HStack } from "@/styled-system/jsx";

type LinkCardButtonProps = {
  editor: Editor;
};

export function LinkCardButton({ editor }: LinkCardButtonProps) {
  const [url, setUrl] = useState("");
  const [open, setOpen] = useState(false);

  const handleSetLinkCard = () => {
    if (url === "") {
      return;
    }

    (editor.chain().focus() as any).setLinkCard({ href: url }).run();
    setOpen(false);
    setUrl("");
  };

  return (
    <Popover.Root open={open} onOpenChange={(details) => setOpen(details.open)}>
      <Popover.Trigger asChild>
        <Button
          type="button"
          size="xs"
          variant="ghost"
          title="Insert link card"
          onClick={() => setOpen(true)}
        >
          🔗
        </Button>
      </Popover.Trigger>

      <Popover.Positioner>
        <Popover.Content>
          <HStack gap="1" alignItems="stretch">
            <Input
              size="xs"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="Enter URL for card"
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault();
                  handleSetLinkCard();
                }
                if (e.key === "Escape") {
                  setOpen(false);
                }
              }}
              autoFocus
            />
            <Button type="button" size="xs" onClick={handleSetLinkCard}>
              Insert
            </Button>
          </HStack>
        </Popover.Content>
      </Popover.Positioner>
    </Popover.Root>
  );
}
