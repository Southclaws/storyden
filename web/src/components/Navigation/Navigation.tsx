import {
  Box,
  Button,
  Flex,
  HStack,
  IconButton,
  Input,
  SlideFade,
  VStack,
} from "@chakra-ui/react";
import {
  BellIcon,
  Cog8ToothIcon,
  PlusIcon,
  XMarkIcon,
} from "@heroicons/react/24/outline";
import { map } from "lodash/fp";
import { Category } from "src/api/openapi/schemas";
import { Anchor } from "../site/Anchor";
import { Unready } from "../Unready";
import { MenuIcon } from "./components/MenuIcon";
import { useNavigation } from "./useNavigation";

const mapCategories = (selected?: string) =>
  map((c: Category) => (
    <Anchor key={c.id} href={`/c/${c.name}`} w="full">
      <Button bgColor={c.name === selected ? "blackAlpha.200" : ""} w="full">
        {c.name}
      </Button>
    </Anchor>
  ));

export function Navigation() {
  const { error, isAuthenticated, isExpanded, onExpand, categories, category } =
    useNavigation();

  if (error) return <Unready {...error} />;

  return (
    <Box
      id="navigation-overlay"
      position="fixed"
      bottom="env(safe-area-inset-bottom)"
      width="100vw"
      height="100vh"
      pointerEvents="none"
      zIndex="overlay"
    >
      <Flex
        id="flex-container"
        height="full"
        gap={3}
        p={2}
        justifyContent="end"
        alignItems="center"
        flexDir="column"
      >
        <Flex
          px={{ base: isExpanded ? 2 : 4 }}
          py={{ base: 2 }}
          gap={2}
          flexDirection="column"
          borderRadius={{ base: 16 }}
          backdropFilter="blur(1em)"
          transitionProperty="background-color"
          transitionDuration="0.5s"
          backgroundColor="hsla(210, 38.5%, 94.9%, 0.8)"
          border="1px solid hsla(209, 100%, 20%, 0.02)"
          width={isExpanded ? "100%" : "min-content"}
          maxW={{ base: "23em", md: "container.sm" }}
          justifyContent="space-between"
          alignItems="center"
          pointerEvents="auto"
        >
          <SlideFade
            in={isExpanded}
            style={{
              maxHeight: "100%",
              width: "100%",
              display: isExpanded ? "flex" : "none",
            }}
          >
            <VStack width="full" p={2}>
              <Flex
                maxHeight="80vh"
                flexDir="column"
                justifyContent="center"
                alignItems="center"
                maxW="container.sm"
                width="full"
                pointerEvents="auto"
                gap={2}
              >
                <HStack
                  width="full"
                  justifyContent={isAuthenticated ? "space-between" : "end"}
                >
                  {isAuthenticated && (
                    <Anchor variant="outline" size="sm">
                      <IconButton
                        aria-label="Settings"
                        borderRadius="50%"
                        icon={<Cog8ToothIcon width="1em" />}
                      />
                    </Anchor>
                  )}

                  <IconButton
                    aria-label="Close menu"
                    borderRadius="50%"
                    icon={<XMarkIcon width="1em" />}
                    onClick={onExpand}
                    colorScheme="gray"
                    size="xs"
                  />
                </HStack>

                <Flex
                  height="full"
                  width="full"
                  gap={2}
                  flexDir="column"
                  justifyContent="space-between"
                  alignItems="start"
                  overflowY="scroll"
                >
                  {mapCategories(category)(categories)}
                </Flex>
              </Flex>
            </VStack>
          </SlideFade>

          <HStack gap={4} w="full" justifyContent="space-between">
            <Anchor href="/notifications">
              <IconButton
                aria-label=""
                borderRadius={12}
                icon={<BellIcon width="1em" />}
              />
            </Anchor>

            {isExpanded ? (
              <Input
                variant="outline"
                border="none"
                placeholder="Search anything..."
              />
            ) : (
              <IconButton
                aria-label="main menu"
                borderRadius={12}
                icon={<MenuIcon />}
                onClick={onExpand}
              />
            )}

            <Anchor href="/new">
              <IconButton
                aria-label=""
                borderRadius={12}
                icon={<PlusIcon width="1em" />}
              />
            </Anchor>
          </HStack>
        </Flex>
      </Flex>
    </Box>
  );
}
