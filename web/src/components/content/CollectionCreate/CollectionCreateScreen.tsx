import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form-control";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { Input } from "@/components/ui/input";
import { Text } from "@/components/ui/text";
import { VStack, WStack, styled } from "@/styled-system/jsx";
import { UseDisclosureProps } from "@/utils/useDisclosure";

import { Props, useCollectionCreate } from "./useCollectionCreate";

export function CollectionCreateScreen(props: Props) {
  const { register, onSubmit } = useCollectionCreate(props);

  return (
    <VStack alignItems="start" gap="4">
      <Text variant="supporting">
        Use collections to curate content from the community. Collections can
        include threads, pages and other items from the community knowledgebase.
      </Text>
      <styled.form
        display="flex"
        flexDir="column"
        gap="4"
        w="full"
        onSubmit={onSubmit}
      >
        <FormControl>
          <FormLabel>Name*</FormLabel>
          <Input {...register("name")} type="text" />
          <FormHelperText>The name for your collection</FormHelperText>
        </FormControl>
        <FormControl>
          <FormLabel>Description</FormLabel>

          {/* TODO: Make a larger textarea component for this. */}
          <Input {...register("description")} type="text" />
          <FormHelperText>
            Optional description for your collection.
          </FormHelperText>
        </FormControl>

        <WStack>
          <Button
            flexGrow="1"
            type="button"
            variant="outline"
            onClick={props.onClose}
          >
            Cancel
          </Button>
          <Button flexGrow="1" type="submit">
            Create
          </Button>
        </WStack>
      </styled.form>
    </VStack>
  );
}
