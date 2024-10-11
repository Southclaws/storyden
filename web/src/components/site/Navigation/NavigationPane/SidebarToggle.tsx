"use client";

import { Sidebar } from "src/components/graphics/Sidebar/Sidebar";

import { Button } from "@/components/ui/button";

import { useSidebar } from "./useSidebar";

type Props = {
  initialValue: boolean;
};

export function SidebarToggle({ initialValue }: Props) {
  const { setShowLeftBar, showLeftBar } = useSidebar(initialValue);

  return (
    <Button size="md" p="0" variant="ghost" onClick={setShowLeftBar}>
      <Sidebar open={showLeftBar} />
    </Button>
  );
}
