import { List, ListItem, VStack } from "@chakra-ui/react";
import { AuthProvider } from "src/api/openapi/schemas";
import { Unready } from "src/components/Unready";
import { AuthSelectionOption } from "../AuthSelectionOption";
import { useAuthSelection } from "./useAuthSelection";

export function AuthSelection() {
  const { data, error } = useAuthSelection();

  if (!data) <Unready {...error} />;

  // sort by alphabetical because lazy
  // TODO: allow the order to be configured by the admin.
  data?.providers?.sort((a, b) => a.provider.localeCompare(b.provider));

  return (
    <VStack w="full" gap={4}>
      <List display="flex" flexDir="column" gap={4} w="full">
        {data?.providers?.map((v: AuthProvider) => (
          <ListItem key={v.provider}>
            <AuthSelectionOption
              name={v.name}
              method={v.provider}
              icon={v.logo_url}
              link={v.link || undefined}
            />
          </ListItem>
        ))}
      </List>
    </VStack>
  );
}
