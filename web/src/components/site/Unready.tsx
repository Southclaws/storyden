"use client";

import { PropsWithChildren } from "react";

import { APIError } from "src/api/openapi-schema";

import { Skeleton } from "@/components/ui/skeleton";
import { Flex } from "@/styled-system/jsx";

import ErrorBanner from "./ErrorBanner";

export function Unready(props: PropsWithChildren<Partial<APIError>>) {
  if (!props.error) {
    return (
      <Flex
        flexDirection="column"
        width="full"
        justifyContent="center"
        p="4"
        gap="4"
      >
        {props.children ?? (
          <>
            <Skeleton />
            <Skeleton />
            <Skeleton />
          </>
        )}
      </Flex>
    );
  }

  return <ErrorBanner {...props} />;
}
