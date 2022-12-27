import { Box, HStack } from "@chakra-ui/react";
import { Profile } from "../Profile/Profile";
import { SearchBar } from "../SearchBar/SearchBar";
import { StorydenLogo } from "../StorydenLogo";

export function Navigation() {
  return (
    <HStack
      width="full"
      padding="1em"
      justifyContent="center"
      bgColor="#E5E5E5"
    >
      <HStack width="full" maxW="container.lg" justifyContent="space-around">
        <StorydenLogo />

        <SearchBar />

        <Profile />
      </HStack>
    </HStack>
  );
}
