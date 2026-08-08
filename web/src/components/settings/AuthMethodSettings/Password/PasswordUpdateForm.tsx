import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { Input } from "@/components/ui/input";
import { SectionHeading } from "@/components/ui/section-heading";
import { styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { usePasswordUpdate } from "./usePasswordUpdate";

export function PasswordUpdateForm() {
  const { form, handlePasswordChange } = usePasswordUpdate();

  return (
    <styled.form className={lstack()} gap="2" onSubmit={handlePasswordChange}>
      <SectionHeading>Password</SectionHeading>

      <FormControl>
        <Input
          maxW="xs"
          type="password"
          autoComplete="current-password"
          placeholder="current password"
          {...form.register("old")}
        />
        <FormErrorText>{form.formState.errors["old"]?.message}</FormErrorText>
      </FormControl>

      <FormControl>
        <Input
          maxW="xs"
          type="password"
          autoComplete="new-password"
          placeholder="new password"
          {...form.register("new")}
        />
        <FormErrorText>{form.formState.errors["new"]?.message}</FormErrorText>
        <FormErrorText>{form.formState.errors["root"]?.message}</FormErrorText>
      </FormControl>

      <Button type="submit" variant="subtle" size="sm">
        Change password
      </Button>
    </styled.form>
  );
}
