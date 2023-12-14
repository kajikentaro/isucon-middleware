import { useOnExecute } from "@/hooks/use-execute";
import {
  closeComparisonPopup,
  selectComparisonPopup,
} from "@/store/comparison-popup";
import { selectExecutionResponse } from "@/store/execution-response";
import { useAppDispatch, useAppSelector } from "@/store/main";
import { selectRecordedTransaction } from "@/store/recorded-transaction";
import { stringifyHeader } from "@/utils/stringify-header";

export default function ComparisonPopup() {
  const popupState = useAppSelector(selectComparisonPopup);
  const dispatch = useAppDispatch();
  if (!popupState.isVisible) {
    return null;
  }

  const onClose = () => {
    dispatch(closeComparisonPopup());
  };

  return (
    <div
      className="fixed z-50 top-0 left-0 w-full h-full flex justify-center items-center bg-black bg-opacity-50 px-10 py-20"
      onClick={onClose}
    >
      <div
        className="bg-white p-6 rounded-md w-full h-full relative"
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
        <button
          onClick={onClose}
          className="absolute top-2 right-2 text-gray-600 hover:text-gray-800 text-xl w-12 h-12 bg-slate-100 rounded-full"
        >
          X
        </button>
        <h2 className="text-2xl font-bold mb-4">Comparison</h2>
        <ModalContents />
      </div>
    </div>
  );
}

function ModalContents() {
  const popupState = useAppSelector(selectComparisonPopup);
  const recordedTransaction = useAppSelector(
    selectRecordedTransaction(popupState.ulid)
  );
  const executionResponse = useAppSelector(
    selectExecutionResponse(popupState.ulid)
  );
  const onExecute = useOnExecute(popupState.ulid);

  return (
    <div className="flex justify-center">
      <div className="w-full flex">
        <Transaction
          header={recordedTransaction.ResHeader}
          statusCode={recordedTransaction.StatusCode}
          body={recordedTransaction.ResBody}
          title="Recorded Response"
        />
        {executionResponse ? (
          <Transaction
            header={executionResponse.ActualResHeader}
            statusCode={executionResponse.StatusCode}
            body={executionResponse.ActualResBody}
            title="Actual Response"
          />
        ) : (
          <div className="w-1/2 border border-gray-300 p-4 rounded-md mb-4">
            <p>This transaction have not been executed yet</p>
            <button
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-2 rounded-full mt-10 m-auto block"
              onClick={(e) => {
                e.stopPropagation();
                onExecute();
              }}
            >
              Execute
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

interface TransactionProps {
  header: { [key: string]: string[] };
  statusCode: number;
  body: string;
  title: string;
}
function Transaction(props: TransactionProps) {
  const { header, statusCode, body, title } = props;
  return (
    <div className="w-1/2 border border-gray-300 p-4 rounded-md mb-4">
      <h3 className="text-lg font-semibold mb-2">{title}</h3>
      <p>Response Header:</p>
      <code className="block bg-black text-white p-2 rounded-md mb-2">
        {stringifyHeader(header)}
      </code>
      <p>Status Code:</p>
      <code className="block bg-black text-white p-2 rounded-md mb-2">
        {statusCode}
      </code>
      <p>Response Body:</p>
      <code className="block bg-black text-white p-2 rounded-md mb-2">
        {body}
      </code>
    </div>
  );
}
