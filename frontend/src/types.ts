export interface RecordedTransaction {
  IsReqText: boolean;
  IsResText: boolean;
  StatusCode: number;
  Ulid: string;
  ResBody: string;
  ResHeader: { [key: string]: string[] };
  ReqBody: string;
  Url: string;
  ReqHeader: { [key: string]: string[] };
  Method: string;
}

export interface ExecutionResponse {
  IsSameResBody: true;
  IsSameResHeader: true;
  IsSameStatusCode: true;
  ActualResHeader: { [key: string]: string[] };
  ActualResBody: string;
  IsBodyText: boolean;
  StatusCode: number;
}
