import { DeleteWithConfirmationButton } from "@/components/site/DeleteConfirmationButton";
import { ColourPickerField } from "@/components/ui/ColourPickerField";
import { FormControl } from "@/components/ui/FormControl";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { FormLabel } from "@/components/ui/FormLabel";
import { Button } from "@/components/ui/button";
import { CardGroupSelect } from "@/components/ui/form/CardGroupSelect";
import { Input } from "@/components/ui/input";
import { PermissionList } from "@/lib/permission/permission";
import { HStack, LStack, styled } from "@/styled-system/jsx";
import { LStack as lstack } from "@/styled-system/patterns";

import { Props, useRoleEditScreen } from "./useRoleEdit";

export function RoleEditScreen(props: Props) {
  const {
    form,
    handlers: { handleSave, handleDelete },
  } = useRoleEditScreen(props);

  return (
    <styled.form
      className={lstack()}
      h="full"
      justifyContent="space-between"
      onSubmit={handleSave}
    >
      <LStack px="0.5" maxH="full" overflowY="scroll">
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
            items={PermissionList.map((p) => ({
              value: p.value,
              label: p.name,
              description: p.description,
            }))}
          />
          <FormErrorText>{form.formState.errors.name?.message}</FormErrorText>
        </FormControl>
      </LStack>

      <HStack w="full">
        <DeleteWithConfirmationButton
          type="button"
          w="full"
          onDelete={handleDelete}
        />
        <Button w="full" variant="outline" disabled={!form.formState.isDirty}>
          Save
        </Button>
      </HStack>
    </styled.form>
  );
}
