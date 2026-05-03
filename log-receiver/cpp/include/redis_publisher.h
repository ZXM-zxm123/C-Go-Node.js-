#pragma once

#include <memory>
#include <string>

namespace log_receiver {

class RedisPublisher {
public:
    RedisPublisher(const std::string& host, int port, const std::string& stream);
    ~RedisPublisher();

    bool connect();
    void disconnect();
    bool publish(const std::string& data, const std::string& source);

private:
    struct Impl;
    std::unique_ptr<Impl> impl_;
};

} // namespace log_receiver