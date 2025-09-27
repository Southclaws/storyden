import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { usePathname } from "next/navigation";

import { Category } from "src/api/openapi-schema";
import { Anchor } from "src/components/site/Anchor";

import { css } from "@/styled-system/css";
import { Box, HStack, styled } from "@/styled-system/jsx";

import { CategoryMenu } from "../CategoryMenu/CategoryMenu";

type Props = {
  category: Category;
  isAdmin: boolean;
};

export function CategoryListItem({ category, isAdmin }: Props) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: category.id });
  const pathname = usePathname();

  const href = `/d/${category.slug}`;
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
      key={category.id}
      ref={setNodeRef}
      w="full"
      height="8"
      py="1"
      px="2"
      justifyContent="space-between"
      borderRadius="md"
      bgColor={selected ? "bg.selected" : undefined}
      _hover={{
        backgroundColor: "bg.emphasized",
        color: "fg.emphasized",
      }}
      color={selected ? "fg.selected" : "fg.subtle"}
      cursor={isDragging ? "grabbing" : undefined}
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
          {category.name}
        </styled.h2>
      </Anchor>

      {isAdmin && (
        <Box
          // Only show the edit button when hovering over the item.
          // Why not display: none? Well, the menu anchor must be mounted on the
          // DOM in order to calculate the position of the menu correctly.
          // Group-hover is a Panda CSS feature that applies styles to a child
          // when the parent of className="group" is hovered.
          opacity={{
            base: "0",
            _groupHover: "full",
          }}
        >
          <CategoryMenu category={category} />
        </Box>
      )}
    </HStack>
  );
}
