import { Box, Divider, Flex, VStack } from "@chakra-ui/react";
import { useNavigation } from "../useNavigation";
import { Authbar } from "./components/Authbar";
import { CategoryList } from "./components/CategoryList";
import { Title } from "./components/Title";
import { Toolbar } from "./components/Toolbar";

export function Sidebar() {
  const { title } = useNavigation();

  return (
    <Flex
      id="desktop-nav-container"
      as="header"
      justifyContent="end"
      width="full"
      height="full"
    >
      <Box
        id="desktop-nav-box"
        maxWidth="2xs"
        minWidth={{
          base: "full",
          lg: "3xs",
        }}
        height="full"
      >
        <VStack
          as="nav"
          height="full"
          py={4}
          gap={2}
          justifyContent="space-between"
          alignItems="start"
        >
          <VStack width="full" alignItems="start" overflow="hidden">
            <Title>{title}</Title>

            <Toolbar />

            <Divider />

            <Box overflowY="scroll" width="full">
              <CategoryList />
            </Box>
          </VStack>

          <Divider />

          <VStack alignItems="start">
            <Authbar />
          </VStack>
        </VStack>
      </Box>
    </Flex>
  );
}
