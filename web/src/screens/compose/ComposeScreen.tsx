import {
  Button,
  FormErrorMessage,
  HStack,
  Heading,
  Input,
  VStack,
} from "@chakra-ui/react";
import { Controller } from "react-hook-form";
import { Editor } from "src/components/Editor";
import { CategorySelect } from "./components/CategorySelect/CategorySelect";
import { useComposeScreen } from "./useComposeScreen";

export function ComposeScreen() {
  const {
    onSubmit,
    handleSubmit,
    control,
    register,
    isValid,
    errors,
    isSubmitting,
  } = useComposeScreen();

  return (
    <VStack
      alignItems="start" //
      gap={2}
      w="full"
      py={5}
    >
      <Heading>New thread</Heading>

      <VStack
        as="form"
        onSubmit={handleSubmit(onSubmit)}
        alignItems="start"
        w="full"
        gap={2}
      >
        <Input placeholder="Title" {...register("title")} />
        <FormErrorMessage>{errors.title?.message}</FormErrorMessage>

        <Controller
          render={({ field }) => <Editor onChange={field.onChange} />}
          control={control}
          name="body"
        />
        <FormErrorMessage>{errors.body?.message}</FormErrorMessage>

        <HStack>
          <Controller
            render={({ field }) => (
              <>
                <CategorySelect {...field} />
                <FormErrorMessage>{errors.category?.message}</FormErrorMessage>
              </>
            )}
            control={control}
            name="category"
          />

          <Button
            type="submit" //
            isDisabled={!isValid}
            isLoading={isSubmitting}
          >
            Post
          </Button>
        </HStack>
      </VStack>
    </VStack>
  );
}
