import { useCollectionGet } from "src/api/openapi/collections";

export type Props = {
  handle: string;
  collection: string;
};

export function useCollectionScreen(props: Props) {
  const { data, error, isLoading } = useCollectionGet(props.collection);

  return { data, error, isLoading };
}
