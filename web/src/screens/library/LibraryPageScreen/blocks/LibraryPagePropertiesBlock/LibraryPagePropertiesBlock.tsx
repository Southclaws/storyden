import { ChangeEvent } from "react";

import { PropertyName, PropertyType } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { IconButton } from "@/components/ui/icon-button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { Input } from "@/components/ui/input";
import { Center, HStack, LStack, styled } from "@/styled-system/jsx";

import { useLibraryPageContext } from "../../Context";
import { useWatch } from "../../store";
import { useEditState } from "../../useEditState";

export function LibraryPagePropertiesBlock() {
  const { editing } = useEditState();
  const properties = useWatch((s) => s.draft.properties);

  if (editing) {
    return <LibraryPagePropertiesBlockEditable />;
  }

  return (
    <styled.dl display="table" borderCollapse="collapse">
      {properties.map((p) => {
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
  const { store } = useLibraryPageContext();
  const {
    addProperty,
    removePropertyByName,
    setPropertyName,
    setPropertyValue,
  } = store.getState();
  const current = useWatch((s) => s.draft.properties);

  function handleAddProperty() {
    addProperty("Field", PropertyType.text);
  }

  const handleRemoveProperty = (name: PropertyName) => () => {
    removePropertyByName(name);
  };

  const handlePropertyNameChange =
    (name: PropertyName) => (e: ChangeEvent<HTMLInputElement>) => {
      setPropertyName(name, e.target.value);
    };

  const handlePropertyValueChange =
    (name: PropertyName) => (e: ChangeEvent<HTMLInputElement>) => {
      setPropertyValue(name, e.target.value);
    };

  return (
    <LStack w="64">
      {current.length > 0 && (
        <styled.dl display="table" borderCollapse="collapse">
          {current.map((p) => {
            return (
              <HStack key={p.fid} display="table-row">
                <styled.dt display="table-cell" p="1" color="fg.muted">
                  <Input
                    variant="ghost"
                    defaultValue={p.name}
                    onChange={handlePropertyNameChange(p.name)}
                  />
                </styled.dt>
                <styled.dd display="table-cell" p="1">
                  <Input
                    variant="ghost"
                    defaultValue={p.value}
                    onChange={handlePropertyValueChange(p.name)}
                  />
                </styled.dd>

                <Center>
                  <IconButton
                    type="button"
                    variant="ghost"
                    color="fg.destructive"
                    size="sm"
                    onClick={handleRemoveProperty(p.name)}
                  >
                    <DeleteIcon />
                  </IconButton>
                </Center>
              </HStack>
            );
          })}
        </styled.dl>
      )}
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
}
