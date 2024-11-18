import { PasswordCreateForm } from "./PasswordCreateForm";
import { PasswordUpdateForm } from "./PasswordUpdateForm";

export type Props = {
  active: boolean;
};

export function Password(props: Props) {
  return props.active ? <PasswordUpdateForm /> : <PasswordCreateForm />;
}
