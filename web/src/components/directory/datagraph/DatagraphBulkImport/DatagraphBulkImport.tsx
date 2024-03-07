import { Button } from "src/theme/components/Button";
import { Card, CardRows } from "src/theme/components/Card";
import { Input } from "src/theme/components/Input";

import { ClusterCard } from "../ClusterCard";

import { css } from "@/styled-system/css";
import { HStack, LStack, styled } from "@/styled-system/jsx";

import { Props, State, useDatagraphBulkImport } from "./useDatagraphBulkImport";

export function DatagraphBulkImport(props: Props) {
  const { directoryPath, form, handlers, data } = useDatagraphBulkImport(props);
  const placeholder = "Enter URL(s)";
  const isDirty = data.items.length > 0;

  return (
    <LStack>
      <styled.form display="flex" w="full" onSubmit={handlers.handleSubmission}>
        <Input
          w="full"
          borderRight="none"
          borderRightRadius="none"
          type="search"
          placeholder={placeholder}
          {...form.register("url")}
        />

        {isDirty && (
          <Button borderRadius="none" onClick={handlers.handleClear}>
            Clear
          </Button>
        )}
        <Button
          flexShrink="0"
          borderLeft="none"
          borderLeftRadius="none"
          type="submit"
          // TODO: Fix the disabled button styling, this looks goofy af
          // disabled={data.urls.length === 0}
          width="min"
        >
          {data.urls.length > 0 ? (
            <>
              Prepare {data.urls.length}{" "}
              {data.urls.length === 1 ? "link" : "links"} for import
            </>
          ) : (
            <>Prepare for import</>
          )}
        </Button>
      </styled.form>

      <CardRows>
        {data.items.map((task) => {
          const { url } = task;

          switch (task.state) {
            case "idle":
              return (
                <Card
                  key={url}
                  id={url}
                  shape="row"
                  title="idle"
                  url=""
                  controls={<StateBadge state={task.state} />}
                />
              );

            case "loading":
              return (
                <Card
                  key={url}
                  id={url}
                  shape="row"
                  title="loading"
                  url=""
                  controls={<StateBadge state={task.state} />}
                />
              );

            case "success": {
              const { link } = task;
              return (
                <Card
                  key={url}
                  id={link.slug}
                  shape="row"
                  title={link.title || link.url}
                  text={link.description || "(no description found)"}
                  url={link.url}
                  image={link.assets[0]?.url}
                  controls={<StateBadge state={task.state} />}
                >
                  <HStack>
                    <Button
                      size="xs"
                      kind="primary"
                      onClick={() => handlers.handleCreateNodeFromLink(link)}
                    >
                      {props.node ? `Create in ${props.node?.name}` : "Create"}
                    </Button>
                  </HStack>
                </Card>
              );
            }

            case "created":
              return (
                <ClusterCard
                  key={task.node.id}
                  directoryPath={directoryPath}
                  context="directory"
                  cluster={task.node}
                  shape="row"
                />
              );

            case "error":
              return (
                <Card
                  key={url}
                  id={url}
                  shape="row"
                  title="Unable to fetch link"
                  text={task.error}
                  url=""
                  controls={
                    <HStack>
                      <Button
                        size="xs"
                        onClick={() => handlers.handleRemove(url)}
                      >
                        Remove
                      </Button>
                      <StateBadge state={task.state} />
                    </HStack>
                  }
                />
              );
          }
        })}
      </CardRows>
    </LStack>
  );
}

const badge = css({
  fontSize: "xs",
  fontWeight: "bold",
  color: "fg.default",
  px: "1",
  py: "0.5",
  borderRadius: "sm",
});

function StateBadge({ state }: { state: State }) {
  switch (state) {
    case "idle":
      return (
        <styled.div className={badge} bgColor="gray.100">
          Idle
        </styled.div>
      );

    case "loading":
      return (
        <styled.div className={badge} bgColor="blue.400">
          Loading
        </styled.div>
      );

    case "success":
      return (
        <styled.div className={badge} bgColor="green.500">
          Success
        </styled.div>
      );

    case "created":
      return (
        <styled.div className={badge} bgColor="green.500">
          Created
        </styled.div>
      );

    case "error":
      return (
        <styled.div className={badge} bgColor="red.500">
          Error
        </styled.div>
      );
  }
}
