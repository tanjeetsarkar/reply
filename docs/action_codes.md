# Action Codes

## Authentication
**AUTH_REQUEST**: Client requests authentication.  
**AUTH_RESPONSE**: Server responds to authentication request.

## User Presence
**USER_JOIN**: User joins the chat.  
**USER_LEAVE**: User leaves the chat.  
**USER_STATUS_UPDATE**: User updates their online/offline status.

## Text Messaging
**TEXT_MESSAGE**: Regular text message.  
**SYSTEM_MESSAGE**: System-generated message (e.g., notifications, alerts).

## Message Acknowledgment
**MESSAGE_ACK**: Acknowledgment of successful message delivery.  
**MESSAGE_NACK**: Notification of failed message delivery.

## File Transfer
**FILE_REQUEST**: Client requests file transfer.  
**FILE_RESPONSE**: Server responds to file transfer request.  
**FILE_TRANSFER**: Actual file transfer.

## Error Handling
**ERROR**: General error message.  
**AUTH_ERROR**: Authentication-related error.

## Heartbeat/Ping
**PING**: Client sends a ping to the server.  
**PONG**: Server responds to the client's ping.

## Offline Message Handling
**OFFLINE_MESSAGE**: Message sent while the recipient was offline.

## Encryption
**ENCRYPTED_MESSAGE**: Indication that the message content is encrypted.

## User Interaction
**USER_TYPING**: Notification that a user is typing.  
**READ_RECEIPT**: Notification that a message has been read.

## Metadata
**METADATA_UPDATE**: Update to additional metadata fields.

## Protocol Management
**PROTOCOL_VERSION**: Exchange of protocol version information.

## User Blocking/Filtering
**BLOCK_USER**: Client blocks another user.  
**UNBLOCK_USER**: Client unblocks a previously blocked user.

## User Feedback
**FEEDBACK**: User provides feedback on the application.

## Connection Management
**CONNECTION_ESTABLISHED**: Notification that a new connection has been established.  
**CONNECTION_CLOSED**: Notification that a connection has been closed.

## Server-to-Client Notifications
**SERVER_NOTIFICATION**: General notification from the server to the client.

## User Information
**REQUEST_USER_INFO**: Client requests information about another user.  
**RESPONSE_USER_INFO**: Server responds with user information.

## Group Chat (if applicable)
**GROUP_MESSAGE**: Message sent in a group chat.  
**GROUP_INVITE**: Invitation to join a group.