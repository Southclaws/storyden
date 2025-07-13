"use client";

import {
  ComboboxValueChangeDetails,
  createListCollection,
} from "@ark-ui/react";
import { map } from "lodash/fp";
import { useState } from "react";

import { useNodeList } from "@/api/openapi-client/nodes";
import { Node } from "@/api/openapi-schema";
import { Unready } from "@/components/site/Unready";
import * as Combobox from "@/components/ui/combobox";
import { IconButton } from "@/components/ui/icon-button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { ChevronUpDownIcon } from "@/components/ui/icons/Chevron";
import { Input } from "@/components/ui/input";

type Props = Omit<
  Combobox.RootProps,
  "onChange" | "defaultValue" | "value" | "collection"
> & {
  defaultValue?: string;
  value?: string;
  onChange: (node: Node | undefined) => void;
};

export function LibraryPageSelect(props: Props) {
  const { data, error } = useNodeList({
    visibility: ["published"],
  });
  if (!data) {
    return <Unready error={error} />;
  }

  return <LibraryPageSelectCombobox nodes={data.nodes} {...props} />;
}

type Item = {
  label: string;
  value: string;
};

function toItem(node: Node): Item {
  return {
    label: node.name,
    value: node.id,
  };
}

const toItems = map(toItem);

function LibraryPageSelectCombobox({
  defaultValue,
  value,
  onChange,
  nodes,
  ...rest
}: Props & { nodes: Node[] }) {
  const initialCollection = createListCollection({
    items: toItems(nodes),
  });

  const [collection, setCollection] = useState(initialCollection);

  const handleInputChange = ({
    inputValue,
  }: Combobox.InputValueChangeDetails) => {
    const filtered = initialCollection.items.filter((item) =>
      item.label.toLowerCase().includes(inputValue.toLowerCase()),
    );

    setCollection(
      filtered.length > 0
        ? createListCollection({ items: filtered })
        : initialCollection,
    );
  };

  const handleOpenChange = () => {
    setCollection(initialCollection);
  };

  function handleChange({ value }: ComboboxValueChangeDetails) {
    if (!value || value.length === 0) {
      return;
    }

    const selectedNode = nodes.find((item) => item.id === value[0]);

    onChange(selectedNode);
  }

  return (
    <Combobox.Root
      {...rest}
      collection={collection}
      defaultValue={defaultValue ? [defaultValue] : undefined}
      onInputValueChange={handleInputChange}
      onOpenChange={handleOpenChange}
      onValueChange={handleChange}
      value={value ? [value] : []}
      size="xs"
    >
      <Combobox.Control>
        <Combobox.Input placeholder="Select a Page" asChild>
          <Input size="xs" />
        </Combobox.Input>
        <Combobox.Trigger asChild>
          <IconButton variant="link" aria-label="open" size="xs">
            <ChevronUpDownIcon />
          </IconButton>
        </Combobox.Trigger>
      </Combobox.Control>
      <Combobox.Positioner>
        <Combobox.Content>
          <Combobox.ItemGroup>
            <Combobox.Item key="unset" item="unset">
              <Combobox.ItemText>All root level pages</Combobox.ItemText>
              <Combobox.ItemIndicator>
                <CheckIcon />
              </Combobox.ItemIndicator>
            </Combobox.Item>
          </Combobox.ItemGroup>
          <Combobox.ItemGroup>
            {collection.items.map((item) => (
              <Combobox.Item key={item.value} item={item}>
                <Combobox.ItemText>{item.label}</Combobox.ItemText>
                <Combobox.ItemIndicator>
                  <CheckIcon />
                </Combobox.ItemIndicator>
              </Combobox.Item>
            ))}
          </Combobox.ItemGroup>
        </Combobox.Content>
      </Combobox.Positioner>
    </Combobox.Root>
  );
}
