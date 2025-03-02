export interface Response<T> {
  statusCode: number;
  statusMessage: string;
  data: T;
}
