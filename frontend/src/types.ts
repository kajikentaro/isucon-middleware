export interface SearchResponse {
  transactions: RecordedTransaction[];
  totalHit: number;
}

export interface RecordedTransaction {
  isReqText: boolean;
  isResText: boolean;
  statusCode: number;
  ulid: string;
  resBody: string;
  resHeader: Header;
  reqBody: string;
  url: string;
  reqHeader: Header;
  method: string;
  reqLength: number;
  resLength: number;
}

export interface ExecutionResponse {
  isSameResBody: true;
  isSameResHeader: true;
  isSameStatusCode: true;
  actualResHeader: Header;
  actualResBody: string;
  isBodyText: boolean;
  statusCode: number;
  actualResLength: number;
}

export interface IsRecording {
  isRecording: boolean;
}

export interface TotalTransactions {
  count: number;
}

export type Header = { [key: string]: string[] };
