import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { usePathname } from "next/navigation";

import { Category } from "src/api/openapi/schemas";
import { DragHandleIcon } from "src/components/graphics/DragHandleIcon";
import { Anchor } from "src/components/site/Anchor";

import { CategoryEdit } from "../CategoryEdit/CategoryEdit";

import { css } from "@/styled-system/css";
import { Box, HStack, styled } from "@/styled-system/jsx";

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

  const href = `/c/${props.slug}`;
  const selected = href === pathname;

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <Box
      id="category-list-item"
      style={style}
      key={props.id}
      ref={setNodeRef}
      borderRadius="md"
      bgColor={selected ? "blackAlpha.100" : ""}
      _hover={{
        backgroundColor: "blackAlpha.50",
      }}
      w="full"
    >
      <HStack justifyContent="space-between">
        <Anchor
          href={href}
          className={css({
            w: "full",
            _hover: { textDecoration: "none" },
          })}
        >
          <styled.h2 p={2} role="navigation" w="full" fontWeight="bold">
            {props.name}
          </styled.h2>
        </Anchor>

        {props.isAdmin && (
          <HStack gap={0}>
            <CategoryEdit {...props} />

            <Box
              {...attributes}
              {...listeners}
              cursor={isDragging ? "grabbing" : "grab"}
            >
              <DragHandleIcon />
            </Box>
          </HStack>
        )}
      </HStack>
    </Box>
  );
}
