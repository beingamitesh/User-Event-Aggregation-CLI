# User-Event-Aggregation-CLI
A command-line utility for a social media app backend, aggregating user activity events to generate daily summary reports. Reports encompass post counts, received likes, etc., updated in real-time with new event data.

## Overview

This Go application processes user events stored in a JSON file and aggregates daily summaries of events based on user IDs and event types. The application supports updating existing summaries with new events, ensuring that each event is considered only once.

## Features

- **Input/Output**: The application takes input from a JSON file containing user events and produces daily summaries in another JSON file.
- **Update Existing Summaries**: An optional flag `--update` allows updating existing summaries with new events.
- **Event Deduplication**: The application checks for duplicate events using a hash mechanism to avoid processing the same event multiple times.

## Usage

### Command-line Arguments

- `-i` (Required): Path to the input JSON file containing user events.
- `-o` (Required): Path to the output JSON file where daily summaries will be stored.
- `--update` (Optional): If provided, updates existing summaries with new events.

### Example

```bash
go build .
./aggregate_events -i input.json -o output.json --update
