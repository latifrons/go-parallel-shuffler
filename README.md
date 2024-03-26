## GoParallelShuffler

GoParallelShuffler is a Go library that provides a solution for executing disordered messages in parallel and then re-ordering the output by their message IDs.

In various distributed systems and message queues, messages can arrive out of order due to network conditions, load balancing, or other factors. GoMessageShuffler allows you to process these disordered messages concurrently, taking advantage of multi-core CPUs, and then outputs the results in the correct order based on their message IDs.

**Key Features:**

- **Parallel Execution**: Utilizes workers to process multiple messages simultaneously, improving overall throughput.
- **Disordered Input**: Handles messages that arrive in a random or shuffled order.
- **Ordered Output**: Guarantees that the output results are sorted by the message IDs, restoring the order by message ID.
- **Efficient**: Employs buffered channels and worker pools to maximize CPU utilization and minimize resource contention.
- **Customizable**: Allows you to define your own message processing logic while benefiting from the concurrent execution and ordering capabilities.

GoParallelShuffler can be useful in various scenarios, such as distributed data processing, message-based architectures, and real-time event handling systems, where maintaining the correct order of results is essential despite receiving disordered input.

## Example

Let's say we have a set of messages with numerical IDs that need to be processed:

Input Messages (Disordered):

[5] Message 5

[2] Message 2

[4] Message 4

[1] Message 1

[3] Message 3

Using GoParallelShuffler, we can process these messages concurrently and obtain the results in the correct order:

Output Results (Ordered by ID):

[1] Processed Message 1

[2] Processed Message 2

[3] Processed Message 3

[4] Processed Message 4

[5] Processed Message 5
