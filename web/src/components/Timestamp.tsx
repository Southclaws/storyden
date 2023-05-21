import { LinkProps, Text } from "@chakra-ui/react";
import { Anchor } from "./site/Anchor";

type Props = {
  created: string;
  updated?: string | undefined;
} & LinkProps;

export function Timestamp({ created, updated, ...props }: Props) {
  return (
    <Text as="span" px={2}>
      {props.href ? (
        <Anchor href={props.href}>{created} ago</Anchor>
      ) : (
        <Text as="span">{created}</Text>
      )}
      {updated && <> (updated {updated} ago)</>}
    </Text>
  );
}
