import { MenuSelectionDetails, Portal } from "@ark-ui/react";

import { handle } from "@/api/client";
import { useCategoryList } from "@/api/openapi-client/categories";
import { ThreadReference } from "@/api/openapi-schema";
import { Unready } from "@/components/site/Unready";
import { CategoryIcon } from "@/components/ui/icons/Category";
import { SubmenuIcon } from "@/components/ui/icons/Submenu";
import * as Menu from "@/components/ui/menu";
import { useI18n } from "@/i18n/provider";
import { useThreadMutations } from "@/lib/thread/mutation";
import { HStack } from "@/styled-system/jsx";

import { CategoryBadge } from "../CategoryBadge";
import {
  CategoryCreateTrigger,
  CreateCategoryID,
  CreateCategoryMenuItem,
} from "../CategoryCreate/CategoryCreateTrigger";

type Props = {
  thread: ThreadReference;
};

export function useCategoryMoveMenu({ thread }: Props) {
  const { revalidate, updateCategory } = useThreadMutations(thread);
  const { t } = useI18n();

  async function handleSelect({ value }: MenuSelectionDetails) {
    // If the user clicked the create category prompt, only present when there
    // are no categories available, then exit early.
    if (value === CreateCategoryID) {
      return;
    }

    await handle(
      async () => {
        await updateCategory(value);
      },
      {
        promiseToast: {
          loading: t("Moving thread..."),
          success: t("Moved!"),
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
  const { t } = useI18n();

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
          {t("Move")}
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
  const { t } = useI18n();
  const { data, error } = useCategoryList();
  if (!data) {
    return <Unready error={error} />;
  }

  const { categories } = data;

  const hasAnyCategories = categories.length > 0;
  if (!hasAnyCategories) {
    return (
      <Menu.Content minW="48" userSelect="none">
        <Menu.ItemGroup id="move-no-categories">
          <Menu.ItemGroupLabel>
            {t("No categories to move to")}
          </Menu.ItemGroupLabel>

          <CreateCategoryMenuItem />
        </Menu.ItemGroup>
      </Menu.Content>
    );
  }

  return (
    <Menu.Content minW="48" userSelect="none">
      <Menu.ItemGroup id="move">
        <Menu.ItemGroupLabel>{t("Move thread")}</Menu.ItemGroupLabel>

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
