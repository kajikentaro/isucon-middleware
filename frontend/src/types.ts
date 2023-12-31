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
  ReqLength: number;
  ResLength: number;
}

export interface ExecutionResponse {
  IsSameResBody: true;
  IsSameResHeader: true;
  IsSameStatusCode: true;
  ActualResHeader: Header;
  ActualResBody: string;
  IsBodyText: boolean;
  StatusCode: number;
  ActualResLength: number;
}

export interface IsRecording {
  IsRecording: boolean;
}

export type Header = { [key: string]: string[] };
