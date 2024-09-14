import { ReactNode, useEffect } from "react";

interface Props {
  children: ReactNode;
  isVisible: boolean;
  closePopup: () => void;
  title: string;
}

export default function Modal(props: Props) {
  if (!props.isVisible) {
    return null;
  }

  return <ModalContent {...props} />;
}

function ModalContent({ children, closePopup, title }: Props) {
  const closePopupOnEscapePressed = (event: KeyboardEvent) => {
    if (event.key === "Escape") {
      closePopup();
    }
  };

  useEffect(() => {
    window.addEventListener("keydown", closePopupOnEscapePressed);
    return () => {
      window.removeEventListener("keydown", closePopupOnEscapePressed);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div
      className="fixed z-50 top-0 left-0 w-full h-full flex justify-center items-start bg-black bg-opacity-50 px-10 py-20"
      onClick={closePopup}
    >
      <div className="relative h-full">
        <div
          className="bg-white p-6 rounded-md max-w-full max-h-full overflow-y-auto"
          onClick={(e) => {
            e.stopPropagation();
          }}
        >
          <button
            onClick={closePopup}
            className="absolute top-5 right-5 text-gray-800 text-xl w-12 h-12 bg-slate-200 hover:opacity-80 rounded-full"
          >
            X
          </button>
          <h2 className="text-2xl font-bold mt-2">{title}</h2>
          {children}
        </div>
      </div>
    </div>
  );
}
