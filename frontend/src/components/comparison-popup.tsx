import { useExecute } from "@/hooks/use-execute";
import Code from "@/parts/code";
import { ExecuteButton } from "@/parts/execute-button";
import Modal from "@/parts/modal";
import ProgressIcon from "@/parts/progress-icon";
import { TagBinary } from "@/parts/tag-binary";
import { useAppDispatch, useAppSelector } from "@/store";
import { selectExecutionProgress } from "@/store/execution-progress";
import { selectExecutionResponse } from "@/store/execution-response";
import { selectRecordedTransaction } from "@/store/recorded-transaction";
import {
  closeComparisonPopup,
  selectComparisonPopup,
} from "@/store/ui/comparison-popup";
import { Header } from "@/types";
import { shouldBeNever } from "@/utils/assert-never";
import { BodyType } from "@/utils/get-url";
import { stringifyHeader } from "@/utils/stringify-header";

export default function ComparisonPopup() {
  const dispatch = useAppDispatch();

  const closePopup = () => {
    dispatch(closeComparisonPopup());
  };

  const { isVisible } = useAppSelector(selectComparisonPopup);

  return (
    <Modal isVisible={isVisible} closePopup={closePopup} title={"Comparison"}>
      <ModalContents />
    </Modal>
  );
}

function ModalContents() {
  const popupState = useAppSelector(selectComparisonPopup);
  const recordedTransaction = useAppSelector(
    selectRecordedTransaction(popupState.ulid)
  );

  return (
    <div>
      <Request
        statusCode={recordedTransaction.statusCode}
        ulid={popupState.ulid}
        body={recordedTransaction.reqBody}
        header={recordedTransaction.reqHeader}
        isText={recordedTransaction.isReqText}
        contentLength={recordedTransaction.reqLength}
      />
      <span className="mb-4 h-0.5 bg-gray-300 block" />
      <div className="flex justify-center">
        <div className="w-full flex">
          <Response
            transaction={{
              statusCode: recordedTransaction.statusCode,
              ulid: popupState.ulid,
              body: recordedTransaction.resBody,
              header: recordedTransaction.resHeader,
              isText: recordedTransaction.isResText,
              contentLength: recordedTransaction.resLength,
            }}
            type="res-body"
            title={
              <h3 className="text-lg font-semibold my-2">Recorded Response</h3>
            }
          />
          <span className="w-0.5 bg-gray-300" />
          <ActualResponse />
        </div>
      </div>
    </div>
  );
}

function ActualResponse() {
  const popupState = useAppSelector(selectComparisonPopup);
  const executionResponse = useAppSelector(
    selectExecutionResponse(popupState.ulid)
  );
  const executionProgress = useAppSelector(
    selectExecutionProgress(popupState.ulid)
  );

  const onExecute = useExecute(popupState.ulid);

  const onExecuteClick: React.MouseEventHandler<HTMLButtonElement> = (e) => {
    e.stopPropagation();
    onExecute();
  };

  switch (executionProgress) {
    case "init":
      return (
        <div className="w-1/2 p-4 rounded-md mb-4">
          <p>This transaction have not been executed yet</p>
          <div className="mt-10 flex justify-center">
            <ExecuteButton onClick={onExecuteClick}>Execute</ExecuteButton>
          </div>
        </div>
      );
    case "fail":
    case "waitingResponse":
    case "waitingQueue":
      return (
        <div className="w-1/2 p-4 rounded-md mb-4">
          <p>Executing</p>
          <div className="py-2 px-2 mt-10 flex justify-center">
            <ProgressIcon ulid={popupState.ulid} />
          </div>
        </div>
      );
    case "bodyNotSame":
    case "headerNotSame":
    case "statusCodeNotSame":
    case "success":
      return (
        <Response
          transaction={{
            statusCode: executionResponse.statusCode,
            ulid: popupState.ulid,
            body: executionResponse.actualResBody,
            header: executionResponse.actualResHeader,
            isText: executionResponse.isBodyText,
            contentLength: executionResponse.actualResLength,
          }}
          type="reproduced-res-body"
          title={
            <div className="flex flex-row gap-3 items-center">
              <h3 className="text-lg font-semibold my-2">Actual Response</h3>
              <ProgressIcon ulid={popupState.ulid} />
              <div className="ml-auto">
                <ExecuteButton onClick={onExecuteClick}>
                  Execute Again
                </ExecuteButton>
              </div>
            </div>
          }
        />
      );
    default:
      shouldBeNever(executionProgress);
  }
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
      <Code className="my-2">{stringifyHeader(header)}</Code>
      <p>Request Body:</p>
      {isText ? (
        <Code className="my-2">{body}</Code>
      ) : (
        <TagBinary
          ulid={ulid}
          type="req-body"
          className="my-2"
          contentLength={contentLength}
        />
      )}
    </div>
  );
}

function Response(props: {
  transaction: TransactionProps;
  title: JSX.Element;
  type: BodyType;
}) {
  const {
    type,
    title,
    transaction: { header, statusCode, isText, body, ulid, contentLength },
  } = props;
  return (
    <div className="w-1/2 p-4 rounded-md mb-4">
      {title}
      <p>Response Header:</p>
      <Code className="my-2">{stringifyHeader(header)}</Code>
      <p>Status Code:</p>
      <Code className="my-2">{statusCode}</Code>
      <p>Response Body:</p>
      {isText ? (
        <Code className="my-2">{body}</Code>
      ) : (
        <TagBinary
          ulid={ulid}
          type={type}
          className="my-2"
          contentLength={contentLength}
        />
      )}
    </div>
  );
}
