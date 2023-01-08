import { Flex } from "@chakra-ui/react";
import { SearchBar } from "../SearchBar/SearchBar";
import { HomeLink } from "./components/HomeLink";
import { Options } from "./components/Options/Options";

export function Navigation() {
  return (
    <Flex py="1em" width="full" justifyContent="center" bgColor="teal.200">
      <Flex
        width="full"
        px={4}
        maxW="container.lg"
        justifyContent="space-between"
        alignItems="center"
      >
        <HomeLink />

        <SearchBar />

        <Options />
      </Flex>
    </Flex>
  );
}
