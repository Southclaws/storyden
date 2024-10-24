import { createListCollection, useTagsInput } from "@ark-ui/react";
import { without } from "lodash";
import { XIcon } from "lucide-react";
import { useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import * as Combobox from "@/components/ui/combobox";
import { IconButton } from "@/components/ui/icon-button";
import { Input } from "@/components/ui/input";
import * as TagsInput from "@/components/ui/tags-input";

export type Props = {
  onQuery: (query: string) => Promise<string[]>;
};

export function Combotags(props: Props) {
  const [query, setQuery] = useState("");
  const [items, setItems] = useState<string[]>([]);
  const [isComboboxOpen, setIsComboboxOpen] = useState(false);

  const tagsInput = useTagsInput({
    inputValue: query,
    onInputValueChange: async ({ inputValue }) => {
      setQuery(inputValue);
      const result = await props.onQuery(inputValue);

      const currentItems = [...tagsInput.value];

      const newItems = without(result, ...currentItems);

      console.log({ inputValue, result, currentItems, newItems });

      setItems(newItems);
      setIsComboboxOpen(() => true);
    },
  });

  const collection = createListCollection({ items });

  const handleSelect = (value: string) => () => {
    tagsInput.addValue(value);
    setQuery("");
    setIsComboboxOpen(() => false);
  };

  const ref = useRef<HTMLDivElement>(null);

  const rect = ref.current?.getBoundingClientRect();
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
              <TagsInput.Control ref={ref}>
                {api.value.map((value, index) => (
                  <TagsInput.Item key={index} index={index} value={value}>
                    <TagsInput.ItemPreview>
                      <TagsInput.ItemText>{value}</TagsInput.ItemText>
                      <TagsInput.ItemDeleteTrigger asChild>
                        <IconButton variant="link" size="xs">
                          <XIcon />
                        </IconButton>
                      </TagsInput.ItemDeleteTrigger>
                    </TagsInput.ItemPreview>
                    <TagsInput.ItemInput />
                    <TagsInput.HiddenInput />
                  </TagsInput.Item>
                ))}
                <TagsInput.Input placeholder="Tags..." />
              </TagsInput.Control>
            </>
          )}
        </TagsInput.Context>

        <Combobox.Root
          position="absolute"
          style={{
            top: offset,
          }}
          positioning={{}}
          open={isComboboxOpen}
          autoFocus={false}
          collection={collection}
          inputValue={query}
          // highlightedValue={highlighted?.id}
          inputBehavior="autohighlight"
          size="lg"
        >
          <Combobox.Content maxH="64" overflowY="scroll">
            {items.map((item) => (
              <Combobox.Item
                key={item}
                id={item}
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
                  {item}
                </Combobox.ItemText>
              </Combobox.Item>
            ))}
          </Combobox.Content>
        </Combobox.Root>
      </TagsInput.RootProvider>
    </>
  );
}
