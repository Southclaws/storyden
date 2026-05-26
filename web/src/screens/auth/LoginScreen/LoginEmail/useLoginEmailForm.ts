"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import * as z from "zod";

import { handle } from "@/api/client";
import { useAccountGet } from "@/api/openapi-client/accounts";
import {
  authEmailPasswordSignin,
  authPasswordSignin,
} from "@/api/openapi-client/auth";
import {
  ExistingPasswordSchema,
  UsernameOrEmailSchema,
} from "@/lib/auth/schemas";
import { refreshFeed } from "@/lib/feed/refresh";

const FormSchema = z.object({
  identifier: UsernameOrEmailSchema,
  password: ExistingPasswordSchema,
});
type Form = z.infer<typeof FormSchema>;

export function useLoginEmailForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const returnURL = searchParams.get("return_url") ?? "/";

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const { mutate } = useAccountGet();

  const handleSubmit = form.handleSubmit(async (payload: Form) => {
    await handle(async () => {
      const isEmail = payload.identifier.includes("@");

      if (isEmail) {
        await authEmailPasswordSignin({
          email: payload.identifier,
          password: payload.password,
        });
      } else {
        await authPasswordSignin({
          identifier: payload.identifier,
          token: payload.password,
        });
      }

      mutate();
      refreshFeed();
      router.push(returnURL);
    });
  });

  return {
    form,
    handlers: {
      handleSubmit,
    },
  };
}
