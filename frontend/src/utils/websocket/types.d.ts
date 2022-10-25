interface IResponse<T> {
  statusCode: number;
  payload: T;
}

type EventHandler = (payload: any) => void;

interface IEventHandlerBase {
  id?: number;
  event: string;
  resolve?: (response: IResponse<any>) => void;
  reject?: (response: IResponse<any>) => void;
  callback?: EventHandler;
}

interface IEventHandler extends IEventHandlerBase {
  id: number;
  requestID: string;
}
