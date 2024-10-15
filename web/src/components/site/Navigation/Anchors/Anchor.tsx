import Link from "next/link";
import React from "react";

import { LinkButton } from "@/components/ui/link-button";
import { Item } from "@/components/ui/menu";

type Props = {
  id: string;
  route: string | (() => string);
  icon: React.ReactElement;
  label: string;
} & AnchorProps;

export type AnchorProps = {
  hideLabel?: boolean;
};

export function Anchor({ id, route, icon, label, hideLabel, ...props }: Props) {
  const href = typeof route === "function" ? route() : route;

  return (
    <LinkButton href={href} size="xs" p="1" variant="ghost" {...props}>
      {React.cloneElement(icon, {
        width: "1.5rem",
      } as any)}
      {!hideLabel && (
        <>
          &nbsp;<span>{label}</span>
        </>
      )}
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
    <Link href={href}>
      <Item value={id}>
        {icon}
        {!hideLabel && (
          <>
            &nbsp;<span>{label}</span>
          </>
        )}
      </Item>
    </Link>
  );
}
