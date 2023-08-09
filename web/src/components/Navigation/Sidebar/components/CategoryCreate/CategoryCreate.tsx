import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverContent,
  PopoverFooter,
  PopoverTrigger,
  Portal,
} from "@chakra-ui/react";
import { PlusIcon } from "@heroicons/react/24/solid";

import { useCategoryCreate } from "./useCategoryCreate";

export function CategoryCreate() {
  const { register, onSubmit, errors } = useCategoryCreate();
  return (
    <Popover placement="right">
      <PopoverTrigger>
        <Button
          w="full"
          size="xs"
          variant="outline"
          leftIcon={<PlusIcon width="1.125em" />}
        >
          New category
        </Button>
      </PopoverTrigger>
      <Portal>
        <PopoverContent
          as="form"
          backgroundColor="cyan.50"
          boxShadow="lg"
          onSubmit={onSubmit}
        >
          <PopoverArrow />
          <PopoverBody display="flex" flexDirection="column" gap={4}>
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
          </PopoverBody>
          <PopoverFooter
            border="0"
            display="flex"
            alignItems="center"
            justifyContent="end"
            pb={3}
            gap={4}
          >
            <Button variant="outline" size="sm">
              Cancel
            </Button>
            <Button colorScheme="green" size="sm" type="submit">
              Create
            </Button>
          </PopoverFooter>
        </PopoverContent>
      </Portal>
    </Popover>
  );
}
