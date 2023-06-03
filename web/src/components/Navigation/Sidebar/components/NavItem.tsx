import { LinkProps } from "@chakra-ui/react";
import { usePathname } from "next/navigation";
import { Anchor } from "src/components/site/Anchor";

export function NavItem(props: LinkProps) {
  const pathname = usePathname();
  const selected = props.href === pathname;

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
