1. **Authentication:**
   - `AUTH_REQUEST`
     ```json
     {
       "action": "AUTH_REQUEST",
       "username": "user123",
       "password": "password123"
     }
     ```
   - `AUTH_RESPONSE`
     ```json
     {
       "action": "AUTH_RESPONSE",
       "success": true,
       "user_id": "12345"
     }
     ```

2. **User Presence:**
   - `USER_JOIN`
     ```json
     {
       "action": "USER_JOIN",
       "user_id": "12345",
       "timestamp": "2023-01-01T12:00:00Z"
     }
     ```
   - `USER_LEAVE`
     ```json
     {
       "action": "USER_LEAVE",
       "user_id": "12345",
       "timestamp": "2023-01-01T12:05:00Z"
     }
     ```
   - `USER_STATUS_UPDATE`
     ```json
     {
       "action": "USER_STATUS_UPDATE",
       "user_id": "12345",
       "status": "online",
       "timestamp": "2023-01-01T12:10:00Z"
     }
     ```

3. **Text Messaging:**
   - `TEXT_MESSAGE`
     ```json
     {
       "action": "TEXT_MESSAGE",
       "sender_id": "12345",
       "receiver_id": "67890",
       "content": "Hello!",
       "timestamp": "2023-01-01T12:15:00Z"
     }
     ```
   - `SYSTEM_MESSAGE`
     ```json
     {
       "action": "SYSTEM_MESSAGE",
       "content": "New user joined the chat.",
       "timestamp": "2023-01-01T12:20:00Z"
     }
     ```

4. **Message Acknowledgment:**
   - `MESSAGE_ACK`
     ```json
     {
       "action": "MESSAGE_ACK",
       "message_id": "abc123",
       "timestamp": "2023-01-01T12:25:00Z"
     }
     ```
   - `MESSAGE_NACK`
     ```json
     {
       "action": "MESSAGE_NACK",
       "message_id": "def456",
       "error": "Recipient not found",
       "timestamp": "2023-01-01T12:30:00Z"
     }
     ```

5. **File Transfer:**
   - `FILE_REQUEST`
     ```json
     {
       "action": "FILE_REQUEST",
       "sender_id": "12345",
       "receiver_id": "67890",
       "file_name": "document.pdf",
       "timestamp": "2023-01-01T12:35:00Z"
     }
     ```
   - `FILE_RESPONSE`
     ```json
     {
       "action": "FILE_RESPONSE",
       "file_id": "xyz789",
       "status": "approved",
       "timestamp": "2023-01-01T12:40:00Z"
     }
     ```
   - `FILE_TRANSFER`
     ```json
     {
       "action": "FILE_TRANSFER",
       "file_id": "xyz789",
       "file_data": "<binary_data>",
       "timestamp": "2023-01-01T12:45:00Z"
     }
     ```

6. **Error Handling:**
   - `ERROR`
     ```json
     {
       "action": "ERROR",
       "error_code": 500,
       "error_message": "Internal Server Error",
       "timestamp": "2023-01-01T12:50:00Z"
     }
     ```
   - `AUTH_ERROR`
     ```json
     {
       "action": "AUTH_ERROR",
       "error_code": 401,
       "error_message": "Invalid credentials",
       "timestamp": "2023-01-01T12:55:00Z"
     }
     ```

7. **Heartbeat/Ping:**
   - `PING`
     ```json
     {
       "action": "PING",
       "timestamp": "2023-01-01T13:00:00Z"
     }
     ```
   - `PONG`
     ```json
     {
       "action": "PONG",
       "timestamp": "2023-01-01T13:05:00Z"
     }
     ```

8. **Offline Message Handling:**
   - `OFFLINE_MESSAGE`
     ```json
     {
       "action": "OFFLINE_MESSAGE",
       "sender_id": "12345",
       "receiver_id": "67890",
       "content": "I missed you!",
       "timestamp": "2023-01-01T13:10:00Z"
     }
     ```

9. **Encryption:**
   - `ENCRYPTED_MESSAGE`
     ```json
     {
       "action": "ENCRYPTED_MESSAGE",
       "sender_id": "12345",
       "receiver_id": "67890",
       "encrypted_content": "<encrypted_data>",
       "timestamp": "2023-01-01T13:15:00Z"
     }
     ```

10. **User Interaction:**
    - `USER_TYPING`
      ```json
      {
        "action": "USER_TYPING",
        "sender_id": "12345",
        "receiver_id": "67890",
        "timestamp": "2023-01-01T13:20:00Z"
      }
      ```
    - `READ_RECEIPT`
      ```json
      {
        "action": "READ_RECEIPT",
        "message_id": "abc123",
        "timestamp": "2023-01-01T13:25:00Z"
      }
      ```

11. **Metadata:**
    - `METADATA_UPDATE`
      ```json
      {
        "action": "METADATA_UPDATE",
        "user_id": "12345",
        "metadata": {
          "nickname": "CoolUser",
          "avatar_url": "https://example.com/avatar.jpg"
        },
        "timestamp": "2023-01-01T13:30:00Z"
      }
      ```

12. **Protocol Management:**
    - `PROTOCOL_VERSION`
      ```json
      {
        "action": "PROTOCOL_VERSION",
        "version": "1.0",
        "timestamp": "2023-01-01T13:35:00Z"
      }
      ```

13. **User Blocking/Filtering:**
    - `BLOCK_USER`
      ```json
      {
        "action": "BLOCK_USER",
        "blocker_id": "12345",
        "blocked_id": "67890",
        "timestamp": "2023-01-01T13:40:00Z"
      }
      ```
    - `UNBLOCK_USER`
      ```json
      {
        "action": "UNBLOCK_USER",
        "blocker_id": "12345",
        "unblocked_id": "67890",
        "timestamp": "2023-01-01T13:45:00Z"
      }
      ```

14. **User Feedback:**
   

 - `FEEDBACK`
      ```json
      {
        "action": "FEEDBACK",
        "user_id": "12345",
        "feedback_type": "suggestion",
        "feedback_content": "Add more emoji options!",
        "timestamp": "2023-01-01T13:50:00Z"
      }
      ```

15. **Connection Management:**
    - `CONNECTION_ESTABLISHED`
      ```json
      {
        "action": "CONNECTION_ESTABLISHED",
        "timestamp": "2023-01-01T13:55:00Z"
      }
      ```
    - `CONNECTION_CLOSED`
      ```json
      {
        "action": "CONNECTION_CLOSED",
        "reason": "Client logout",
        "timestamp": "2023-01-01T14:00:00Z"
      }
      ```

16. **Server-to-Client Notifications:**
    - `SERVER_NOTIFICATION`
      ```json
      {
        "action": "SERVER_NOTIFICATION",
        "content": "Maintenance scheduled at 2 AM UTC",
        "timestamp": "2023-01-01T14:05:00Z"
      }
      ```

17. **User Information:**
    - `REQUEST_USER_INFO`
      ```json
      {
        "action": "REQUEST_USER_INFO",
        "requester_id": "12345",
        "target_id": "67890",
        "timestamp": "2023-01-01T14:10:00Z"
      }
      ```
    - `RESPONSE_USER_INFO`
      ```json
      {
        "action": "RESPONSE_USER_INFO",
        "user_id": "67890",
        "username": "user678",
        "nickname": "FriendlyUser",
        "avatar_url": "https://example.com/friendly_avatar.jpg",
        "timestamp": "2023-01-01T14:15:00Z"
      }
      ```

18. **Group Chat (if applicable):**
    - `GROUP_MESSAGE`
      ```json
      {
        "action": "GROUP_MESSAGE",
        "group_id": "group123",
        "sender_id": "12345",
        "content": "Hello group!",
        "timestamp": "2023-01-01T14:20:00Z"
      }
      ```
    - `GROUP_INVITE`
      ```json
      {
        "action": "GROUP_INVITE",
        "inviter_id": "12345",
        "invitee_id": "67890",
        "group_id": "group123",
        "timestamp": "2023-01-01T14:25:00Z"
      }
      ```