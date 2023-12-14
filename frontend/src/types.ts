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
  // TODO: add is res text
  ActualResBody: string;
  ActualResHeader: { [key: string]: string[] };
  IsSameResBody: true;
  IsSameResHeader: true;
  IsSameStatusCode: true;
}
