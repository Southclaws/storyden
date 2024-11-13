import { MenuSelectionDetails, Portal } from "@ark-ui/react";

import { handle } from "@/api/client";
import { useCategoryList } from "@/api/openapi-client/categories";
import { ThreadReference } from "@/api/openapi-schema";
import { Unready } from "@/components/site/Unready";
import { CategoryIcon } from "@/components/ui/icons/Category";
import { SubmenuIcon } from "@/components/ui/icons/Submenu";
import * as Menu from "@/components/ui/menu";
import { useThreadMutations } from "@/lib/thread/mutation";
import { HStack } from "@/styled-system/jsx";

import { CategoryBadge } from "../CategoryBadge";

type Props = {
  thread: ThreadReference;
};

export function useCategoryMoveMenu({ thread }: Props) {
  const { revalidate, updateCategory } = useThreadMutations(thread);

  async function handleSelect({ value }: MenuSelectionDetails) {
    await handle(
      async () => {
        await updateCategory(value);
      },
      {
        promiseToast: {
          loading: "Moving thread...",
          success: "Moved!",
        },
        async cleanup() {
          await revalidate();
        },
      },
    );
  }

  return {
    handlers: {
      handleSelect,
    },
  };
}

export function CategoryMoveMenu(props: Props) {
  const { handlers } = useCategoryMoveMenu(props);

  return (
    <Menu.Root
      size="xs"
      //   lazyMount
      positioning={{ placement: "right-start", gutter: -2 }}
      onSelect={handlers.handleSelect}
    >
      <Menu.TriggerItem justifyContent="space-between">
        <HStack gap="1">
          <CategoryIcon />
          Move
        </HStack>
        <SubmenuIcon />
      </Menu.TriggerItem>

      <Portal>
        <Menu.Positioner>
          <LazyLoadedCategoryMoveMenuContent />
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

function LazyLoadedCategoryMoveMenuContent() {
  const { data, error } = useCategoryList();
  if (!data) {
    return <Unready error={error} />;
  }

  const { categories } = data;

  return (
    <Menu.Content minW="48" userSelect="none">
      <Menu.ItemGroup id="account">
        <Menu.ItemGroupLabel>Move thread</Menu.ItemGroupLabel>

        <Menu.Separator />

        {categories.map((c) => {
          return (
            <Menu.Item key={c.id} value={c.id}>
              <CategoryBadge category={c} asLink={false} />
            </Menu.Item>
          );
        })}
      </Menu.ItemGroup>
    </Menu.Content>
  );
}
