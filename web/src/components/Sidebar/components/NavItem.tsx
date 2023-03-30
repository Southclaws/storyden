import { LinkProps } from "@chakra-ui/react";
import { useRouter } from "next/router";
import { Anchor } from "src/components/site/Anchor";

export function NavItem(props: LinkProps) {
  const { asPath } = useRouter();
  const selected = props.href === asPath;

  return (
    <Anchor
      borderRadius="md"
      p={2}
      bgColor={selected ? "blackAlpha.100" : ""}
      _hover={{
        backgroundColor: "blackAlpha.50",
      }}
      {...props}
    >
      {props.children}
    </Anchor>
  );
}
