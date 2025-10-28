import "@typespec/http";

using Http;
@service(#{ title: "Notification Service" })
@server("/api/v1", "API v1")
namespace NotificationService;

@error
model Error {
  message: string;
}

model WorkerResponse {
  message: string;
}

@route("/worker")
@tag("Worker")
@summary("Handles worker lifecycle")
interface Worker {
  @route("start")
  @get
  @summary("Starts the polling worker")
  start(): WorkerResponse | Error;

  @route("stop")
  @get
  @summary("Stops the polling worker")
  stop(): WorkerResponse | Error;
}



model Message {
  @visibility(Lifecycle.Read) id: integer;
  @summary("Phone number of the recipient")
  @example("+40700000000")
  to: string;
  @summary("Message content")
  @example("How are you on this beautiful day?")
  content: string;
  @visibility(Lifecycle.Read) createdAt: utcDateTime;
  @visibility(Lifecycle.Read) updatedAt: utcDateTime;
}

@route("/messages")
@tag("Messages")
interface Messages {
  @get sent(@query limit: integer, @query offset: integer): Message[] | Error;

  @post add(...Message): CreatedResponse | Error;
}
