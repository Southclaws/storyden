"use client";

import { Portal } from "@ark-ui/react";
import Link from "next/link";

import { CheckIcon } from "@/components/ui/icons/Check";
import { ChevronDownIcon } from "@/components/ui/icons/Chevron";
import * as Popover from "@/components/ui/popover";
import { cx } from "@/styled-system/css";
import { sectionNavigation } from "@/styled-system/recipes";

export type SectionNavigationItem = {
  href: string;
  label: string;
};

export type SectionNavigationGroup = {
  label?: string;
  items: readonly SectionNavigationItem[];
};

export type SectionNavigationInternalProps = {
  activeHref: string;
  className?: string;
  groups: readonly SectionNavigationGroup[];
  label: string;
  onNavigate: () => void;
  onOpenChange: (open: boolean) => void;
  open: boolean;
};

export function SectionNavigationInternal({
  activeHref,
  className,
  groups,
  label,
  onNavigate,
  onOpenChange,
  open,
}: SectionNavigationInternalProps) {
  const styles = sectionNavigation();
  const items = groups.flatMap((group) => group.items);
  const activeItem =
    items.find((item) => item.href === activeHref) ?? items.at(0);

  if (!activeItem) {
    return null;
  }

  return (
    <div className={cx(styles.root, className)}>
      <Popover.Root
        open={open}
        onOpenChange={({ open: nextOpen }) => onOpenChange(nextOpen)}
        positioning={{ placement: "bottom-start", sameWidth: true }}
      >
        <Popover.Trigger
          className={styles.trigger}
          aria-label={`${label}: ${activeItem.label}. Choose section`}
        >
          <span className={styles.triggerLabel}>{activeItem.label}</span>
          <span className={styles.triggerIndicator} aria-hidden="true">
            <ChevronDownIcon />
          </span>
        </Popover.Trigger>

        <Portal>
          <Popover.Positioner className={styles.positioner}>
            <Popover.Content className={styles.content} aria-label={label}>
              <nav className={styles.navigation} aria-label={label}>
                {groups.map((group, groupIndex) => (
                  <section
                    className={styles.section}
                    key={group.label ?? groupIndex}
                  >
                    {group.label && (
                      <div className={styles.sectionLabel}>{group.label}</div>
                    )}
                    <ul className={styles.items}>
                      {group.items.map((item) => {
                        const active = item.href === activeHref;

                        return (
                          <li className={styles.item} key={item.href}>
                            <Link
                              aria-current={active ? "page" : undefined}
                              className={styles.link}
                              href={item.href}
                              onClick={onNavigate}
                            >
                              <span>{item.label}</span>
                              {active && (
                                <span
                                  className={styles.linkIndicator}
                                  aria-hidden="true"
                                >
                                  <CheckIcon />
                                </span>
                              )}
                            </Link>
                          </li>
                        );
                      })}
                    </ul>
                  </section>
                ))}
              </nav>
            </Popover.Content>
          </Popover.Positioner>
        </Portal>
      </Popover.Root>
    </div>
  );
}
