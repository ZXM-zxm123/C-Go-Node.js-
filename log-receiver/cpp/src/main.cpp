#include <iostream>
#include <thread>
#include <chrono>
#include "log_receiver.h"
#include "zero_copy_queue.h"
#include "redis_publisher.h"

int main(int argc, char* argv[]) {
    const int udp_port = 514;
    const int tcp_port = 515;
    const size_t queue_size = 100000;
    const std::string redis_host = "localhost";
    const int redis_port = 6379;
    const std::string redis_stream = "log_stream";

    log_receiver::ZeroCopyQueue queue(queue_size);
    log_receiver::RedisPublisher publisher(redis_host, redis_port, redis_stream);

    if (!publisher.connect()) {
        std::cerr << "Failed to connect to Redis" << std::endl;
        return 1;
    }

    auto callback = [&queue](const std::string& data, const std::string& source) {
        queue.push(data.data(), data.size(), source);
    };

    log_receiver::LogReceiver receiver(udp_port, tcp_port, callback);
    receiver.start();

    std::thread consumer([&queue, &publisher]() {
        log_receiver::ZeroCopyQueue::Chunk chunk;
        while (true) {
            if (queue.pop(chunk)) {
                std::string data(chunk.data, chunk.size);
                publisher.publish(data, chunk.source);
            } else {
                std::this_thread::sleep_for(std::chrono::microseconds(10));
            }
        }
    });

    std::cout << "Log receiver started. UDP: " << udp_port << ", TCP: " << tcp_port << std::endl;

    std::string line;
    std::getline(std::cin, line);

    receiver.stop();
    consumer.join();

    return 0;
}