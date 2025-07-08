"use client";

import { Command } from "cmdk";
import { useEffect, useRef, useState } from "react";

import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { IntelligenceIcon } from "@/components/ui/icons/Intelligence";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { LinkIcon } from "@/components/ui/icons/Link";

import "./styles.css";

export function CommandPalette() {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const dialogRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((open) => !open);
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  useEffect(() => {
    function handleClickOutside(event) {
      if (
        dialogRef.current &&
        !dialogRef.current.contains(event.target) &&
        // Only do outside click handling if the input is empty.
        search === ""
      ) {
        setOpen(false);
      }
    }

    if (open) {
      document.addEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [open, setOpen, search]);

  return (
    <Command.Dialog
      ref={dialogRef}
      open={open}
      label="Command Menu"
      onKeyDown={(e) => {
        if (e.key === "Escape") {
          e.preventDefault();
          setOpen(false);
        }
      }}
      aria-description="The command palette allows you to quickly navigate and perform actions in Storyden."
    >
      <Command.Input
        value={search}
        onValueChange={setSearch}
        placeholder="Ask, search or command..."
      />
      <Command.List>
        <Command.Item>
          <DiscussionIcon />
          Post a new thread
        </Command.Item>
        <Command.Item>
          <LibraryIcon />
          Create a new page
        </Command.Item>
        <Command.Item>
          <LinkIcon />
          Share a link
        </Command.Item>
        <Command.Item>
          <IntelligenceIcon />
          Ask AI...
        </Command.Item>
      </Command.List>
    </Command.Dialog>
  );
}
