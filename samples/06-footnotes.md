# Footnotes and Definition Lists

## Footnotes

Kubernetes[^1] is the de facto standard for container orchestration. OpenTelemetry[^2] provides vendor-neutral observability.

[^1]: Originally developed by Google, now maintained by the CNCF.
[^2]: Formed by merging OpenTracing and OpenCensus in 2019.

## Definition Lists

Kubernetes
:   An open-source container orchestration platform for automating deployment, scaling, and management.

OpenTelemetry
:   A collection of tools, APIs, and SDKs for generating, collecting, and exporting telemetry data.
:   Supports traces, metrics, and logs.

Jaeger
:   An open-source distributed tracing system, originally developed by Uber.

## Combined

The observability stack[^3] typically consists of three pillars.

Traces
:   End-to-end request flow across services.

Metrics
:   Numeric measurements over time (counters, histograms).

Logs
:   Timestamped text records of discrete events.

[^3]: As defined by the OpenTelemetry project.
