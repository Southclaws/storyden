import { Box, Heading, List, ListItem, Text, VStack } from "@chakra-ui/react";
import { KeyIcon } from "@heroicons/react/20/solid";
import { useAuthProviderList } from "src/api/openapi/auth";
import { AuthProvider } from "src/api/openapi/schemas";
import { AuthSelectionOption } from "./AuthSelectionOption";

export function AuthSelection() {
  const { data, error } = useAuthProviderList();

  if (error) throw error;

  return (
    <VStack gap={4}>
      <Box>
        <Heading size="md">
          Sign up
          <br />
        </Heading>
        <Text size="sm" fontWeight="medium" color="blackAlpha.600">
          or sign in
        </Text>
      </Box>

      <List display="flex" flexDir="column" gap={4}>
        {data?.map((v: AuthProvider) => (
          <ListItem key={v.provider}>
            <AuthSelectionOption
              name={v.name}
              method={v.provider}
              icon={<KeyIcon height="100%" />}
              link={v.link || undefined}
            />
          </ListItem>
        ))}
      </List>
    </VStack>
  );
}
