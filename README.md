# Broadcast and Discovery
Let your servers discover each other on the network

## Description
Application Server (TCP) is started on 8080. This server handles all persistent
connections.

A "Presence" Message is broadcasted to all networks (255.255.255.255) on port
9999.  After that, a Discovery Server (UDP) is started on port 9999, listening
on IPv4zero (`0.0.0.0`) to receive "Presence" message from peers.

Due to Discovery message being sent out, any other instance of this application
running on the same network will reach the message and can connect to this
instance as a Peer via TCP for a persistent connection.

## TODO: Add Gossip
Currently there is only Broadcast and Discovery of instances running in the same
network.

Gossip needs to happen to allow instances to share info to one another.
