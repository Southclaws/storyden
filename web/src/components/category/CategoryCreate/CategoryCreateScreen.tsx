import { Button } from "src/theme/components/Button";
import { FormControl } from "src/theme/components/FormControl";
import { FormHelperText } from "src/theme/components/FormHelperText";
import { FormLabel } from "src/theme/components/FormLabel";
import { Input } from "src/theme/components/Input";
import { UseDisclosureProps } from "src/utils/useDisclosure";

import { HStack, VStack, styled } from "@/styled-system/jsx";

import { useCategoryCreate } from "./useCategoryCreate";

export function CategoryCreateScreen(props: UseDisclosureProps) {
  const { register, onSubmit } = useCategoryCreate(props);

  return (
    <VStack alignItems="start" gap="4">
      <styled.p>
        Use categories to organise posts. A post can only have one category,
        unlike tags. So it&apos;s best to keep categories high-level and
        different enough so that it&apos;s not easy to get confused between
        them.
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
          <FormHelperText>The name for your category</FormHelperText>
        </FormControl>
        <FormControl>
          <FormLabel>Description</FormLabel>

          {/* TODO: Make a larger textarea component for this. */}
          <Input {...register("description")} type="text" />
          <FormHelperText>Describe your category</FormHelperText>
        </FormControl>

        <HStack w="full" justify="space-between">
          <Button w="full" type="submit" onClick={props.onClose}>
            Cancel
          </Button>
          <Button w="full" type="submit" kind="primary">
            Create
          </Button>
        </HStack>
      </styled.form>
    </VStack>
  );
}
