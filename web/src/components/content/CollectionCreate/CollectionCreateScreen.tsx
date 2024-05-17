import { Button } from "src/theme/components/Button";
import { FormControl } from "src/theme/components/FormControl";
import { FormHelperText } from "src/theme/components/FormHelperText";
import { FormLabel } from "src/theme/components/FormLabel";
import { Input } from "src/theme/components/Input";
import { UseDisclosureProps } from "src/utils/useDisclosure";

import { HStack, VStack, styled } from "@/styled-system/jsx";

import { useCollectionCreate } from "./useCollectionCreate";

export function CollectionCreateScreen(props: UseDisclosureProps) {
  const { register, onSubmit } = useCollectionCreate(props);

  return (
    <VStack alignItems="start" gap="4">
      <styled.p>
        Use collections to curate content from the community. Collections can
        include threads, posts and other items from the community database.
      </styled.p>
      <styled.form
        display="flex"
        flexDir="column"
        gap="4"
        w="full"
        onSubmit={onSubmit}
      >
        <FormControl>
          <FormLabel>Name</FormLabel>
          <Input {...register("name")} type="text" />
          <FormHelperText>The name for your collection</FormHelperText>
        </FormControl>
        <FormControl>
          <FormLabel>Description</FormLabel>

          {/* TODO: Make a larger textarea component for this. */}
          <Input {...register("description")} type="text" />
          <FormHelperText>Describe your collection</FormHelperText>
        </FormControl>

        <HStack w="full" justify="space-between">
          <Button w="full" type="submit">
            Cancel
          </Button>
          <Button w="full" type="submit">
            Create
          </Button>
        </HStack>
      </styled.form>
    </VStack>
  );
}
