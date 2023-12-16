export interface RecordedTransaction {
  IsReqText: boolean;
  IsResText: boolean;
  StatusCode: number;
  Ulid: string;
  ResBody: string;
  ResHeader: Header;
  ReqBody: string;
  Url: string;
  ReqHeader: Header;
  Method: string;
}

export interface ExecutionResponse {
  IsSameResBody: true;
  IsSameResHeader: true;
  IsSameStatusCode: true;
  ActualResHeader: Header;
  ActualResBody: string;
  IsBodyText: boolean;
  StatusCode: number;
}

export type Header = { [key: string]: string[] };
