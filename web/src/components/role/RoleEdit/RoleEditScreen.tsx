import { DeleteWithConfirmationButton } from "@/components/site/DeleteConfirmationButton";
import { ColourPickerField } from "@/components/ui/ColourPickerField";
import { FormControl } from "@/components/ui/FormControl";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { FormLabel } from "@/components/ui/FormLabel";
import { Button } from "@/components/ui/button";
import { CardGroupSelect } from "@/components/ui/form/CardGroupSelect";
import { Input } from "@/components/ui/input";
import {
  PermissionList,
  buildPermissionList,
} from "@/lib/permission/permission";
import {
  isDefaultRole,
  isEditableDefaultRole,
  isGuestRole,
  isStoredDefaultRole,
  readPermissions,
} from "@/lib/role/defaults";
import { LStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { Props, useRoleEditScreen } from "./useRoleEdit";

export function RoleEditScreen(props: Props) {
  const {
    form,
    handlers: { handleSave, handleDelete, handleReset },
  } = useRoleEditScreen(props);

  const permissionList = isGuestRole(props.role)
    ? buildPermissionList(...readPermissions)
    : PermissionList;

  const isStored = isStoredDefaultRole(props.role);
  const isDefault = isDefaultRole(props.role);

  return (
    <styled.form
      className={lstack()}
      h="full"
      justifyContent="space-between"
      onSubmit={handleSave}
    >
      <LStack px="0.5" maxH="full" pb="1" overflowY="scroll">
        <FormControl>
          <FormLabel>Name</FormLabel>
          <Input {...form.register("name")} />
          <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Colour</FormLabel>
          <ColourPickerField control={form.control} name="colour" />
          <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
        </FormControl>

        <FormControl>
          <FormLabel>Permissions</FormLabel>
          <CardGroupSelect
            control={form.control}
            name="permissions"
            items={permissionList.map((p) => ({
              value: p.value,
              label: p.name,
              description: p.description,
            }))}
          />
          <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
        </FormControl>
      </LStack>

      <WStack>
        {isDefault ? (
          <>
            {isStored ? (
              <DeleteWithConfirmationButton
                type="button"
                flexGrow="1"
                onDelete={handleReset}
              >
                Reset
              </DeleteWithConfirmationButton>
            ) : null}
          </>
        ) : (
          <DeleteWithConfirmationButton
            type="button"
            flexGrow="1"
            onDelete={handleDelete}
          />
        )}
        <Button
          flexGrow="1"
          variant="outline"
          disabled={!form.formState.isDirty}
        >
          Save
        </Button>
      </WStack>
    </styled.form>
  );
}
