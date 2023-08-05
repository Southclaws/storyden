import { LinkProps, Text } from "@chakra-ui/react";

import { Anchor } from "./site/Anchor";

type Props = {
  created: string;
  updated?: string | undefined;
} & LinkProps;

export function Timestamp({ created, updated, ...props }: Props) {
  return (
    <Text as="span">
      {props.href ? (
        <Anchor href={props.href}>{created}</Anchor>
      ) : (
        <Text as="span">{created}</Text>
      )}
      {updated && <> (updated {updated})</>}
    </Text>
  );
}
