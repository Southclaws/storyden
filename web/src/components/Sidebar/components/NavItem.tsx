import { Button, LinkProps, Text } from "@chakra-ui/react";
import { Anchor } from "src/components/site/Anchor";

type Props = LinkProps & {
  selected: boolean;
};

export function NavItem(props: Props) {
  return (
    <Anchor {...props}>
      <Button
        as="span"
        role="navigation"
        variant="ghost"
        bgColor={props.selected ? "blackAlpha.200" : ""}
        w="full"
      >
        {props.children}
      </Button>
    </Anchor>
  );
}
