import { MenuSelectionDetails, Portal } from "@ark-ui/react";

import { Category, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { MoreAction } from "@/components/site/Action/More";
import { Heading } from "@/components/ui/heading";
import * as Menu from "@/components/ui/menu";
import { WEB_ADDRESS } from "@/config";
import { styled } from "@/styled-system/jsx";
import { useShare } from "@/utils/client";
import { hasPermission } from "@/utils/permissions";
import { useCopyToClipboard } from "@/utils/useCopyToClipboard";

import { CategoryCreateMenuItem } from "../CategoryCreate/CategoryCreateMenuItem";
import { CategoryDeleteMenuItem } from "../CategoryDelete/CategoryDeleteMenuItem";
import { CategoryEditMenuItem } from "../CategoryEdit/CategoryEdit";

type Props = {
  category: Category;
};

export function useCategoryMenu({ category }: Props) {
  const account = useSession();
  const [, copyToClipboard] = useCopyToClipboard();

  const isEditingEnabled = hasPermission(account, Permission.MANAGE_CATEGORIES);

  const isSharingEnabled = useShare();

  const permalink = `${WEB_ADDRESS}/d/${category.slug}`;

  async function handleCopyLink() {
    await copyToClipboard(permalink);
  }

  async function handleShare() {
    await navigator.share({
      title: `Discussion category: ${category.name}`,
      url: permalink,
      text: category.description,
    });
  }

  async function handleSelect({ value }: MenuSelectionDetails) {
    switch (value) {
      case "copy-link":
        return handleCopyLink();

      case "share":
        return handleShare();

      case "edit":
        // Handled by item component
        return;

      case "create-subcategory":
        // Handled by item component
        return;

      case "delete":
        // Handled by item component
        return;
    }
  }

  return {
    isSharingEnabled,
    isEditingEnabled,
    handlers: {
      handleSelect,
    },
  };
}

export function CategoryMenu(props: Props) {
  const { isSharingEnabled, isEditingEnabled, handlers } =
    useCategoryMenu(props);

  const { category } = props;

  return (
    <Menu.Root onSelect={handlers.handleSelect}>
      <Menu.Trigger asChild>
        <MoreAction size="xs" />
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="48" userSelect="none">
            <Menu.ItemGroup id="account">
              <Menu.ItemGroupLabel>
                <Heading size="sm">{category.name}</Heading>
                <styled.span color="fg.subtle">discussion category</styled.span>
              </Menu.ItemGroupLabel>

              <Menu.Separator />

              <Menu.Item value="copy-link">Copy link</Menu.Item>
              {isSharingEnabled && <Menu.Item value="share">Share</Menu.Item>}
            </Menu.ItemGroup>

            {isEditingEnabled && (
              <>
                <Menu.Separator />

                <Menu.ItemGroup id="manage">
                  <CategoryCreateMenuItem parentCategory={category} />
                  <CategoryEditMenuItem {...category} />
                  <CategoryDeleteMenuItem {...category} />
                </Menu.ItemGroup>
              </>
            )}
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}
