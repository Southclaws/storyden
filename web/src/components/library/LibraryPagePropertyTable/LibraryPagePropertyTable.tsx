import { uniqueId } from "lodash/fp";
import { useRef, useState } from "react";
import {
  Controller,
  ControllerFieldState,
  ControllerProps,
  ControllerRenderProps,
  FieldValues,
  UseFormStateReturn,
} from "react-hook-form";

import { handle } from "@/api/client";
import { tagList } from "@/api/openapi-client/tags";
import {
  Identifier,
  InstanceCapability,
  Node,
  NodeWithChildren,
  Property,
  PropertyName,
  PropertySchema,
  TagNameList,
} from "@/api/openapi-schema";
import { IntelligenceAction } from "@/components/site/Action/Intelligence";
import { TagBadgeList } from "@/components/tag/TagBadgeList";
import { Button } from "@/components/ui/button";
import { Combotags, CombotagsHandle } from "@/components/ui/combotags";
import { IconButton } from "@/components/ui/icon-button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { Input } from "@/components/ui/input";
import { useLibraryMutation } from "@/lib/library/library";
import { usePropertyMutation } from "@/lib/library/properties";
import { useCapability } from "@/lib/settings/capabilities";
import { Form } from "@/screens/library/LibraryPageScreen/useLibraryPageScreen";
import { HStack, LStack, styled } from "@/styled-system/jsx";

export type Props<T extends FieldValues> = Omit<
  ControllerProps<T>,
  "render"
> & {
  editing: boolean;
  node: NodeWithChildren;
};

export function LibraryPagePropertyTable<T extends FieldValues>({
  editing,
  node,
  ...props
}: Props<T>) {
  if (editing) {
    return (
      <LibraryPagePropertyTableEditable
        {...props}
        editing={editing}
        node={node}
      />
    );
  }

  console.log(node.properties);

  return (
    <styled.dl display="table" borderCollapse="collapse">
      {node.properties?.map((p) => {
        if (p.value == null) {
          return null;
        }

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

export function LibraryPagePropertyTableEditable({
  editing,
  node,
  ...props
}: Props<Form>) {
  return (
    <Controller<Form>
      control={props.control}
      name="properties"
      render={({ field, fieldState, formState }) => {
        const fieldValue = field.value as Form["properties"];
        const initialValue = formState.defaultValues?.[
          "properties"
        ] as Form["properties"];
        const current = fieldValue ?? initialValue ?? [];

        async function handleAddProperty() {
          // await addField({ name: "", type: "string", sort: "99999999" });

          console.log(field, fieldState, formState);

          field.onChange([
            ...current,
            {
              name: uniqueId("Field"),
              type: "text",
              value: "",
            },
          ]);
        }

        async function handleRemoveProperty(name: PropertyName) {
          const next = current.filter((f) => f.name !== name);
          field.onChange(next);
        }

        console.log({
          defaults: formState.defaultValues,
          fieldValue,
          initialValue,
          current,
        });

        return (
          <LStack w="64">
            <pre>properties: {JSON.stringify(current)}</pre>
            <styled.dl display="table" borderCollapse="collapse">
              {current.map((p) => {
                function handleRemove() {
                  handleRemoveProperty(p.name);
                }
                return (
                  <HStack key={p.name} display="table-row">
                    <styled.dt
                      display="table-cell"
                      w="full"
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
                      <Input defaultValue={p.name} />
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

                    <IconButton
                      type="button"
                      variant="subtle"
                      onClick={handleRemove}
                    >
                      <DeleteIcon />
                    </IconButton>
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
