import { UseDisclosureProps } from "src/utils/useDisclosure";

import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormHelperText } from "@/components/ui/form/FormHelperText";
import { FormLabel } from "@/components/ui/form/FormLabel";
import { Input } from "@/components/ui/input";
import { VStack, WStack, styled } from "@/styled-system/jsx";

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

        <WStack>
          <Button flexGrow="1" type="button" onClick={props.onClose}>
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
