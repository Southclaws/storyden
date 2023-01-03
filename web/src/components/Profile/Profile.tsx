import {
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
} from "@chakra-ui/react";
import { UserIcon } from "@heroicons/react/24/solid";
import Link from "../site/Link";
import { useProfile } from "./useProfile";

export function Profile() {
  const { authenticated } = useProfile();

  if (!authenticated) {
    return (
      <Link href="/login" padding="0.1em" _hover={{ cursor: "pointer" }}>
        <UserIcon aria-label="authenticate" width="2em" />
      </Link>
    );
  }

  return (
    <Menu>
      <MenuButton
        as={IconButton}
        padding={1}
        aria-label="Profile"
        icon={<UserIcon />}
        variant="outline"
      />

      <MenuList>
        <MenuItem>Profile</MenuItem>
      </MenuList>
    </Menu>
  );
}
