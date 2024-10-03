import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { usePathname } from "next/navigation";

import { Category } from "src/api/openapi-schema";
import { Anchor } from "src/components/site/Anchor";

import { css } from "@/styled-system/css";
import { Box, HStack, styled } from "@/styled-system/jsx";

import { CategoryEdit } from "../CategoryEdit/CategoryEdit";

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

  const href = `/d/${props.slug}`;
  const selected = href === pathname;

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <HStack
      className="group"
      id="category-list-item"
      style={style}
      key={props.id}
      ref={setNodeRef}
      w="full"
      height="8"
      p="1"
      justifyContent="space-between"
      borderRadius="md"
      bgColor={selected ? "gray.a2" : undefined}
      _hover={{
        backgroundColor: "gray.a2",
      }}
      cursor={isDragging ? "grabbing" : "grab"}
    >
      <Anchor
        href={href}
        className={css({
          w: "full",
          _hover: { textDecoration: "none" },

          // Disable pointer events that trigger navigation while dragging.
          ...(isDragging && {
            pointerEvents: "none",
          }),
        })}
        // Only listen for drag events around the actual category name. This
        // prevents the item from being dragged when clicking the edit button.
        {...attributes}
        {...listeners}
      >
        <styled.h2 role="navigation" w="full" fontWeight="medium" fontSize="xs">
          {props.name}
        </styled.h2>
      </Anchor>

      {props.isAdmin && (
        <Box
          display={{
            base: "none",
            _groupHover: "block",
          }}
        >
          <CategoryEdit {...props} />
        </Box>
      )}
    </HStack>
  );
}
