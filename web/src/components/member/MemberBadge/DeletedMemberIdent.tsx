import { DeletedMemberIcon } from "@/components/ui/icons/DeletedMember";
import { Box, HStack, styled } from "@/styled-system/jsx";

export type Props = {
  size: "xs" | "sm" | "md" | "lg";
};

function iconSize(size: Props["size"]) {
  switch (size) {
    case "xs":
      return 14;
    case "sm":
      return 18;
    case "md":
      return 24;
    case "lg":
      return 32;
  }
}

export function DeletedMemberIdent({ size }: Props) {
  const width = iconSize(size);

  return (
    <HStack
      className="deleted-member-ident__container"
      minW="0"
      w="full"
      alignItems="center"
      overflow="hidden"
      gap={size === "lg" ? "2" : "1"}
    >
      <Box flexShrink="0" style={{ width, height: width }}>
        <DeletedMemberIcon w="full" h="full" color="fg.muted" />
      </Box>
      <styled.p
        className="deleted-member-ident__label"
        fontSize={size}
        fontWeight={size === "lg" ? "medium" : "normal"}
        overflowX="hidden"
        overflowY="clip"
        textWrap="nowrap"
        textOverflow="ellipsis"
        lineHeight="tight"
        color="fg.subtle"
      >
        @deleted
      </styled.p>
    </HStack>
  );
}
