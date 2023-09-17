import { Menu, MenuButton, MenuItem, MenuList } from "@chakra-ui/react";

import { Settings } from "../site/Action/Action";

export function SettingsMenu() {
  return (
    <Menu>
      <MenuButton as={Settings}>Settings</MenuButton>

      <MenuList>
        <MenuItem>Edit mode</MenuItem>
      </MenuList>
    </Menu>
  );
}
