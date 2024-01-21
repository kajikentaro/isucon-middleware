import { useExecute } from "@/hooks/use-execute";
import Code from "@/parts/code";
import { TagBinary } from "@/parts/tag-binary";
import { useAppDispatch, useAppSelector } from "@/store";
import { selectExecutionResponse } from "@/store/execution-response";
import { selectRecordedTransaction } from "@/store/recorded-transaction";
import {
  closeComparisonPopup,
  selectComparisonPopup,
} from "@/store/ui/comparison-popup";
import { Header } from "@/types";
import { BodyType } from "@/utils/get-url";
import { stringifyHeader } from "@/utils/stringify-header";
import { useEffect } from "react";

export default function ComparisonPopup() {
  const popupState = useAppSelector(selectComparisonPopup);

  if (!popupState.isVisible) {
    return null;
  }

  // split main content to avoid conditional call of useEffect
  return <ComparisonPopupMain />;
}

function ComparisonPopupMain() {
  const dispatch = useAppDispatch();

  useEffect(() => {
    window.addEventListener("keydown", closePopupOnEscapePressed);
    return () => {
      window.removeEventListener("keydown", closePopupOnEscapePressed);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const closePopup = () => {
    dispatch(closeComparisonPopup());
  };

  const closePopupOnEscapePressed = (event: KeyboardEvent) => {
    if (event.key === "Escape") {
      closePopup();
    }
  };

  return (
    <div
      className="fixed z-50 top-0 left-0 w-full h-full flex justify-center items-center bg-black bg-opacity-50 px-10 py-20"
      onClick={closePopup}
    >
      <div
        className="bg-white p-6 rounded-md w-full h-full overflow-y-auto"
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
        <button
          onClick={closePopup}
          className="absolute top-30 right-20 text-gray-800 text-xl w-12 h-12 bg-slate-200 hover:opacity-80 rounded-full"
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
  const onExecute = useExecute(popupState.ulid);

  return (
    <div>
      <Request
        statusCode={recordedTransaction.StatusCode}
        ulid={popupState.ulid}
        body={recordedTransaction.ReqBody}
        header={recordedTransaction.ReqHeader}
        isText={recordedTransaction.IsReqText}
        contentLength={recordedTransaction.ReqLength}
      />
      <span className="mb-4 h-0.5 bg-gray-300 block" />
      <div className="flex justify-center">
        <div className="w-full flex">
          <Response
            transaction={{
              statusCode: recordedTransaction.StatusCode,
              ulid: popupState.ulid,
              body: recordedTransaction.ResBody,
              header: recordedTransaction.ResHeader,
              isText: recordedTransaction.IsResText,
              contentLength: recordedTransaction.ResLength,
            }}
            type="res-body"
            title="Recorded Response"
          />
          <span className="w-0.5 bg-gray-300" />
          {executionResponse ? (
            <Response
              transaction={{
                statusCode: executionResponse.StatusCode,
                ulid: popupState.ulid,
                body: executionResponse.ActualResBody,
                header: executionResponse.ActualResHeader,
                isText: executionResponse.IsBodyText,
                contentLength: executionResponse.ActualResLength,
              }}
              type="reproduced-res-body"
              title="Actual Response"
            />
          ) : (
            <div className="w-1/2 p-4 rounded-md mb-4">
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
    </div>
  );
}

interface TransactionProps {
  statusCode: number;
  ulid: string;
  body: string;
  header: Header;
  isText: boolean;
  contentLength: number;
}

function Request(props: TransactionProps) {
  const { header, isText, body, ulid, contentLength } = props;
  return (
    <div className="p-4 rounded-md">
      <h3 className="text-lg font-semibold my-2"></h3>
      <p>Request Header:</p>
      <Code>{stringifyHeader(header)}</Code>
      <p>Request Body:</p>
      {isText ? (
        <Code>{body}</Code>
      ) : (
        <TagBinary
          ulid={ulid}
          type="req-body"
          className="mt-2"
          contentLength={contentLength}
        />
      )}
    </div>
  );
}

function Response(props: {
  transaction: TransactionProps;
  title: string;
  type: BodyType;
}) {
  const {
    type,
    title,
    transaction: { header, statusCode, isText, body, ulid, contentLength },
  } = props;
  return (
    <div className="w-1/2 p-4 rounded-md mb-4">
      <h3 className="text-lg font-semibold my-2">{title}</h3>
      <p>Response Header:</p>
      <Code>{stringifyHeader(header)}</Code>
      <p>Status Code:</p>
      <Code>{statusCode}</Code>
      <p>Response Body:</p>
      {isText ? (
        <Code>{body}</Code>
      ) : (
        <TagBinary
          ulid={ulid}
          type={type}
          className="mt-2"
          contentLength={contentLength}
        />
      )}
    </div>
  );
}
