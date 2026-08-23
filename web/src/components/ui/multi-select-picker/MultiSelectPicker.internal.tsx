"use client";

import { createListCollection } from "@ark-ui/react";
import {
  FloatingPortal,
  autoUpdate,
  flip,
  size as floatingSize,
  offset,
  shift,
  useFloating,
} from "@floating-ui/react";
import { type CSSProperties, type Ref, useRef, useState } from "react";
import { cx } from "styled-system/css";

import { Unready } from "@/components/site/Unready";
import * as Combobox from "@/components/ui/combobox";
import { styled } from "@/styled-system/jsx";
import {
  type InputVariantProps,
  input,
  multiSelectPicker,
} from "@/styled-system/recipes";
import { deriveColour } from "@/utils/colour";

import { Badge, badgeColourPalette, badgeColours } from "../badge";
import { IconButton } from "../icon-button";
import { AddIcon } from "../icons/Add";
import { CheckIcon } from "../icons/Check";
import { SelectIcon } from "../icons/Select";
import { Text } from "../text";

type ControlSize = "sm" | "md" | "lg";

export type MultiSelectPickerItem = {
  label: string;
  value: string;
  colour?: string;
};

type RootProps = Omit<
  Combobox.RootProps,
  | "children"
  | "collection"
  | "defaultInputValue"
  | "defaultValue"
  | "inputValue"
  | "multiple"
  | "onChange"
  | "onInputValueChange"
  | "onValueChange"
  | "open"
  | "value"
>;

export type MultiSelectPickerProps = RootProps & {
  value: MultiSelectPickerItem[];
  onChange: (items: MultiSelectPickerItem[]) => Promise<void>;
  onQuery: (query: string) => void;
  queryResults?: MultiSelectPickerItem[];
  allowNewValues?: boolean;
  inputPlaceholder?: string;
  ariaLabel?: string;
  autoColour?: boolean;
  queryError?: string | null;
  size?: ControlSize;
  inputVariant?: InputVariantProps["variant"];
};

type PickerCollectionItem = MultiSelectPickerItem & {
  kind?: "create";
};

const CREATE_VALUE_PREFIX = "multi-select-picker:create:";

function selectedBadgeSize(size: ControlSize | undefined) {
  if (size === "lg") return "lg";
  if (size === "md") return "md";
  return "sm";
}

function itemColourStyles(item: MultiSelectPickerItem) {
  const colour = badgeColours(
    item.colour ? item.colour : deriveColour(item.value),
  );

  return colour ? badgeColourPalette(colour) : undefined;
}

function uniqueItems(items: PickerCollectionItem[]) {
  const seen = new Set<string>();

  return items.filter((item) => {
    if (seen.has(item.value)) return false;
    seen.add(item.value);
    return true;
  });
}

export function MultiSelectPicker({
  value,
  onChange,
  onQuery,
  queryResults,
  allowNewValues = false,
  inputPlaceholder = "Select items...",
  ariaLabel,
  autoColour = false,
  queryError,
  size = "sm",
  inputVariant = "ghost",
  positioning,
  onOpenChange,
  defaultOpen = false,
  ...rootProps
}: MultiSelectPickerProps) {
  const [queryInput, setQueryInput] = useState("");
  const [isOpen, setIsOpen] = useState(defaultOpen);
  const inputRef = useRef<HTMLInputElement>(null);
  const {
    refs: { setFloating, setReference },
    floatingStyles,
  } = useFloating({
    open: isOpen,
    placement: positioning?.placement ?? "bottom-start",
    strategy: "fixed",
    whileElementsMounted: autoUpdate,
    middleware: [
      offset(positioning?.gutter ?? 4),
      flip({ padding: positioning?.overflowPadding ?? 8 }),
      shift({ padding: positioning?.overflowPadding ?? 8 }),
      floatingSize({
        padding: positioning?.overflowPadding ?? 8,
        apply({ availableHeight, availableWidth, elements, rects }) {
          const width = Math.min(
            rects.reference.width,
            512,
            Math.max(0, availableWidth),
          );

          Object.assign(elements.floating.style, {
            boxSizing: "border-box",
            maxHeight: `${Math.min(384, Math.max(0, availableHeight))}px`,
            maxWidth: `${Math.max(0, availableWidth)}px`,
            width: `${width}px`,
          });
        },
      }),
    ],
  });
  const trimmedQuery = queryInput.trim();
  const selectedValues = new Set(value.map((item) => item.value));
  const filteredQueryResults =
    queryResults?.filter((result) => !selectedValues.has(result.value)) ?? [];
  const showCreateNew =
    allowNewValues &&
    trimmedQuery.length > 0 &&
    !selectedValues.has(trimmedQuery) &&
    !filteredQueryResults.some((result) => result.value === trimmedQuery);
  const createItem: PickerCollectionItem | undefined = showCreateNew
    ? {
        kind: "create",
        label: `Create “${trimmedQuery}”`,
        value: `${CREATE_VALUE_PREFIX}${encodeURIComponent(trimmedQuery)}`,
      }
    : undefined;
  const collectionItems = uniqueItems([
    ...value,
    ...filteredQueryResults,
    ...(createItem ? [createItem] : []),
  ]);
  const collection = createListCollection({ items: collectionItems });
  const itemsByValue = new Map(
    collectionItems.map((item) => [item.value, item]),
  );
  const badgeSize = selectedBadgeSize(size);
  const inputLabel = ariaLabel ?? "Search for items";
  const triggerLabel = ariaLabel ?? "Select items";

  function resetQuery() {
    setQueryInput("");
    onQuery("");
  }

  function addColour(item: MultiSelectPickerItem) {
    return autoColour && !item.colour
      ? { ...item, colour: deriveColour(item.value) }
      : item;
  }

  async function createNewValue() {
    if (!showCreateNew) return;

    await onChange([
      ...value,
      {
        label: trimmedQuery,
        value: trimmedQuery,
        colour: autoColour ? deriveColour(trimmedQuery) : undefined,
      },
    ]);
    resetQuery();
  }

  function handleValueChange(nextValues: string[]) {
    if (createItem && nextValues.includes(createItem.value)) {
      void createNewValue();
      return;
    }

    const nextItems = nextValues.flatMap((itemValue) => {
      const item = itemsByValue.get(itemValue);
      if (!item || item.kind === "create") return [];

      return [selectedValues.has(itemValue) ? item : addColour(item)];
    });

    void onChange(nextItems);
    resetQuery();
  }

  return (
    <Combobox.Root
      {...rootProps}
      collection={collection}
      closeOnSelect={false}
      inputBehavior="autohighlight"
      inputValue={queryInput}
      multiple
      open={isOpen}
      openOnClick
      positioning={{ placement: "bottom-start", ...positioning }}
      size={size}
      value={value.map((item) => item.value)}
      onOpenChange={(details) => {
        setIsOpen(details.open);
        if (!details.open && queryInput) resetQuery();
        onOpenChange?.(details);
      }}
      onValueChange={({ value: nextValues }) => handleValueChange(nextValues)}
    >
      <Combobox.Label srOnly>{inputLabel}</Combobox.Label>
      <Combobox.Control
        ref={setReference}
        className={multiSelectPicker({ size })}
        onClick={(event) => {
          if (!(event.target as HTMLElement).closest("button")) {
            inputRef.current?.focus();
          }
        }}
      >
        <styled.div
          alignItems="center"
          display="flex"
          flex="1"
          flexWrap="wrap"
          gap="1"
          minW="0"
        >
          {value.map((item) => (
            <Badge
              key={item.value}
              aria-hidden="true"
              flexShrink="0"
              size={badgeSize}
              style={itemColourStyles(item)}
            >
              {item.label}
            </Badge>
          ))}

          <Combobox.Input
            ref={inputRef}
            aria-label={inputLabel}
            autoComplete="off"
            className={cx(input({ size, variant: inputVariant }))}
            flex="1"
            h="auto!"
            minW={value.length > 0 ? "24" : "32"}
            placeholder={inputPlaceholder}
            pr="0!"
            width="auto"
            onChange={(event) => {
              const nextQuery = event.currentTarget.value;
              setQueryInput(nextQuery);
              onQuery(nextQuery);
            }}
            onKeyDown={(event) => {
              if (
                event.key === "Backspace" &&
                event.currentTarget.value === "" &&
                value.length > 0
              ) {
                event.preventDefault();
                void onChange(value.slice(0, -1));
                return;
              }

              if (event.key === "Enter" && showCreateNew) {
                event.preventDefault();
                void createNewValue();
              }
            }}
          />
        </styled.div>

        <Combobox.Trigger asChild>
          <IconButton
            aria-label={triggerLabel}
            position="absolute"
            right="0"
            size={size}
            top="[50%]"
            transform="translateY(-50%)"
            type="button"
            variant="ghost"
          >
            <SelectIcon />
          </IconButton>
        </Combobox.Trigger>
      </Combobox.Control>

      {isOpen && (
        <PickerContent
          createItem={createItem}
          floatingRef={setFloating}
          floatingStyles={floatingStyles}
          queryError={queryError}
          queryInput={queryInput}
          queryResults={filteredQueryResults}
          selectedItems={value}
        />
      )}
    </Combobox.Root>
  );
}

type PickerContentProps = {
  createItem?: PickerCollectionItem;
  floatingRef: Ref<HTMLDivElement>;
  floatingStyles: CSSProperties;
  queryError?: string | null;
  queryInput: string;
  queryResults: MultiSelectPickerItem[];
  selectedItems: MultiSelectPickerItem[];
};

function PickerContent({
  createItem,
  floatingRef,
  floatingStyles,
  queryError,
  queryInput,
  queryResults,
  selectedItems,
}: PickerContentProps) {
  const hasOptions = queryResults.length > 0 || createItem !== undefined;

  return (
    <FloatingPortal>
      <Combobox.Content
        ref={floatingRef}
        style={floatingStyles}
        zIndex="popover"
      >
        <Combobox.List>
          {selectedItems.length > 0 && (
            <Combobox.ItemGroup>
              <PickerGroupLabel>Selected</PickerGroupLabel>
              {selectedItems.map((item) => (
                <PickerOption key={item.value} item={item} />
              ))}
            </Combobox.ItemGroup>
          )}

          {!queryError && queryResults.length > 0 && (
            <Combobox.ItemGroup>
              <PickerGroupLabel>
                {queryInput ? "Results" : "Available"}
              </PickerGroupLabel>
              {queryResults.map((item) => (
                <PickerOption key={item.value} item={item} />
              ))}
            </Combobox.ItemGroup>
          )}

          {!queryError && createItem && (
            <Combobox.ItemGroup>
              <PickerGroupLabel>Create new</PickerGroupLabel>
              <Combobox.Item item={createItem}>
                <Combobox.ItemText>{createItem.label}</Combobox.ItemText>
                <AddIcon aria-hidden="true" />
              </Combobox.Item>
            </Combobox.ItemGroup>
          )}
        </Combobox.List>

        {queryError ? (
          <styled.div p="2">
            <Unready error={queryError} />
          </styled.div>
        ) : (
          !hasOptions && (
            <styled.div px="2" py="3" role="status">
              <Text variant="supporting" color="text.muted">
                {queryInput
                  ? `No results for “${queryInput}”`
                  : "No options available"}
              </Text>
            </styled.div>
          )
        )}
      </Combobox.Content>
    </FloatingPortal>
  );
}

function PickerGroupLabel({ children }: { children: React.ReactNode }) {
  return (
    <Combobox.ItemGroupLabel asChild>
      <Text as="div" variant="metadata" fontWeight="medium">
        {children}
      </Text>
    </Combobox.ItemGroupLabel>
  );
}

function PickerOption({ item }: { item: MultiSelectPickerItem }) {
  return (
    <Combobox.Item item={item}>
      <Combobox.ItemText>
        <styled.span lineClamp="1" title={item.label}>
          {item.label}
        </styled.span>
      </Combobox.ItemText>
      <Combobox.ItemIndicator>
        <CheckIcon />
      </Combobox.ItemIndicator>
    </Combobox.Item>
  );
}
