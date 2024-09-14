import { useExecuteChecked } from "@/hooks/use-execute-checked";
import { ExecuteButton } from "@/parts/execute-button";

export default function ExecuteCheckedButton() {
  const executeChecked = useExecuteChecked();

  return (
    <ExecuteButton
      onClick={(e) => {
        executeChecked();
        e.stopPropagation();
      }}
    >
      Execute Checked
    </ExecuteButton>
  );
}
