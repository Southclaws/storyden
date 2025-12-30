import { createListCollection, useTagsInput } from "@ark-ui/react";
import { uniq } from "lodash";
import {
  forwardRef,
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
} from "react";

import * as Combobox from "@/components/ui/combobox";
import { IconButton } from "@/components/ui/icon-button";
import * as TagsInput from "@/components/ui/tags-input";

import { DeleteSmallIcon } from "./icons/Delete";

export type CombotagsItem = {
  id: string;
  label: string;
};

export type Props = {
  placeholder?: string;
  initialValue?: string[];
  onQuery: (query: string) => Promise<CombotagsItem[]>;
  onChange: (values: string[]) => Promise<void>;
};

export type CombotagsHandle = {
  append: (tags: string[]) => void;
  setValue: (tags: string[]) => void;
};

// Combotags provides a mix of a tags input and a combobox where the tags input
// field is used to filter the combobox results. The combobox results are then
// used to add items to the tags input.
export const Combotags = forwardRef<CombotagsHandle, Props>((props, ref) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<CombotagsItem[]>([]);
  const [isComboboxOpen, setIsComboboxOpen] = useState(false);
  const [selectedItems, setSelectedItems] = useState<Map<string, string>>(
    new Map(),
  );

  // NOTE: Because we're combining the combobox with a tags input, we need to
  // use the context provider here for easier control of the tags input values.
  const tagsInput = useTagsInput({
    defaultValue: props.initialValue,
    inputValue: searchQuery,
    addOnPaste: true,
    onInputValueChange: handleInputValueChange,
    onValueChange: handleValueChange,
    onInteractOutside: handleInteractOutside,
  });

  // Used by the combobox event handler to update the values of the tags input.
  const tagsInputRef = useRef(tagsInput);
  useEffect(() => {
    tagsInputRef.current = tagsInput;
  }, [tagsInput]);

  useImperativeHandle<CombotagsHandle, CombotagsHandle>(ref, () => {
    function append(tags: string[]) {
      const newValue = uniq([...tagsInputRef.current.value, ...tags]);
      tagsInputRef.current.setValue(newValue);
    }

    function setValue(tags: string[]) {
      tagsInputRef.current.setValue(tags);
    }

    return {
      append,
      setValue,
    };
  });

  // Used for positioning the combobox by computing the height of the input.
  const inputControlRef = useRef<HTMLDivElement>(null);

  async function handleInputValueChange({ inputValue }) {
    setSearchQuery(inputValue);
    const result = await props.onQuery(inputValue);

    setSearchResults(result);
    setIsComboboxOpen(() => true);
  }

  async function handleValueChange({ value }) {
    // Update selected items map
    const newMap = new Map<string, string>();
    value.forEach((id: string) => {
      newMap.set(id, selectedItems.get(id) || id);
    });
    setSelectedItems(newMap);

    // Immediately update the local list, filtering out the newly added values.
    setSearchResults(searchResults.filter((item) => !value.includes(item.id)));

    // NOTE: Not awaited to facilitate optimistic updates.
    props.onChange(value);
  }

  function handleInteractOutside() {
    setIsComboboxOpen(() => false);
  }

  const handleSelect = (item: CombotagsItem) => () => {
    if (tagsInputRef.current.value.includes(item.id)) {
      // Don't add duplicates.
      return;
    }

    // Update selected items map with label
    setSelectedItems((prev) => new Map(prev).set(item.id, item.label));

    // This is necessary because `addValue` is broken at the moment.
    const newValue = [...tagsInputRef.current.value, item.id];
    tagsInputRef.current.setValue(newValue);

    setSearchQuery("");
  };

  const collection = createListCollection({
    items: searchResults,
    itemToValue(item) {
      return item.id;
    },
  });
  const rect = inputControlRef.current?.getBoundingClientRect();
  const offset = rect?.height ?? 0;

  return (
    <>
      <TagsInput.RootProvider
        value={tagsInput}
        position="relative"
        w="full"
        size="sm"
      >
        <TagsInput.Context>
          {(api) => (
            <>
              <TagsInput.Control ref={inputControlRef}>
                {api.value.map((id, index) => (
                  <TagsInput.Item key={index} index={index} value={id}>
                    <TagsInput.ItemPreview>
                      <TagsInput.ItemText>
                        {selectedItems.get(id) || id}
                      </TagsInput.ItemText>
                      <TagsInput.ItemDeleteTrigger asChild>
                        <IconButton variant="link" size="xs">
                          <DeleteSmallIcon />
                        </IconButton>
                      </TagsInput.ItemDeleteTrigger>
                    </TagsInput.ItemPreview>
                    <TagsInput.ItemInput />
                    <TagsInput.HiddenInput />
                  </TagsInput.Item>
                ))}
                <TagsInput.Input placeholder={props.placeholder || "Tags..."} />
              </TagsInput.Control>
            </>
          )}
        </TagsInput.Context>

        <Combobox.Root
          position="absolute"
          style={{
            top: offset,
          }}
          open={isComboboxOpen}
          autoFocus={false}
          collection={collection}
          inputValue={searchQuery}
          inputBehavior="autohighlight"
          size="sm"
        >
          <Combobox.Content maxH="64" overflowY="scroll">
            {searchResults.map((item) => (
              <Combobox.Item
                key={item.id}
                id={item.id}
                item={item}
                onClick={handleSelect(item)}
              >
                <Combobox.ItemText
                  alignItems="center"
                  lineHeight="tight"
                  textWrap="nowrap"
                  overflow="hidden"
                  textOverflow="ellipsis"
                >
                  {item.label}
                </Combobox.ItemText>
              </Combobox.Item>
            ))}
          </Combobox.Content>
        </Combobox.Root>
      </TagsInput.RootProvider>
    </>
  );
});

Combotags.displayName = "Combotags";
