## JavaScript Agent demo

This directory contains a simple agent job that calls a JS script. As you can see, all the script needs to do is output a valid payload that is applied to the existing flow content using a PATCH operation.

In practice, this means that all the top-level properties your script outputs will completely overwrite the corresponding top-level properties in the Telemetry data flow associated with your job.