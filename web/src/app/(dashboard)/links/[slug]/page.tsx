import { LinkScreen } from "src/screens/library/links/LinkScreen/LinkScreen";

type Props = {
  params: Promise<{
    slug: string;
  }>;
};

export default async function Page(props: Props) {
  return <LinkScreen slug={(await props.params).slug} />;
}
