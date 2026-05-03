#pragma once

#include <functional>
#include <string>

namespace log_receiver {

class LogReceiver {
public:
    using LogCallback = std::function<void(const std::string&, const std::string&)>;

    LogReceiver(int udp_port, int tcp_port, LogCallback callback);
    ~LogReceiver();

    void start();
    void stop();

private:
    struct Impl;
    std::unique_ptr<Impl> impl_;
};

} // namespace log_receiver