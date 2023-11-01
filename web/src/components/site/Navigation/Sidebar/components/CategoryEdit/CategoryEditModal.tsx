import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";
import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  HStack,
  Input,
  VStack,
} from "src/theme/components";

import { Props, useCategoryEdit } from "./useCategoryEdit";

export function CategoryEditModal(props: Props) {
  const { register, onSubmit, errors } = useCategoryEdit(props);

  return (
    <ModalDrawer
      isOpen={props.isOpen}
      onClose={props.onClose}
      title="Edit category"
    >
      <VStack
        as="form"
        justify="space-between"
        align="start"
        height="full"
        onSubmit={onSubmit}
      >
        <VStack w="full">
          <FormControl>
            <FormLabel>Name</FormLabel>
            <Input {...register("name")} type="text" />
            <FormErrorMessage>{errors["name"]?.message}</FormErrorMessage>
          </FormControl>

          <FormControl>
            <FormLabel>Description</FormLabel>
            <Input {...register("description")} type="text" />
            <FormErrorMessage>
              {errors["description"]?.message}
            </FormErrorMessage>
          </FormControl>
        </VStack>

        <HStack w="full" align="center" justify="end" pb={3} gap={4}>
          <Button variant="outline" size="sm" onClick={props.onClose}>
            Cancel
          </Button>
          <Button colorScheme="green" size="sm" type="submit">
            Save
          </Button>
        </HStack>
      </VStack>
    </ModalDrawer>
  );
}
