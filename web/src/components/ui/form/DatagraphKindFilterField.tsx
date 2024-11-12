import { Portal, ToggleGroupValueChangeDetails } from "@ark-ui/react";
import { JSX } from "react";
import { Controller, ControllerProps, FieldValues } from "react-hook-form";

import * as ToggleGroup from "@/components/ui/toggle-group";
import * as Tooltip from "@/components/ui/tooltip";
import { HStack } from "@/styled-system/jsx";

type CollectionItem = {
  label: string;
  icon: JSX.Element;
  description: string;
  value: string;
};

type Props<T extends FieldValues> = Omit<ControllerProps<T>, "render"> & {
  items: CollectionItem[];
};

export function DatagraphKindFilterField<T extends FieldValues>({
  items,
  ...props
}: Props<T>) {
  return (
    <Controller<T>
      {...props}
      render={({ formState, field }) => {
        function handleChangeFilter({ value }: ToggleGroupValueChangeDetails) {
          field.onChange(value);
        }

        return (
          <ToggleGroup.Root
            multiple
            size="xs"
            onValueChange={handleChangeFilter}
            defaultValue={formState.defaultValues?.[props.name]}
          >
            {items.map((item) => (
              <ToggleGroup.Item
                key={item.value}
                value={item.value}
                aria-label={item.description}
              >
                <Tooltip.Root
                  lazyMount
                  openDelay={0}
                  positioning={{
                    slide: true,
                    shift: -48,
                    placement: "right-end",
                  }}
                >
                  <Tooltip.Trigger asChild>
                    <HStack gap="1">
                      {item.icon} {item.label}
                    </HStack>
                  </Tooltip.Trigger>

                  <Portal>
                    <Tooltip.Positioner>
                      <Tooltip.Arrow>
                        <Tooltip.ArrowTip />
                      </Tooltip.Arrow>

                      <Tooltip.Content>{item.description}</Tooltip.Content>
                    </Tooltip.Positioner>
                  </Portal>
                </Tooltip.Root>
              </ToggleGroup.Item>
            ))}
          </ToggleGroup.Root>
        );
      }}
    />
  );
}
