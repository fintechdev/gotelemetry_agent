## URL-based Agent demo

In this directory, you will find an example of running Agent jobs over HTTP. When the job for `test_flow` is executed, the Agent makes a request over HTTP and uses the returned information to patch the flow.

If the HTTP request returns a value greater than 399, the Agent bubbles the returned data to the API as a job error that the system will display in your account's log.

In this example, the web functionality is provided by a simple Node.js-based server running locally, but, of course, you could query any server, as long as it is accessible from the machine on which the Agent runs, and it returns the right information.