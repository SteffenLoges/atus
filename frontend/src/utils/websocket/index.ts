// event-driven, promise-based websocket libary
// 09/23/21 Steffen Loges

let conn: WebSocket;

const decoder = new TextDecoder("utf-8");
const encoder = new TextEncoder();
const messageRegExp = new RegExp(
  "^(?<event>[A-Za-z0-9-_]+?):?(?<requestID>[A-Za-z0-9]{15})?:(?<statusCode>[0-9]{3})?:?(?<payload>.+)?$"
);

// eventHandlers stores all persistent event handlers
let eventHandlers: IEventHandler[] = [];

// temporaryEventHandlers are used to store event handlers for a single request using send()
// the handler is removed after we received the response
let temporaryEventHandlers: {
  [requestID: string]: IEventHandlerBase;
} = {};

// native websocket event handlers
let _onOpen: WebSocket["onopen"] = () => {};
export const onOpen = (event: WebSocket["onopen"]) =>
  (_onOpen = event);

let _onMessage: WebSocket["onmessage"] = () => {};
export const onMessage = (event: WebSocket["onmessage"]) =>
  (_onMessage = event);

let _onClose: WebSocket["onclose"] = () => {};
export const onClose = (event: WebSocket["onclose"]) =>
  (_onClose = event);

let _onError: WebSocket["onerror"] = () => {};
export const onError = (event: WebSocket["onerror"]) =>
  (_onError = event);

export function connect({
  url = "",
  reconnect = true,
  reconnectTimeout = 1e3,
} = {}) {
  conn = new WebSocket(url);
  conn.binaryType = "arraybuffer";

  // -- setup native handlers -------------------

  conn.onmessage = function (event) {
    const decoded = decoder.decode(
      new Uint8Array(event.data)
    );

    for (const msg of decoded.split("\n")) {
      const matches = msg.match(messageRegExp);

      if (!matches) {
        throw new Error(
          `websocket: invalid message format: ${decoded}`
        );
      }

      const { event, requestID }: any = matches.groups;
      const statusCode = parseInt(
        matches.groups?.statusCode || "200"
      );

      let payload = matches.groups?.payload || "null";
      try {
        payload = JSON.parse(payload);
      } catch (e) {
        console.warn("websocket: invalid payload", payload);
        continue;
      }

      const statusIsOK = statusCode === 200;
      let handlers: IEventHandlerBase[] | IEventHandler[];
      // check if we have a temporary event handler for this request
      if (requestID && temporaryEventHandlers[requestID]) {
        handlers = [temporaryEventHandlers[requestID]];
        // remove from temporaryEventHandlers
        delete temporaryEventHandlers[requestID];
      }
      // check if we have a persistent event handler for this event
      else {
        handlers = eventHandlers.filter(
          (h) => h.event === event
        );
      }

      if (handlers.length === 0) {
        console.warn(
          `websocket: no handler for event: ${event}`
        );
      }

      for (const handler of handlers) {
        const resp: IResponse<any> = {
          statusCode,
          payload,
        };

        if (handler.callback) {
          handler.callback(resp);
        }

        if (handler.resolve && statusIsOK) {
          handler.resolve(resp);
        }

        if (handler.reject && !statusIsOK) {
          handler.reject(resp);
        }
      }
    }

    _onMessage?.call(this, event);
  };

  conn.onopen = _onOpen;

  conn.onerror = function (event) {
    console.error("websocket: error", event);
    _onError?.call(this, event);
    close();
  };

  conn.onclose = function (event) {
    console.log("websocket: closed", event);

    _onClose?.call(this, event);

    // reconnect if connection was closed unexpectedly and reconnect is enabled
    if (!event.wasClean && reconnect) {
      setTimeout(
        () => connect({ url, reconnect, reconnectTimeout }),
        reconnectTimeout
      );
    }
  };
}

export function addEventHandler(
  event: string,
  callback: EventHandler
) {
  const id = getLastEventHandlerID() + 1;

  eventHandlers.push({
    id,
    event,
    callback,
  } as IEventHandler);

  return id;
}

export function removeEventHandler(id: number) {
  const handlerIndex = eventHandlers.findIndex(
    (h) => h.id === id
  );
  if (handlerIndex != -1) {
    eventHandlers.splice(handlerIndex, 1);
  }
}

export function removeAllEventHandlers(
  event?: string
): void {
  if (!event) {
    eventHandlers = [];
    return;
  }

  eventHandlers = eventHandlers.filter(
    (h) => h.event !== event
  );
}

function getLastEventHandlerID() {
  // return eventHandlers.at(-1)?.id || 0;
  return eventHandlers.length > 0
    ? eventHandlers[eventHandlers.length - 1].id
    : 0;
}

export function close() {
  conn.close();
}

// genRequestID generates a random request ID for a websocket message
// is used to identify the response to a request
function genRequestID() {
  return (
    Date.now().toString(36) +
    Math.random().toString(36).substring(2)
  ).substring(0, 15);
}

// send sends a message to the websocket server
// returns a promise that resolves when the server responds
export function send(
  event: string,
  payload?: any
): Promise<any> {
  if (!conn) {
    throw new Error(
      "cannot send message, websocket not connected"
    );
  }

  let p = payload || "";
  // stringify payload if it is not already a string
  if (typeof p !== "string") {
    p = JSON.stringify(p);
  }

  const requestID = genRequestID();

  let msg = `${event}:${requestID}`;
  if (p !== "") {
    msg += `:${p}`;
  }

  conn.send(encoder.encode(msg));

  return new Promise((resolve, reject) => {
    temporaryEventHandlers[requestID] = {
      event,
      resolve: resolve,
      reject: reject,
    };
  });
}

export default {
  connect,
  onOpen,
  onMessage,
  onClose,
  onError,
  addEventHandler,
  removeEventHandler,
  removeAllEventHandlers,
  close,
  send,
};
