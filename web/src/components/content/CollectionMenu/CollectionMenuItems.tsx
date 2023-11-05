import { MinusIcon, PlusIcon } from "@heroicons/react/24/solid";

import { Checkbox } from "src/theme/components/Checkbox";
import {
  MenuContent,
  MenuItem,
  MenuItemGroup,
  MenuItemGroupLabel,
  MenuSeparator,
} from "src/theme/components/Menu";

import { CollectionCreateTrigger } from "../CollectionCreate/CollectionCreateTrigger";

import { Center, HStack } from "@/styled-system/jsx";

import { Props, useCollectionMenuItems } from "./useCollectionMenuItems";

export function CollectionMenuItems(props: Props) {
  const { collections, onSelect, multiSelect } = useCollectionMenuItems(props);

  return (
    <MenuContent>
      <MenuItemGroup id="group">
        <MenuItemGroupLabel htmlFor="group">
          Add to collections
        </MenuItemGroupLabel>

        <MenuSeparator />

        {collections.map((c) => (
          <MenuItem id={c.id} key={c.id} onClick={onSelect(c)}>
            <HStack>
              {multiSelect ? (
                <Checkbox checked={c.hasPost} />
              ) : (
                <Center w="5">
                  {c.hasPost ? <MinusIcon /> : <PlusIcon />}
                </Center>
              )}
              {c.name}
            </HStack>
          </MenuItem>
        ))}
      </MenuItemGroup>

      <MenuItemGroup id="create">
        <MenuItem id="create-collection" closeOnSelect={false}>
          <CollectionCreateTrigger kind="blank" />
        </MenuItem>
      </MenuItemGroup>
    </MenuContent>
  );
}
