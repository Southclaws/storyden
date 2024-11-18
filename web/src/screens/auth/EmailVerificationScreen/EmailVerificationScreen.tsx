"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { handle } from "@/api/client";
import { authEmailVerify } from "@/api/openapi-client/auth";
import { Account } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form/FormControl";
import { FormErrorText } from "@/components/ui/form/FormErrorText";
import { PinInputField } from "@/components/ui/form/PinInputField";
import { Heading } from "@/components/ui/heading";
import { Input } from "@/components/ui/input";
import { vstack } from "@/styled-system/patterns";

export const FormSchema = z.object({
  email: z.string().email("Please enter a valid email address"),
  code: z.string().length(6, "Please enter the 6 digit code from your email."),
});
export type Form = z.infer<typeof FormSchema>;

type Props = {
  initialAccount?: Account;
  returnURL?: string;
};

export function EmailVerificationScreen(props: Props) {
  const router = useRouter();
  const probablyEmail = props.initialAccount?.email_addresses.find(
    (e) => e.verified === false,
  );

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      email: probablyEmail?.email_address,
    },
  });

  const handleSubmit = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        const r = await authEmailVerify({
          email: data.email,
          code: data.code,
        });

        router.push(props.returnURL ?? "/");
      },
      {
        errorToast: false,
        async onError(e) {
          form.setError("root", {
            message: "Invalid email or code.",
          });
        },
      },
    );
  });

  return (
    <form className={vstack()} onSubmit={handleSubmit}>
      <Heading>Verify your email address.</Heading>
      <p>Check your email for a 6 digit code.</p>

      <FormControl>
        <Input
          type="email"
          placeholder="Email address..."
          {...form.register("email")}
        />
        <FormErrorText>{form.formState.errors["email"]?.message}</FormErrorText>
      </FormControl>

      <FormControl>
        <PinInputField control={form.control} name="code" length={6} />
        <FormErrorText>{form.formState.errors["code"]?.message}</FormErrorText>
      </FormControl>

      <Button
        w="full"
        loading={form.formState.isSubmitting}
        disabled={!form.formState.isValid}
      >
        Verify
      </Button>

      <FormErrorText>{form.formState.errors["root"]?.message}</FormErrorText>

      {props.returnURL && (
        <Link className="link" href={props.returnURL}>
          Back to previous page
        </Link>
      )}
    </form>
  );
}
