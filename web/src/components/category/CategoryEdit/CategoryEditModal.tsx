import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";
import { Button } from "src/theme/components/Button";
import { FormControl } from "src/theme/components/FormControl";
import { FormFeedback } from "src/theme/components/FormFeedback";
import { FormLabel } from "src/theme/components/FormLabel";
import { Input } from "src/theme/components/Input";

import { HStack, VStack, styled } from "@/styled-system/jsx";

import { Props, useCategoryEdit } from "./useCategoryEdit";

export function CategoryEditModal(props: Props) {
  const { register, onSubmit, errors } = useCategoryEdit(props);

  return (
    <ModalDrawer
      isOpen={props.isOpen}
      onClose={props.onClose}
      title="Edit category"
    >
      <styled.form
        display="flex"
        flexDir="column"
        justifyContent="space-between"
        alignItems="start"
        height="full"
        onSubmit={onSubmit}
        gap="2"
      >
        <VStack w="full">
          <FormControl>
            <FormLabel>Name</FormLabel>
            <Input {...register("name")} type="text" />
            <FormFeedback error={errors["name"]?.message}>
              The name of the category.
            </FormFeedback>
          </FormControl>

          <FormControl>
            <FormLabel>Description</FormLabel>
            <Input {...register("description")} type="text" />
            <FormFeedback error={errors["description"]?.message}>
              The description for the category.
            </FormFeedback>
          </FormControl>
        </VStack>

        <HStack w="full" alignItems="center" justify="end" pb="3" gap="4">
          <Button kind="ghost" size="sm" onClick={props.onClose}>
            Cancel
          </Button>
          <Button kind="primary" size="sm" type="submit">
            Save
          </Button>
        </HStack>
      </styled.form>
    </ModalDrawer>
  );
}
