#include <iostream>
#include <thread>
#include <chrono>
#include "log_receiver.h"
#include "zero_copy_queue.h"
#include "redis_publisher.h"
#include "config.h"

int main(int argc, char* argv[]) {
    std::string configPath = log_receiver::ConfigLoader::GetConfigPath();
    log_receiver::Config config = log_receiver::ConfigLoader::LoadConfig(configPath);

    const int udp_port = config.cpp_receiver.udp_port;
    const int tcp_port = config.cpp_receiver.tcp_port;
    const size_t queue_size = config.cpp_receiver.queue_size;
    const std::string redis_host = config.cpp_receiver.redis_host;
    const int redis_port = config.cpp_receiver.redis_port;
    const std::string redis_stream = config.cpp_receiver.redis_stream;

    log_receiver::ZeroCopyQueue queue(queue_size);
    log_receiver::RedisPublisher publisher(redis_host, redis_port, redis_stream);

    if (!publisher.Connect()) {
        std::cerr << "Failed to connect to Redis at " << redis_host << ":" << redis_port << std::endl;
        return 1;
    }

    auto callback = [&queue](const std::string& data, const std::string& source) {
        queue.Push(data.data(), data.size(), source);
    };

    log_receiver::LogReceiver receiver(udp_port, tcp_port, callback);
    receiver.Start();

    std::thread consumer([&queue, &publisher]() {
        log_receiver::ZeroCopyQueue::Chunk chunk;
        while (true) {
            if (queue.Pop(chunk)) {
                std::string data(chunk.data, chunk.size);
                publisher.Publish(data, chunk.source);
            } else {
                std::this_thread::sleep_for(std::chrono::microseconds(10));
            }
        }
    });

    std::cout << "Log receiver started. UDP: " << udp_port << ", TCP: " << tcp_port << std::endl;

    std::string line;
    std::getline(std::cin, line);

    receiver.Stop();
    consumer.join();

    return 0;
}