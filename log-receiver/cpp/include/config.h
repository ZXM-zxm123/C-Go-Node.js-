#pragma once

#include <string>
#include <iostream>
#include <cstdlib>

namespace log_receiver {

struct CppReceiverConfig {
    int udp_port;
    int tcp_port;
    int queue_size;
    std::string redis_host;
    int redis_port;
    std::string redis_stream;
};

struct Config {
    CppReceiverConfig cpp_receiver;
};

class ConfigLoader {
public:
    static Config LoadConfig(const std::string& path);
    static void PrintError(const std::string& file, int line, 
                          const std::string& type, const std::string& message);
    static std::string GetConfigPath();

private:
    static void ValidateConfig(const Config& config);
    static std::string GetAbsolutePath(const std::string& path);
    static bool FileExists(const std::string& path);
    static void Trim(std::string& s);
    static std::string GetErrorMessage(int lineNum, const std::string& type);
};

} // namespace log_receiver