import { Portal } from "@ark-ui/react";

import { Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { Button, ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CreateIcon } from "@/components/ui/icons/Create";
import { Item } from "@/components/ui/menu";
import { hasPermission } from "@/utils/permissions";
import { useDisclosure } from "@/utils/useDisclosure";

import { CategoryCreateModal } from "./CategoryCreateModal";

type Props = ButtonProps & {
  parentSlug?: string;
  hideLabel?: boolean;
};

export const CreateCategoryID = "create-category";
export const CreateCategoryLabel = "Create";
export const CreateCategoryIcon = <CreateIcon />;

export function CategoryCreateTrigger({
  parentSlug,
  hideLabel,
  ...props
}: Props) {
  const session = useSession();
  const useDisclosureProps = useDisclosure();
  const Trigger = hideLabel ? IconButton : Button;

  if (!hasPermission(session, Permission.MANAGE_CATEGORIES)) {
    return null;
  }

  return (
    <>
      <Trigger
        type="button"
        variant="ghost"
        aria-label={hideLabel ? CreateCategoryLabel : undefined}
        onClick={useDisclosureProps.onOpen}
        {...props}
      >
        {CreateCategoryIcon}
        {!hideLabel && (
          <>
            <span>{CreateCategoryLabel}</span>
          </>
        )}
      </Trigger>

      <CategoryCreateModal {...useDisclosureProps} defaultParent={parentSlug} />
    </>
  );
}

export function CreateCategoryMenuItem({ hideLabel }: Props) {
  const useDisclosureProps = useDisclosure();

  return (
    <Item
      value={CreateCategoryID}
      aria-label={hideLabel ? CreateCategoryLabel : undefined}
      onClick={useDisclosureProps.onOpen}
    >
      {CreateCategoryIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{CreateCategoryLabel}</span>
        </>
      )}

      <Portal>
        {/* Portal to avoid nested form triggering. */}
        <CategoryCreateModal {...useDisclosureProps} />
      </Portal>
    </Item>
  );
}
