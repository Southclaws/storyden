import { Flex } from "@chakra-ui/react";
import { useSession } from "src/auth";
import { Profile } from "../Profile/Profile";
import { SearchBar } from "../SearchBar/SearchBar";
import { HomeLink } from "./HomeLink";

export function Navigation() {
  return (
    <Flex py="1em" width="full" justifyContent="center" bgColor="teal.200">
      <Flex
        width="full"
        px={2}
        maxW="container.lg"
        justifyContent="space-around"
        alignItems="center"
      >
        <HomeLink />

        <SearchBar />

        <Profile />
      </Flex>
    </Flex>
  );
}
