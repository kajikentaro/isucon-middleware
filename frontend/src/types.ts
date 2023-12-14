export interface RecordedTransaction {
  IsReqText: boolean;
  IsResText: boolean;
  StatusCode: number;
  Ulid: string;
  ResBody: string;
  ResHeader: { [key: string]: string[] };
  ReqBody: string;
  ReqOthers: {
    Url: string;
    Header: { [key: string]: string[] };
    Method: string;
  };
}

export interface ExecutionResponse {
  ActualResBody: string;
  ActualResHeader: { [key: string]: string[] };
  IsSameResBody: true;
  IsSameResHeader: true;
  IsSameStatusCode: true;
  IsBodyText: boolean;
  StatusCode: number;
}
