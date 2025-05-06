import { uniqueId } from "lodash/fp";
import { ChangeEvent } from "react";
import { Controller } from "react-hook-form";

import { PropertyName, PropertyType } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { Input } from "@/components/ui/input";
import { Center, HStack, LStack, styled } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { Form } from "../../form";
import { useEditState } from "../../useEditState";

export function LibraryPagePropertiesBlock() {
  const { editing } = useEditState();
  const { node } = useLibraryPageContext();

  if (editing) {
    return <LibraryPagePropertiesBlockEditable />;
  }

  return (
    <styled.dl display="table" borderCollapse="collapse">
      {node.properties.map((p) => {
        return (
          <HStack key={p.name} display="table-row">
            <styled.dt
              display="table-cell"
              w="32"
              p="1"
              borderRadius="sm"
              textOverflow="ellipsis"
              overflowX="hidden"
              color="fg.muted"
              _hover={{
                color: "fg.default",
                background: "bg.muted",
                cursor: "pointer",
              }}
            >
              {p.name}
            </styled.dt>
            <styled.dd
              display="table-cell"
              p="1"
              w="min"
              borderRadius="sm"
              _hover={{
                color: "fg.default",
                background: "bg.muted",
                cursor: "pointer",
              }}
            >
              {p.value}
            </styled.dd>
          </HStack>
        );
      })}
    </styled.dl>
  );
}

function LibraryPagePropertiesBlockEditable() {
  const { node, form } = useLibraryPageContext();
  return (
    <Controller<Form>
      control={form.control}
      name="properties"
      render={({ field, fieldState, formState }) => {
        const fieldValue = field.value as Form["properties"];
        const initialValue = formState.defaultValues?.[
          "properties"
        ] as Form["properties"];
        const current = fieldValue ?? initialValue ?? [];

        async function handleAddProperty() {
          const existingNames = new Set(current.map((f) => f.name));
          let newName = "Field 1";
          let counter = 1;
          while (existingNames.has(newName)) {
            newName = `Field ${counter++}`;
          }

          field.onChange([
            ...current,
            {
              fid: uniqueId("new_field_"),
              name: newName,
              type: PropertyType.text, // TODO: Add actual UI for types
              sort: "5",
              value: "",
            },
          ]);
        }

        async function handleRemoveProperty(name: PropertyName) {
          const next = current.filter((f) => f.name !== name);
          field.onChange(next);
        }

        function handlePropertyNameChange(name: PropertyName, newName: string) {
          const next = current.map((f) => {
            if (f.name === name) {
              f.name = newName;
            }

            return f;
          });
          field.onChange(next);
        }

        function handlePropertyValueChange(name: PropertyName, value: string) {
          const next = current.map((f) => {
            if (f.name === name) {
              f.value = value;
            }

            return f;
          });
          field.onChange(next);
        }

        return (
          <LStack w="64">
            <styled.dl display="table" borderCollapse="collapse">
              {current.map((p) => {
                function handleRemove() {
                  handleRemoveProperty(p.name);
                }

                function handleNameChange(e: ChangeEvent<HTMLInputElement>) {
                  handlePropertyNameChange(p.name, e.target.value);
                }

                function handleValueChange(e: ChangeEvent<HTMLInputElement>) {
                  handlePropertyValueChange(p.name, e.target.value);
                }

                return (
                  <HStack key={p.fid} display="table-row">
                    <styled.dt display="table-cell" p="1" color="fg.muted">
                      <Input
                        variant="ghost"
                        defaultValue={p.name}
                        onChange={handleNameChange}
                      />
                    </styled.dt>
                    <styled.dd display="table-cell" p="1">
                      <Input
                        variant="ghost"
                        defaultValue={p.value}
                        onChange={handleValueChange}
                      />
                    </styled.dd>

                    <Center>
                      <IconButton
                        type="button"
                        variant="ghost"
                        color="fg.destructive"
                        size="sm"
                        onClick={handleRemove}
                      >
                        <DeleteIcon />
                      </IconButton>
                    </Center>
                  </HStack>
                );
              })}
            </styled.dl>
            <Button
              type="button"
              w="full"
              size="xs"
              variant="subtle"
              onClick={handleAddProperty}
            >
              Add Property
            </Button>
          </LStack>
        );
      }}
    />
  );
}
