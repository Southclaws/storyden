import { useFilter, useListCollection } from "@ark-ui/react";

import * as Combobox from "@/components/ui/combobox";
import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { ChevronUpDownIcon } from "@/components/ui/icons/Chevron";
import { Input } from "@/components/ui/input";
import { NumberInput } from "@/components/ui/number-input";
import { CardBox, HStack, LStack, styled } from "@/styled-system/jsx";

import { formatSeconds } from "./useSystemSettings";

type OperationCostOverrideAddFormProps = {
  availableOperations: string[];
  selectedOperation: string;
  costValue: number;
  rateLimit: number;
  rateLimitPeriod: number;
  onOperationChange: (operation: string) => void;
  onCostChange: (cost: number) => void;
  onAdd: () => void;
};

export function OperationCostOverrideAddForm({
  rateLimit,
  rateLimitPeriod,
  availableOperations,
  selectedOperation,
  costValue,
  onOperationChange,
  onCostChange,
  onAdd,
}: OperationCostOverrideAddFormProps) {
  const { contains } = useFilter({ sensitivity: "base" });

  const items = availableOperations.map((op) => ({ label: op, value: op }));

  const { collection, filter } = useListCollection({
    initialItems: items,
    itemToString: (item) => item.label,
    itemToValue: (item) => item.value,
    filter: contains,
  });

  const handleInputChange = (details: Combobox.InputValueChangeDetails) => {
    filter(details.inputValue);
  };

  const effectiveLimit = Math.floor(rateLimit / costValue);

  return (
    <CardBox borderRadius="sm">
      <LStack w="full" gap="1">
        <HStack w="full" gap="2">
          <Combobox.Root
            collection={collection}
            value={selectedOperation ? [selectedOperation] : []}
            onInputValueChange={handleInputChange}
            onValueChange={(details) =>
              onOperationChange(details.value[0] || "")
            }
            flex="1"
          >
            <Combobox.Control>
              <Combobox.Input placeholder="Search operations..." asChild>
                <Input size="sm" />
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
                  {collection.items.map((item) => (
                    <Combobox.Item key={item.value} item={item}>
                      <Combobox.ItemText>{item.label}</Combobox.ItemText>
                    </Combobox.Item>
                  ))}
                </Combobox.ItemGroup>
              </Combobox.Content>
            </Combobox.Positioner>
          </Combobox.Root>

          <NumberInput
            width="16"
            size="sm"
            min={1}
            max={100}
            value={costValue.toString()}
            onValueChange={(details) =>
              onCostChange(parseInt(details.value, 10) || 1)
            }
          />

          <IconButton type="button" variant="subtle" size="sm" onClick={onAdd}>
            <AddIcon />
          </IconButton>
        </HStack>

        <styled.p fontSize="xs" color="fg.muted">
          Can be performed{" "}
          <styled.strong color="fg.info">{effectiveLimit}</styled.strong> times
          every{" "}
          <styled.strong color="fg.info">
            {formatSeconds(rateLimitPeriod)}
          </styled.strong>
        </styled.p>
      </LStack>
    </CardBox>
  );
}
