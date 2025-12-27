import { SystemIcon } from "@/components/ui/icons/Robot";
import { HStack, styled } from "@/styled-system/jsx";

type Props = {
  size?: "xs" | "sm" | "md" | "lg";
  name?: "hidden" | "visible";
};

const sizeMap = {
  xs: "xs",
  sm: "sm",
  md: "md",
  lg: "lg",
} as const;

export function SystemBadge({ size = "sm", name = "visible" }: Props) {
  return (
    <HStack
      className="system-badge__container"
      minW="0"
      w="full"
      alignItems="center"
      overflow="hidden"
      gap={size === "lg" ? "2" : "1"}
    >
      <styled.div
        display="flex"
        alignItems="center"
        justifyContent="center"
        borderRadius="md"
        bg="bg.muted"
        w={size === "xs" ? "4" : size === "sm" ? "6" : size === "md" ? "8" : "10"}
        h={size === "xs" ? "4" : size === "sm" ? "6" : size === "md" ? "8" : "10"}
      >
        <SystemIcon
          size={size === "xs" ? 12 : size === "sm" ? 16 : size === "md" ? 20 : 24}
        />
      </styled.div>
      {name === "visible" && (
        <styled.span
          fontSize={sizeMap[size]}
          color="fg.subtle"
          fontStyle="italic"
          fontWeight="normal"
        >
          System
        </styled.span>
      )}
    </HStack>
  );
}
