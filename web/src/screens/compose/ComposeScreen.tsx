import { Button, Input, VStack } from "@chakra-ui/react";
import { Controller, useForm } from "react-hook-form";
import { Editor } from "src/components/Editor";
import { useComposeScreen } from "./useComposeScreen";

export function ComposeScreen() {
  const { onCreate } = useComposeScreen();

  const { handleSubmit, control, register } = useForm();

  const onSubmit = (data: any) => {
    onCreate(data.title, "", data.md);
  };

  return (
    <VStack w="full" py={5}>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Input placeholder="Title" {...register("title")} />
        <Controller
          render={({ field }) => <Editor onChange={field.onChange} />}
          control={control}
          name="body"
        />

        <Button type="submit">Post</Button>
      </form>
    </VStack>
  );
}
