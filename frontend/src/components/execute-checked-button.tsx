import { useExecuteChecked } from "@/hooks/use-execute-checked";

export default function ExecuteCheckedButton() {
  const executeChecked = useExecuteChecked();

  return (
    <button
      className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-2 rounded-full flex items-center"
      onClick={(e) => {
        executeChecked();
        e.stopPropagation();
      }}
    >
      Execute Checked
    </button>
  );
}
