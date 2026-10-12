"use client";

import { Title as DialogTitle } from "@radix-ui/react-dialog";
import { Command } from "cmdk";
import dynamic from "next/dynamic";

import { styled } from "@/styled-system/jsx";

import { UnreadyBanner } from "../Unready";

import "./styles.css";

import { CommandPaletteProvider, useCommandPalette } from "./Context";

const CommandPaletteDialogContent = dynamic(
  () =>
    import("./CommandPaletteDialogContent").then(
      (module) => module.CommandPaletteDialogContent,
    ),
  { ssr: false, loading: () => <UnreadyBanner /> },
);

export function CommandPalette() {
  return (
    <CommandPaletteProvider>
      <CommandPaletteDialog />
    </CommandPaletteProvider>
  );
}

function CommandPaletteDialog() {
  const { open, dialogRef } = useCommandPalette();

  return (
    <Command.Dialog
      ref={dialogRef}
      open={open}
      label="Command Menu"
      aria-description="The command palette allows you to quickly navigate and perform actions in Storyden."
    >
      <DialogTitle asChild>
        <styled.h2 srOnly>Command Menu</styled.h2>
      </DialogTitle>
      {open && <CommandPaletteDialogContent />}
    </Command.Dialog>
  );
}
