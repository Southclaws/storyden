import { LinkIndexScreen } from "src/screens/library/links/LinkIndexScreen/LinkIndexScreen";

type Props = {
  searchParams: Promise<{
    q: string;
    page: number;
  }>;
};

export default async function Page(props: Props) {
  return (
    <LinkIndexScreen
      query={(await props.searchParams).q}
      page={(await props.searchParams).page}
    />
  );
}
