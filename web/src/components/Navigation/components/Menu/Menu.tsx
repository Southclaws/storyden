import { Flex, HStack, IconButton } from "@chakra-ui/react";
import { Cog8ToothIcon, UserIcon } from "@heroicons/react/24/outline";
import map from "lodash/fp/map";
import { Category } from "src/api/openapi/schemas";
import { Anchor } from "src/components/site/Anchor";
import { Unready } from "src/components/Unready";
import { useMenu } from "./useMenu";

const mapCategories = map((c: Category) => (
  <Anchor key={c.id} href={`/c/${c.name}`}>
    {c.name}
  </Anchor>
));

export function Menu() {
  const { isAuthenticated, data, error } = useMenu();

  if (!data) return <Unready {...error} />;

  return (
    <Flex
      maxHeight="full"
      flexDir="column"
      justifyContent="center"
      alignItems="center"
      maxW="container.sm"
      width="full"
      p={4}
      borderRadius={32}
      bgColor="hsla(220, 15%, 95%, 0.75)"
      backdropFilter={`blur(1em)`}
      pointerEvents="auto"
    >
      <Flex
        height="full"
        width="full"
        gap={2}
        flexDir="column"
        justifyContent="space-between"
        alignItems="start"
        overflowY="scroll"
      >
        {mapCategories(data.categories)}
      </Flex>

      {isAuthenticated && (
        <HStack width="full" justifyContent="space-between">
          <Anchor variant="outline" size="sm">
            <IconButton
              aria-label="Settings"
              borderRadius="50%"
              icon={<Cog8ToothIcon width="1em" />}
            />
          </Anchor>
          <Anchor variant="outline" size="sm">
            <IconButton
              aria-label="Profile"
              borderRadius="50%"
              icon={<UserIcon width="1em" />}
            />
          </Anchor>
        </HStack>
      )}
    </Flex>
  );
}
