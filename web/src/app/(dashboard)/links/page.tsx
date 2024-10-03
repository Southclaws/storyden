import { LinkIndexScreen } from "src/screens/library/links/LinkIndexScreen/LinkIndexScreen";

type Props = {
  searchParams: {
    q: string;
    page: number;
  };
};

export default function Page(props: Props) {
  return (
    <LinkIndexScreen
      query={props.searchParams.q}
      page={props.searchParams.page}
    />
  );
}
