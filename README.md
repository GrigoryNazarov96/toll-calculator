# Toll Calculator Microservice project

## Overview

Microservice structured app built to aggregate data comes from a simulated OBU device.

Each OBU sends coordinates along with its ID via Web Socket protocol to the receiver.

Receiver sends the data further via Kafka transport straight to the distance calculator.

After calculating distance GRPC aggregation server being invoked, then it sends protobuf data to GRPC client.

Finally API Gateway is available for the end user to fetch an invoice from the inmemory storage, using a regular HTTP protocol.
