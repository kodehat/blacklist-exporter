# Blacklist Exporter for Prometheus

This application exports blacklist statuses to Prometheus, using the MXToolBox API.

## Overview

The exporter makes API requests to MXToolBox to check the blacklist status of specified IP addresses, and then populates Prometheus metrics based on the results. It runs on **port 2112** by default.

## Features

- Uses MXToolBox API to check blacklist status
- Takes API token and IPs from `.env` configuration
- Starts a Prometheus server and populates metrics from MXToolBox responses

## Blacklist Status Categories

The MXToolBox API returns blacklist statuses in the following categories:

- **Passed**: The IP address or domain is not listed on the specific blacklist.
- **Failed**: The IP address or domain is listed on the specific blacklist.
- **Warnings**: There are potential issues or configurations that might lead to future blacklisting or email delivery problems.
- **Timeouts**: The check could not be completed due to a timeout, possibly because the blacklist server did not respond in time.

These categories help identify the status and potential issues affecting email deliverability.

## Setup

1. Clone the repository.
2. Set up the `.env` file with your MXToolBox API token and the IPs to check.
3. Run the exporter to start the Prometheus server on **port 2112**.

## Docker Usage

The exporter can also be run via Docker:

```bash
docker build -t image:tag .
docker run -p 2112:2112 image:tag
