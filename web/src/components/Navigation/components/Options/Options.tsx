import {
  Button,
  HStack,
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
import { ChevronDownIcon, UserIcon } from "@heroicons/react/24/solid";
import { Avatar } from "../../../Avatar/Avatar";
import { useOptions } from "./useOptions";

export function Options() {
  const profile = useOptions();

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
        as={Button}
        padding={1}
        aria-label="menu"
        variant="ghost"
        _hover={{ bgColor: "whiteAlpha.500" }}
      >
        <HStack>
          {/* NOTE: Custom sizing tokens here */}
          <Avatar account={profile.account} boxSize="1.5em" />
          <ChevronDownIcon width="1em" />
        </HStack>
      </MenuButton>

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
