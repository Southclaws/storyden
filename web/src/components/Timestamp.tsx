import { styled } from "@/styled-system/jsx";

import { Anchor } from "./site/Anchor";

type Props = {
  created: string;
  updated?: string | undefined;
  href: string;
};

export function Timestamp({ created, updated, ...props }: Props) {
  return (
    <styled.span>
      {props.href ? (
        <Anchor href={props.href}>{created}</Anchor>
      ) : (
        <styled.span>{created}</styled.span>
      )}
      {updated && <> (updated {updated})</>}
    </styled.span>
  );
}
