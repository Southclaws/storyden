import { Controller, useWatch } from "react-hook-form";

import { DeleteWithConfirmationButton } from "@/components/site/DeleteConfirmationButton";
import { InfoTip } from "@/components/site/InfoTip";
import { ColourPickerField } from "@/components/ui/ColourPickerField";
import { FormControl } from "@/components/ui/FormControl";
import { FormErrorText } from "@/components/ui/FormErrorText";
import { FormLabel } from "@/components/ui/FormLabel";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { CardGroupSelect } from "@/components/ui/form/CardGroupSelect";
import { Input } from "@/components/ui/input";
import {
  PermissionList,
  buildPermissionList,
} from "@/lib/permission/permission";
import {
  isDefaultRole,
  isGuestRole,
  isStoredDefaultRole,
  readPermissions,
} from "@/lib/role/defaults";
import { DefaultRoleMetadata } from "@/lib/role/metadata";
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
  const roleName = useWatch({
    control: form.control,
    name: "name",
  });
  const roleColour = useWatch({
    control: form.control,
    name: "colour",
  });
  const roleBold = useWatch({
    control: form.control,
    name: "meta.bold",
  });
  const roleItalic = useWatch({
    control: form.control,
    name: "meta.italic",
  });
  const roleColoured = useWatch({
    control: form.control,
    name: "meta.coloured",
  });
  const roleMeta = {
    bold: roleBold ?? DefaultRoleMetadata.bold,
    italic: roleItalic ?? DefaultRoleMetadata.italic,
    coloured: roleColoured ?? DefaultRoleMetadata.coloured,
  };

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
          <FormLabel>Name Decorations</FormLabel>

          <LStack gap="1">
            <Controller
              control={form.control}
              name="meta.bold"
              render={({ field }) => (
                <Checkbox
                  size="sm"
                  checked={!!field.value}
                  onCheckedChange={({ checked }) => {
                    field.onChange(checked === true);
                  }}
                >
                  Bold name
                </Checkbox>
              )}
            />

            <Controller
              control={form.control}
              name="meta.italic"
              render={({ field }) => (
                <Checkbox
                  size="sm"
                  checked={!!field.value}
                  onCheckedChange={({ checked }) => {
                    field.onChange(checked === true);
                  }}
                >
                  Italic name
                </Checkbox>
              )}
            />

            <Controller
              control={form.control}
              name="meta.coloured"
              render={({ field }) => (
                <Checkbox
                  size="sm"
                  checked={!!field.value}
                  onCheckedChange={({ checked }) => {
                    field.onChange(checked === true);
                  }}
                >
                  Coloured name{" "}
                  <InfoTip title="Name colour">
                    Roles are ordered by priority, and the highest priority role
                    with coloured name enabled will determine the colour of the
                    member's name.
                  </InfoTip>
                </Checkbox>
              )}
            />
          </LStack>
        </FormControl>

        <FormControl>
          <FormLabel>Name Preview</FormLabel>
          <LStack gap="1.5">
            <RoleNamePreview
              mode="light"
              roleName={roleName}
              roleColour={roleColour}
              roleMeta={roleMeta}
            />

            <RoleNamePreview
              mode="dark"
              roleName={roleName}
              roleColour={roleColour}
              roleMeta={roleMeta}
            />
          </LStack>
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

type RoleNamePreviewProps = {
  mode: "light" | "dark";
  roleName: string;
  roleColour?: string;
  roleMeta: {
    bold: boolean;
    italic: boolean;
    coloured: boolean;
  };
};

function RoleNamePreview({
  mode,
  roleName,
  roleColour,
  roleMeta,
}: RoleNamePreviewProps) {
  const isDark = mode === "dark";
  const backgroundColor = isDark ? "slate.dark.2" : "neutral.light.2";
  const borderColor = isDark ? "neutral.dark.6" : "neutral.light.5";
  const labelColour = isDark ? "neutral.dark.11" : "neutral.light.11";
  const handleColour = isDark ? "neutral.dark.10" : "neutral.light.10";
  const defaultNameColour = isDark ? "neutral.dark.12" : "neutral.light.12";

  return (
    <WStack
      p="2"
      borderWidth="thin"
      borderStyle="solid"
      borderColor={borderColor}
      backgroundColor={backgroundColor}
      borderRadius="md"
    >
      <styled.p fontSize="xs" color={labelColour} textTransform="uppercase">
        {mode}
      </styled.p>

      <styled.p lineHeight="tight" fontSize="md">
        <styled.span
          color={
            !roleColour || !roleMeta.coloured ? defaultNameColour : undefined
          }
          style={
            roleColour && roleMeta.coloured ? { color: roleColour } : undefined
          }
          fontWeight={roleMeta.bold ? "bold" : "medium"}
          fontStyle={roleMeta.italic ? "italic" : "normal"}
        >
          {roleName}
        </styled.span>{" "}
        <styled.span
          color={handleColour}
          fontWeight={roleMeta.bold ? "semibold" : "medium"}
        >
          @sample
        </styled.span>
      </styled.p>
    </WStack>
  );
}
