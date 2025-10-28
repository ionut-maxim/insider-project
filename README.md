# insider-project

## Development setup
This project uses `mise` for environment setup and as a task runner. 
You can install it with this [guide](https://mise.jdx.dev/getting-started.html).  
Mise has extensions for **Goland** and **VS Code** if needed to ensure environment variables are properly loaded.

`Docker` must also be installed on the system.

After you installed mise:
```shell
mise trust
```
Followed by
```shell
mise install
```
At this point all the tools required to build and generate code for the project are installed.

> [!NOTE]
> As a side note you can run the server without `mise` installed also, in this case you will need to manually set the environment variables and have running instances of postgres and redis.
### Starting the server
Now you can run the following command to start the server. It will pull a postgres and redis container and the compile and run the server
```shell
mise run start
```
Once the server starts you can visit http://localhost:8080/docs to start interacting with the API.

### Code generation
For code generation you can run `mise run generate` and any changes made to the `Typespec` file will reflect in a new openapi file and a code generated for the server interface in Golang.

### Environment
The `.env` file already contains sensible defaults for a local development instance.  

## Design/Architecture
### API
The API is simple with 4 actions that can be taken:
 - Worker Start
 - Worker Stop
 - Get sent messages
 - Add a new message to send (this endpoint was added for easier testing) 

I've chosen a schema based approach for the api.  
It has a simple flow: Define schema in Typespec `openapi/main.tsp` -> `openapi.yaml` emmited from Typespec project -> Golang server interface generated from `openapi.yaml`

A handler was added with `swagger-ui` for docs.

Package for the server, that also glues the rest of the system can be found in `internal/server`

### Worker/Scheduler
A simple worker was created to handle polling of the database table for unsent messages. 
Worker package is located here `internal/worker`.

### Message Store
Because it's a shared dependency I've chosen to define the interface for it at the root of the project.

### Main
Main entrypoint in the application: `cmd/server`. It handles configuration parsing and package instancing.
