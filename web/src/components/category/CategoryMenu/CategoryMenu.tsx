import { MenuSelectionDetails, Portal } from "@ark-ui/react";

import { type Account, Category, Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { MoreAction } from "@/components/site/Action/More";
import { ButtonProps } from "@/components/ui/button";
import * as Menu from "@/components/ui/menu";
import { Text } from "@/components/ui/text";
import { WEB_ADDRESS } from "@/config";
import { useShare } from "@/utils/client";
import { hasPermission } from "@/utils/permissions";
import { useCopyToClipboard } from "@/utils/useCopyToClipboard";

import { CategoryCreateMenuItem } from "../CategoryCreate/CategoryCreateMenuItem";
import { CategoryDeleteMenuItem } from "../CategoryDelete/CategoryDeleteMenuItem";
import { CategoryEditMenuItem } from "../CategoryEdit/CategoryEdit";

type Props = {
  category: Category;
  initialSession?: Account;
  triggerProps?: ButtonProps;
};

export function useCategoryMenu({ category, initialSession }: Props) {
  const account = useSession(initialSession);
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

  const { category, triggerProps } = props;

  return (
    <Menu.Root onSelect={handlers.handleSelect}>
      <Menu.Trigger asChild>
        <MoreAction {...triggerProps} />
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="48" userSelect="none">
            <Menu.ItemGroup id="account">
              <Menu.ItemGroupLabel>
                <Text
                  variant="supporting"
                  color="text.default"
                  fontWeight="semibold"
                >
                  {category.name}
                </Text>
                <Text as="span" variant="metadata">
                  discussion category
                </Text>
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
