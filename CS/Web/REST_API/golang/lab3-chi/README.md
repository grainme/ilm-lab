## Lab 2: Learning Notes

> A summary of key concepts and techniques from the lab.

### 1. What is Redis?
`Remote Dictionary Server` is an open-source, in-memory data structure store. It's often called a "data structure server" because its core data types are similar to programming language data structures:

*   **Strings:** The most basic type, can hold any kind of data (text, numbers, binary data).
*   **Lists:** A sequence of strings, ordered by insertion. Think of it like a doubly-linked list.
*   **Sets:** An unordered collection of unique strings.
*   **Hashes:** A map between string fields and string values, perfect for storing objects.
*   **Sorted Sets:** Similar to Sets, but each member has an associated score, which is used to keep the set sorted.

### 2. Common Use Cases

Because Redis is so fast (being in-memory), it's excellent for:

*   **Caching:** This is the most common use case. You can store the results of expensive database queries or API calls in Redis. When the same data is requested again, you can fetch it from Redis instead of hitting the database, which is much faster.
*   **Session Management:** Storing user session data for web applications.
*   **Real-time Analytics:** Counting and tracking real-time events.
*   **Leaderboards:** Sorted Sets are perfect for maintaining leaderboards.
*   **Message Queues:** Lists can be used to implement a simple message queue.
