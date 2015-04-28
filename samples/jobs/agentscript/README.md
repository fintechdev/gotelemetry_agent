## AgentScript and Graphite sample

The examples in this directory show you how you can use the Agent's Graphite integration and AgentScript to manage complex job that support timeseries data and counters.

In the `config.yaml` file, the `graphite` entry instructs the agent to listen for Graphite messages sent to it both over UDP or TCP. Since it supports Graphite's text protocol natively, you can use the Agent as a drop-in replacement and immediately start collecting and aggregating data through it.

The `set_test_flow.asl` file gives you a taste of what AgentScript can do; in this example, data from a timeseries is aggregate over a series of intervals to populate a value flow with both a numeric value and a sparkline. The Agent's anomaly-detection functionality is also used to turn the text in the flow red when a new value is inconsistent with the data that has been used to populate the timeseries thus far.