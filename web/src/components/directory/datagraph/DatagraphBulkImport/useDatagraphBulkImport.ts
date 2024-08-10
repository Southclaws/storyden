import { zodResolver } from "@hookform/resolvers/zod";
import { fromPairs } from "lodash";
import { entries, sortBy } from "lodash/fp";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { linkCreate } from "src/api/openapi-client/links";
import { Link } from "src/api/openapi-schema";
import { DatagraphNodeWithRelations } from "src/components/directory/datagraph/DatagraphNode";
import { useDirectoryPath } from "src/screens/directory/datagraph/useDirectoryPath";
import { deriveError } from "src/utils/error";

// We want to pull many URLs from the input and process them in parallel.
const multipleURLs =
  /\b((https?:\/\/?|www[.])[^\s()<>]+(?:\([\w\d]+\)|([^[:punct:]\s]|\/?)))/gi;

export type Props = {
  node?: DatagraphNodeWithRelations;
  onCreateNodeFromLink: (link: Link) => Promise<DatagraphNodeWithRelations>;
};

const ManyLinkSchema = z
  .string()
  .optional()
  .transform((s: string | undefined) => {
    if (!s) return [];

    const matches = s.match(multipleURLs);

    return matches?.map((url) => url.trim()) ?? [];
  });

const FormSchema = z.object({
  url: ManyLinkSchema,
});
type Form = z.infer<typeof FormSchema>;

type TaskState =
  | { state: "idle" }
  | { state: "loading" }
  | { state: "success"; link: Link }
  | { state: "created"; link: Link; node: DatagraphNodeWithRelations }
  | { state: "error"; error: string };

type TaskStateWithUrl = TaskState & { url: string };

type TaskTable = Record<string, TaskState>;

export type State = TaskState["state"];

const TaskStateSortOrder: Record<TaskState["state"], number> = {
  created: 0,
  success: 1,
  loading: 2,
  idle: 3,
  error: 4,
};

const sortTasks = sortBy<TaskStateWithUrl>((a) => {
  return TaskStateSortOrder[a.state];
});

export function useDatagraphBulkImport(props: Props) {
  const [tasks, setTasks] = useState<TaskTable>({});
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
  });
  const directoryPath = useDirectoryPath();

  const taskList = entries(tasks).map(([url, state]) => ({ url, ...state }));
  const items = sortTasks(taskList);

  const { url } = form.watch();
  const urls = ManyLinkSchema.parse(url);

  const urlToNode = fromPairs(
    props.node?.children.map((node) => [node.link?.url, node]),
  );

  async function handleCreateNodeFromLink(link: Link) {
    const node = await props.onCreateNodeFromLink(link);

    setTasks((current) => ({
      ...current,
      [link.url]: { state: "created", link, node },
    }));
  }

  const handleSubmission = form.handleSubmit((data) => {
    Promise.all(
      data.url.map(async (url) => {
        setTasks((current) => ({ ...current, [url]: { state: "loading" } }));

        try {
          const link = await linkCreate({ url });

          const duplicate = urlToNode[link.url];

          if (duplicate) {
            setTasks((current) => ({
              ...current,
              [url]: {
                state: "created",
                link,
                node: {
                  ...duplicate,
                  type: "node",
                } as DatagraphNodeWithRelations,
              },
            }));
          } else {
            setTasks((current) => ({
              ...current,
              [url]: { state: "success", link },
            }));
          }
        } catch (e: unknown) {
          console.error(e);
          setTasks((current) => ({
            ...current,
            [url]: { state: "error", error: deriveError(e) },
          }));
        }
      }),
    );

    form.reset();
  });

  function handleRemove(url: string) {
    setTasks((current) => {
      const next = { ...current };
      delete next[url];
      return next;
    });
  }

  function handleClear() {
    setTasks({});
    form.reset();
  }

  return {
    directoryPath,
    form,
    data: {
      urls,
      items,
    },
    handlers: {
      handleSubmission,
      handleCreateNodeFromLink,
      handleRemove,
      handleClear,
    },
  };
}
