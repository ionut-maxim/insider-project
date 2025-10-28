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
  @maxLength(400)
  content: string;
  @visibility(Lifecycle.Read) createdAt: utcDateTime;
  @visibility(Lifecycle.Read) updatedAt: utcDateTime;
}

@route("/messages")
@tag("Messages")
interface Messages {
  @summary("Gets all the sent messages")
  @get sent(@query limit?: integer = 10, @query offset?: integer): Message[] | Error;

  @summary("Add a new message. This was added for testing purposes.")
  @post add(...Message): CreatedResponse | Error;
}
