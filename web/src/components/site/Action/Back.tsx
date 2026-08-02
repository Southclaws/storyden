"use client";

import React, { PropsWithChildren } from "react";

import { Button, ButtonProps } from "@/components/ui/button";
import { ArrowLeftIcon } from "@/components/ui/icons/Arrow";
import { useSmartBack } from "@/lib/navigation/smart-back";

type Props = PropsWithChildren<
  ButtonProps & {
    fallbackHref?: string;
  }
>;

export function BackAction({ children, fallbackHref = "/", ...props }: Props) {
  const goBack = useSmartBack(fallbackHref);
  const hasLabel = React.Children.count(children) > 0;
  const label = props["aria-label"] ?? (hasLabel ? undefined : "Back");

  return (
    <Button
      variant="ghost"
      px={hasLabel ? undefined : "0"}
      aria-label={label}
      title={props.title ?? label}
      onClick={goBack}
      {...props}
    >
      <ArrowLeftIcon /> {children}
    </Button>
  );
}
