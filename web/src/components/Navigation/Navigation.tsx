import { Flex } from "@chakra-ui/react";
import { Profile } from "../Profile/Profile";
import { SearchBar } from "../SearchBar/SearchBar";
import { StorydenLogo } from "../StorydenLogo";

export function Navigation() {
  return (
    <Flex py="1em" width="full" justifyContent="center" bgColor="#E5E5E5">
      <Flex
        width="full"
        px={2}
        maxW="container.lg"
        justifyContent="space-around"
      >
        <StorydenLogo />

        <SearchBar />

        <Profile />
      </Flex>
    </Flex>
  );
}
