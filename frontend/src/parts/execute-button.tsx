interface Props {
  children: JSX.Element | string;
  onClick: React.MouseEventHandler<HTMLButtonElement>;
}

export function ExecuteButton({ children, onClick }: Props) {
  return (
    <button
      className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-2 rounded-full block m-auto"
      onClick={onClick}
    >
      {children}
    </button>
  );
}
