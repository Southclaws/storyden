import { Portal } from "@ark-ui/react";

import { useDisclosure } from "src/utils/useDisclosure";

import { Permission } from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { ButtonProps } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { CreateIcon } from "@/components/ui/icons/Create";
import { Item } from "@/components/ui/menu";
import { hasPermission } from "@/utils/permissions";

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
  const { onOpen, isOpen, onClose } = useDisclosure();

  if (!hasPermission(session, Permission.MANAGE_CATEGORIES)) {
    return null;
  }

  return (
    <>
      <IconButton
        type="button"
        size="xs"
        variant="ghost"
        px={hideLabel ? "0" : "1"}
        onClick={onOpen}
        {...props}
      >
        {CreateCategoryIcon}
        {!hideLabel && (
          <>
            <span>{CreateCategoryLabel}</span>
          </>
        )}
      </IconButton>

      <CategoryCreateModal isOpen={isOpen} onClose={onClose} {...props} />
    </>
  );
}

export function CreateCategoryMenuItem({ hideLabel }: Props) {
  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <Item value={CreateCategoryID} onClick={onOpen}>
      {CreateCategoryIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{CreateCategoryLabel}</span>
        </>
      )}

      <Portal>
        {/* Portal to avoid nested form triggering. */}
        <CategoryCreateModal isOpen={isOpen} onClose={onClose} />
      </Portal>
    </Item>
  );
}
