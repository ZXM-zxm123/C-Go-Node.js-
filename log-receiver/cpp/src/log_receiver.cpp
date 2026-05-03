#include "log_receiver.h"
#include <iostream>
#include <thread>
#include <atomic>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <cstring>

namespace log_receiver {

struct LogReceiver::Impl {
    int udp_port;
    int tcp_port;
    LogCallback callback;
    std::thread udp_thread;
    std::thread tcp_thread;
    std::atomic<bool> running{false};
    int udp_socket{-1};
    int tcp_socket{-1};
};

LogReceiver::LogReceiver(int udp_port, int tcp_port, LogCallback callback)
    : impl_(std::make_unique<Impl>()) {
    impl_->udp_port = udp_port;
    impl_->tcp_port = tcp_port;
    impl_->callback = std::move(callback);
}

LogReceiver::~LogReceiver() {
    stop();
}

void udp_listener(int socket, LogReceiver::LogCallback callback) {
    char buffer[65536];
    struct sockaddr_in client_addr;
    socklen_t addr_len = sizeof(client_addr);

    while (true) {
        ssize_t n = recvfrom(socket, buffer, sizeof(buffer) - 1, 0,
                             reinterpret_cast<sockaddr*>(&client_addr), &addr_len);
        if (n <= 0) continue;
        
        buffer[n] = '\0';
        std::string source = inet_ntoa(client_addr.sin_addr);
        callback(std::string(buffer, n), source);
    }
}

void tcp_handler(int client_socket, LogReceiver::LogCallback callback, const std::string& source) {
    char buffer[65536];
    std::string line;
    
    while (true) {
        ssize_t n = read(client_socket, buffer, sizeof(buffer) - 1);
        if (n <= 0) break;
        
        buffer[n] = '\0';
        line += buffer;
        
        size_t pos;
        while ((pos = line.find('\n')) != std::string::npos) {
            std::string log_line = line.substr(0, pos);
            line = line.substr(pos + 1);
            if (!log_line.empty()) {
                callback(log_line, source);
            }
        }
    }
    
    close(client_socket);
}

void tcp_listener(int socket, LogReceiver::LogCallback callback) {
    struct sockaddr_in client_addr;
    socklen_t addr_len = sizeof(client_addr);
    
    while (true) {
        int client_socket = accept(socket, reinterpret_cast<sockaddr*>(&client_addr), &addr_len);
        if (client_socket < 0) continue;
        
        std::string source = inet_ntoa(client_addr.sin_addr);
        std::thread(tcp_handler, client_socket, callback, source).detach();
    }
}

void LogReceiver::start() {
    impl_->running = true;

    impl_->udp_socket = socket(AF_INET, SOCK_DGRAM, 0);
    struct sockaddr_in addr;
    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(impl_->udp_port);
    
    int opt = 1;
    setsockopt(impl_->udp_socket, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));
    
    bind(impl_->udp_socket, reinterpret_cast<sockaddr*>(&addr), sizeof(addr));
    
    impl_->udp_thread = std::thread(udp_listener, impl_->udp_socket, impl_->callback);

    impl_->tcp_socket = socket(AF_INET, SOCK_STREAM, 0);
    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = INADDR_ANY;
    addr.sin_port = htons(impl_->tcp_port);
    
    setsockopt(impl_->tcp_socket, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));
    bind(impl_->tcp_socket, reinterpret_cast<sockaddr*>(&addr), sizeof(addr));
    listen(impl_->tcp_socket, 100);
    
    impl_->tcp_thread = std::thread(tcp_listener, impl_->tcp_socket, impl_->callback);
}

void LogReceiver::stop() {
    impl_->running = false;
    
    if (impl_->udp_socket >= 0) {
        close(impl_->udp_socket);
        impl_->udp_socket = -1;
    }
    
    if (impl_->tcp_socket >= 0) {
        close(impl_->tcp_socket);
        impl_->tcp_socket = -1;
    }
    
    if (impl_->udp_thread.joinable()) {
        impl_->udp_thread.join();
    }
    
    if (impl_->tcp_thread.joinable()) {
        impl_->tcp_thread.join();
    }
}

} // namespace log_receiver