"use client";

import { usePathname } from "next/navigation";
import { useState } from "react";

import {
  type SectionNavigationGroup,
  SectionNavigationInternal,
} from "./SectionNavigation.internal";

export type SectionNavigationProps = {
  activeHref?: string;
  className?: string;
  defaultOpen?: boolean;
  groups: readonly SectionNavigationGroup[];
  label: string;
};

export function SectionNavigation({
  activeHref,
  className,
  defaultOpen = false,
  groups,
  label,
}: SectionNavigationProps) {
  const pathname = usePathname();
  const [open, setOpen] = useState(defaultOpen);

  return (
    <SectionNavigationInternal
      activeHref={activeHref ?? pathname}
      className={className}
      groups={groups}
      label={label}
      open={open}
      onNavigate={() => setOpen(false)}
      onOpenChange={setOpen}
    />
  );
}
