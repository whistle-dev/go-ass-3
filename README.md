# go-ass-3

### Description
You have to implement Chitty-Chat a distributed system, that is providing a chatting service, and keeps track of logical time using Lamport Timestamps.

We call clients of the Chitty-Chat service Participants. 

## Requirements
### System Requirements

1. Chitty-Chat is a distributed service, that enables its clients to chat. The service is using gRPC for communication. You have to design the API, including gRPC methods and data types.  Discuss, whether you are going to use server-side streaming, client-side streaming, or bidirectional streaming?
2. Clients in Chitty-Chat can Publish a valid chat message at any time they wish.  A valid message is a string of UTF-8 encoded text with a maximum length of 128 characters. A client publishes a message by making a gRPC call to Chitty-Chat.
3. The Chitty-Chat service has to broadcast every published message, together with the current Lamport timestamp, to all participants in the system, by using gRPC. It is an implementation decision left to the students, whether a Vector Clock or a Lamport timestamp is sent.
4. When a client receives a broadcasted message, it has to write the message and the current Lamport timestamp to the log
5. Chat clients can join at any time. 
6. A "Participant X  joined Chitty-Chat at Lamport time L" message is broadcast to all Participants when client X joins, including the new Participant.
7. Chat clients can drop out at any time. 
8. A "Participant X left Chitty-Chat at Lamport time L" message is broadcast to all remaining Participants when Participant X leaves.


### Technical Requirements:

1. Use gRPC for all messages passing between nodes
2. Use Golang to implement the service and clients
3. Every client has to be deployed as a separate process
4. Log all service calls (Publish, Broadcast, ...) using the log package
5. Demonstrate that the system can be started with at least 3 client nodes 
6. Demonstrate that a client node can join the system
7. Demonstrate that a client node can leave the system
8. Optional: All elements of the Chitty-Chat service are deployed as Docker containers


### Hand-in requirements:

1. Hand in a single report in a pdf file
2. Describe your system architecture - do you have a server-client architecture, peer-to-peer, or something else?
3. Describe what  RPC methods are implemented, of what type, and what messages types are used for communication
4. Describe how you have implemented the calculation of the Lamport timestamps
5. Provide a diagram, that traces a sequence of RPC calls together with the Lamport timestamps, that corresponds to a chosen sequence of interactions: Client X joins, 6. Client X Publishes, ..., Client X leaves. Include documentation (system logs) in your appendix.
7. Provide a link to a Git repo with your source code in the report
8. Include system logs, that document the requirements are met, in the appendix of your report
