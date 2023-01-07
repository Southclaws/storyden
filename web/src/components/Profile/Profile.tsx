import {
  Box,
  HStack,
  IconButton,
  Link,
  Menu,
  MenuButton,
  MenuDivider,
  MenuItem,
  MenuList,
  Text,
  VStack,
} from "@chakra-ui/react";
import { LockClosedIcon } from "@heroicons/react/20/solid";
import { UserIcon } from "@heroicons/react/24/solid";
import { Avatar } from "../Avatar/Avatar";
import { useProfile } from "./useProfile";

export function Profile() {
  const profile = useProfile();

  if (!profile.authenticated) {
    return (
      <Link href="/auth" padding="0.1em" _hover={{ cursor: "pointer" }}>
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
        variant="ghost"
      />

      <MenuList>
        <HStack alignItems="center" spacing={4} px={4} py={1}>
          <Avatar account={profile.account} boxSize={8} />

          <VStack spacing={0} alignItems="flex-start">
            <Text as="span" fontSize="md">
              {profile.account.name}
            </Text>
            <Text as="span" fontSize="xs" color="grey.gunsmoke">
              {profile.account.handle}
            </Text>
          </VStack>
        </HStack>

        <MenuDivider />

        <Link href="/logout">
          <MenuItem icon={<LockClosedIcon />}>Logout</MenuItem>
        </Link>
      </MenuList>
    </Menu>
  );
}
