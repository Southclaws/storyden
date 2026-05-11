import { Admonition } from "@/components/ui/admonition";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useI18n } from "@/i18n/provider";
import { VStack, styled } from "@/styled-system/jsx";

import { usePasswordCreate } from "./usePasswordCreate";

export function PasswordCreateForm() {
  const {
    form: { register, handlePasswordCreate, errors },
    success,
    handleCloseNotification,
  } = usePasswordCreate();
  const { t } = useI18n();

  return (
    <>
      <p>
        {t(
          "Your account does not currently have a password. You can add a password here.",
        )}
      </p>
      <styled.form
        display="flex"
        flexDir="column"
        w="full"
        gap="2"
        onSubmit={handlePasswordCreate}
      >
        <Input
          maxW="xs"
          type="password"
          autoComplete="new-password"
          placeholder={t("new password")}
          disabled={success}
          {...register("newPassword")}
        />
        <styled.p color="fg.error" fontSize="sm">
          {errors.newPassword?.message}
        </styled.p>
        <Input
          maxW="xs"
          type="password"
          autoComplete="new-password"
          placeholder={t("confirm new password")}
          disabled={success}
          {...register("confirmPassword")}
        />
        <styled.p color="fg.error" fontSize="sm">
          {errors.confirmPassword?.message}
        </styled.p>
        <styled.p color="fg.error" fontSize="sm">
          {errors.root?.message}
        </styled.p>
        <VStack alignItems="start" w="full">
          <Button type="submit" disabled={success}>
            {t("Add password")}
          </Button>
          <Admonition
            value={success}
            onChange={handleCloseNotification}
            kind="success"
            title={t("Success")}
          >
            {t(
              "Your account now has a password! You can now use this to log in, but you can continue to use your other authentication methods as well.",
            )}
          </Admonition>
        </VStack>
      </styled.form>
    </>
  );
}
