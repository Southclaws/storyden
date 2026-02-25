"use client";

import { Sidebar } from "src/components/graphics/Sidebar/Sidebar";

import { Button } from "@/components/ui/button";

import { useSidebar } from "./useSidebar";

type Props = {
  initialValue: boolean;
};

export function SidebarToggle({ initialValue }: Props) {
  const { setShowLeftBar, showLeftBar } = useSidebar(initialValue);
  const isOpen = showLeftBar;
  const label = isOpen
    ? "Close navigation sidebar"
    : "Open navigation sidebar";

  return (
    <Button
      type="button"
      size="md"
      p="0"
      variant="ghost"
      onClick={setShowLeftBar}
      aria-label={label}
      title={label}
      aria-expanded={isOpen}
      aria-controls="navigation__leftbar navigation__rightbar"
      aria-pressed={isOpen}
      data-state={isOpen ? "open" : "closed"}
    >
      <Sidebar open={isOpen} aria-hidden="true" focusable="false" />
    </Button>
  );
}
