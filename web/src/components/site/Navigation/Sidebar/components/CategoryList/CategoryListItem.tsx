import { Box, HStack, Heading } from "@chakra-ui/react";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { usePathname } from "next/navigation";

import { Category } from "src/api/openapi/schemas";
import { DragHandleIcon } from "src/components/graphics/DragHandleIcon";
import { Anchor } from "src/components/site/Anchor";

export function CategoryListItem(props: Category & { isAdmin: boolean }) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: props.id });
  const pathname = usePathname();

  const href = `/c/${props.name}`;
  const selected = href === pathname;

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <Box
      style={style}
      key={props.id}
      ref={setNodeRef}
      borderRadius="md"
      p={2}
      bgColor={selected ? "blackAlpha.100" : ""}
      _hover={{
        backgroundColor: "blackAlpha.50",
      }}
      w="full"
    >
      <HStack justifyContent="space-between">
        <Anchor href={href} w="full" _hover={{ textDecor: "none" }}>
          <Heading size="sm" role="navigation" variant="ghost" w="full">
            {props.name}
          </Heading>
        </Anchor>

        {props.isAdmin && (
          <Box
            {...attributes}
            {...listeners}
            cursor={isDragging ? "grabbing" : "grab"}
          >
            <DragHandleIcon />
          </Box>
        )}
      </HStack>
    </Box>
  );
}
