import { XMarkIcon } from "@heroicons/react/24/outline";

import { Unready } from "src/components/site/Unready";
import { Button } from "src/theme/components/Button";
import { Input } from "src/theme/components/Input";

import { MemberList } from "../MemberList";

import { VStack, styled } from "@/styled-system/jsx";

import { Props, useMemberIndexView } from "./useMemberIndexView";

export function MemberIndexView(props: Props) {
  const { form, data, handlers } = useMemberIndexView(props);

  if (form.formState.isLoading) return <Unready />;

  return (
    <VStack>
      <styled.form
        display="flex"
        w="full"
        onSubmit={handlers.handleSubmission}
        action="/p"
      >
        <Input
          w="full"
          borderRight="none"
          borderRightRadius="none"
          type="search"
          placeholder="Search for members"
          defaultValue={props.query}
          {...form.register("q")}
        />

        {(props.query || data.q) && (
          <Button
            borderX="none"
            borderRadius="none"
            type="reset"
            onClick={handlers.handleReset}
          >
            <XMarkIcon />
          </Button>
        )}
        <Button
          flexShrink="0"
          borderLeft="none"
          borderLeftRadius="none"
          type="submit"
          width="min"
        >
          Search
        </Button>
      </styled.form>

      <MemberList profiles={data.results.profiles} />
    </VStack>
  );
}
