import { UseDisclosureProps } from "src/utils/useDisclosure";

import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { Input } from "@/components/ui/input";
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
          <Button w="full" type="button">
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
