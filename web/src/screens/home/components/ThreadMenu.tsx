import { Menu, MenuButton, MenuItem, MenuList } from "@chakra-ui/react";
import { Bars3Icon } from "@heroicons/react/24/solid";
import { Identifier } from "src/api/openapi/schemas";

type Props = {
  postID: Identifier;
};

// TODO: Implement API endpoints for administrative actions.
export function ThreadMenu(props: Props) {
  return (
    <Menu>
      <MenuButton zIndex={100} aria-label="Options" boxSize="1.2em">
        <Bars3Icon />
      </MenuButton>

      <MenuList>
        <MenuItem>Delete</MenuItem>
      </MenuList>
    </Menu>
  );
}
