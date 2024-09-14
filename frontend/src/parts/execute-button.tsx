interface Props {
  children: JSX.Element | string;
  onClick: React.MouseEventHandler<HTMLButtonElement>;
}

export function ExecuteButton({ children, onClick }: Props) {
  return (
    <button
      className="p-2 rounded-full flex items-center text-white bg-blue-500 hover:bg-blue-700 font-bold"
      onClick={onClick}
    >
      {children}
    </button>
  );
}
