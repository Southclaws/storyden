import Link from "next/link";
import React from "react";

import { LinkButton } from "@/components/ui/link-button";
import { LinkButtonStyleProps } from "@/components/ui/link-button";
import { Item } from "@/components/ui/menu";
import { cx } from "@/styled-system/css";

type Props = {
  id: string;
  route: string | (() => string);
  icon: React.ReactElement;
  label: string;
  className?: string;
} & AnchorProps &
  LinkButtonStyleProps;

export type AnchorProps = {
  hideLabel?: boolean;
};

export function Anchor({
  id,
  route,
  icon,
  label,
  hideLabel,
  className,
  ...props
}: Props) {
  const href = typeof route === "function" ? route() : route;

  return (
    <LinkButton
      href={href}
      size="xs"
      p="1"
      variant="ghost"
      className={cx("navigation-anchor", `navigation-anchor--${id}`, className)}
      data-navigation-anchor={id}
      {...props}
    >
      <span className="navigation-anchor__icon" aria-hidden="true">
        {React.cloneElement(icon, {
          width: "1.5rem",
        } as any)}
      </span>
      {!hideLabel && <span className="navigation-anchor__label">{label}</span>}
    </LinkButton>
  );
}

export function MenuItem({
  id,
  route,
  icon,
  label,
  hideLabel,
  ...props
}: Props) {
  const href = typeof route === "function" ? route() : route;

  return (
    <Link
      className={cx("navigation-menu-anchor", `navigation-menu-anchor--${id}`)}
      data-navigation-anchor={id}
      href={href}
    >
      <Item value={id}>
        <span className="navigation-menu-anchor__icon" aria-hidden="true">
          {icon}
        </span>
        {!hideLabel && (
          <span className="navigation-menu-anchor__label">{label}</span>
        )}
      </Item>
    </Link>
  );
}
