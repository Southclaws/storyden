import { createListCollection } from "@ark-ui/react";
import { useEffect, useMemo, useState } from "react";
import { Control, Controller, FieldPath, FieldValues } from "react-hook-form";

import { RobotModelInfo } from "@/api/openapi-schema";
import * as Combobox from "@/components/ui/combobox";
import { IconButton } from "@/components/ui/icon-button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { ChevronUpDownIcon } from "@/components/ui/icons/Chevron";
import { Input } from "@/components/ui/input";
import { styled } from "@/styled-system/jsx";

type Item = {
  label: string;
  value: string;
};

type Props<T extends FieldValues> = {
  control: Control<T>;
  name: FieldPath<T>;
  models: RobotModelInfo[];
  placeholder: string;
  disabled?: boolean;
};

export function RobotModelComboboxField<T extends FieldValues>({
  control,
  name,
  models,
  placeholder,
  disabled,
}: Props<T>) {
  const initialCollection = useMemo(
    () =>
      createListCollection({
        items: models.map(modelToItem),
      }),
    [models],
  );
  const [collection, setCollection] = useState(initialCollection);

  useEffect(() => {
    setCollection(initialCollection);
  }, [initialCollection]);

  function handleInputChange({ inputValue }: Combobox.InputValueChangeDetails) {
    const query = inputValue.toLowerCase();
    const filtered = initialCollection.items.filter(
      (item) =>
        item.label.toLowerCase().includes(query) ||
        item.value.toLowerCase().includes(query),
    );

    setCollection(createListCollection({ items: filtered }));
  }

  return (
    <Controller
      control={control}
      name={name}
      render={({ field }) => {
        function handleChange({ value }: Combobox.ValueChangeDetails) {
          field.onChange(value[0] ?? undefined);
        }

        return (
          <Combobox.Root
            collection={collection}
            value={field.value ? [field.value] : []}
            onValueChange={handleChange}
            onInputValueChange={handleInputChange}
            onOpenChange={() => setCollection(initialCollection)}
            positioning={{ sameWidth: true, fitViewport: true }}
            disabled={disabled}
            size="sm"
          >
            <Combobox.Control>
              <Combobox.Input placeholder={placeholder} asChild>
                <Input size="sm" />
              </Combobox.Input>
              <Combobox.Trigger asChild>
                <IconButton
                  type="button"
                  variant="link"
                  aria-label="Select model"
                  size="sm"
                >
                  <ChevronUpDownIcon />
                </IconButton>
              </Combobox.Trigger>
            </Combobox.Control>

            <Combobox.Positioner>
              <Combobox.Content>
                <Combobox.List>
                  <Combobox.ItemGroup>
                    {collection.items.map((item) => (
                      <Combobox.Item key={item.value} item={item}>
                        <Combobox.ItemText>
                          <styled.span lineClamp="1">{item.label}</styled.span>
                        </Combobox.ItemText>
                        <Combobox.ItemIndicator>
                          <CheckIcon />
                        </Combobox.ItemIndicator>
                      </Combobox.Item>
                    ))}
                  </Combobox.ItemGroup>
                </Combobox.List>
              </Combobox.Content>
            </Combobox.Positioner>
          </Combobox.Root>
        );
      }}
    />
  );
}

function modelToItem(model: RobotModelInfo): Item {
  return {
    label: model.ref,
    value: model.ref,
  };
}
