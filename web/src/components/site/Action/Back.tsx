"use client";

import { useRouter } from "next/navigation";
import React, { PropsWithChildren } from "react";

import { Button, ButtonProps } from "@/components/ui/button";
import { ArrowLeftIcon } from "@/components/ui/icons/Arrow";

export function BackAction({
  children,
  ...props
}: PropsWithChildren<ButtonProps>) {
  const router = useRouter();

  const hasLabel = React.Children.count(children) > 0;

  function handleBack() {
    router.back();
  }

  return (
    <Button
      size="xs"
      variant="ghost"
      px={hasLabel ? undefined : "0"}
      onClick={handleBack}
      {...props}
    >
      <ArrowLeftIcon /> {children}
    </Button>
  );
}
