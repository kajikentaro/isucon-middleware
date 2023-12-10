
export interface FetchResponse {
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