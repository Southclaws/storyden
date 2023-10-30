interface Props {
  method: string | undefined;
}
export function AuthMethod({ method }: Props) {
  return <p>{method}</p>;
}
