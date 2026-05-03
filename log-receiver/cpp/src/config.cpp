#include "config.h"

#include <fstream>
#include <sstream>
#include <algorithm>
#include <cctype>
#include <vector>
#include <sys/stat.h>

#ifndef _WIN32
#include <unistd.h>
#endif

namespace log_receiver {

std::string ConfigLoader::GetConfigPath() {
    const char* envPath = std::getenv("CONFIG_PATH");
    if (envPath != nullptr) {
        return std::string(envPath);
    }
    return "../config/config.yaml";
}

bool ConfigLoader::FileExists(const std::string& path) {
    struct stat buffer;
    return (stat(path.c_str(), &buffer) == 0 && S_ISREG(buffer.st_mode));
}

void ConfigLoader::Trim(std::string& s) {
    s.erase(s.begin(), std::find_if(s.begin(), s.end(), [](unsigned char ch) {
        return !std::isspace(ch);
    }));
    s.erase(std::find_if(s.rbegin(), s.rend(), [](unsigned char ch) {
        return !std::isspace(ch);
    }).base(), s.end());
}

void ConfigLoader::PrintError(const std::string& file, int line, 
                              const std::string& type, const std::string& message) {
    std::cerr << "ERROR: " << type << " in " << file;
    if (line > 0) {
        std::cerr << " (line " << line << ")";
    }
    std::cerr << ": " << message << std::endl;
}

void ConfigLoader::ValidateConfig(const Config& config) {
    std::vector<std::string> errors;

    if (config.cpp_receiver.udp_port <= 0 || config.cpp_receiver.udp_port > 65535) {
        errors.push_back("cpp_receiver.udp_port must be between 1 and 65535, got " + 
                        std::to_string(config.cpp_receiver.udp_port));
    }
    if (config.cpp_receiver.tcp_port <= 0 || config.cpp_receiver.tcp_port > 65535) {
        errors.push_back("cpp_receiver.tcp_port must be between 1 and 65535, got " + 
                        std::to_string(config.cpp_receiver.tcp_port));
    }
    if (config.cpp_receiver.queue_size <= 0) {
        errors.push_back("cpp_receiver.queue_size must be positive, got " + 
                        std::to_string(config.cpp_receiver.queue_size));
    }
    if (config.cpp_receiver.redis_host.empty()) {
        errors.push_back("cpp_receiver.redis_host is required");
    }
    if (config.cpp_receiver.redis_port <= 0 || config.cpp_receiver.redis_port > 65535) {
        errors.push_back("cpp_receiver.redis_port must be between 1 and 65535, got " + 
                        std::to_string(config.cpp_receiver.redis_port));
    }
    if (config.cpp_receiver.redis_stream.empty()) {
        errors.push_back("cpp_receiver.redis_stream is required");
    }

    if (!errors.empty()) {
        std::cerr << "ERROR: configuration validation failed:" << std::endl;
        for (const auto& err : errors) {
            std::cerr << "  - " << err << std::endl;
        }
        std::exit(EXIT_FAILURE);
    }
}

std::string ConfigLoader::GetAbsolutePath(const std::string& path) {
#ifdef _WIN32
    char full[MAX_PATH];
    if (_fullpath(full, path.c_str(), MAX_PATH) == nullptr) {
        return path;
    }
    return std::string(full);
#else
    char buffer[PATH_MAX];
    if (realpath(path.c_str(), buffer) == nullptr) {
        return path;
    }
    return std::string(buffer);
#endif
}

Config ConfigLoader::LoadConfig(const std::string& path) {
    std::string absPath = GetAbsolutePath(path);
    
    if (!FileExists(absPath)) {
        PrintError(absPath, 0, "configuration file not found", 
                   "the file does not exist or is not a regular file");
        std::exit(EXIT_FAILURE);
    }

    std::ifstream file(absPath);
    if (!file.is_open()) {
        PrintError(absPath, 0, "failed to read config file", 
                   "could not open the file for reading");
        std::exit(EXIT_FAILURE);
    }

    Config config;
    std::string line;
    int lineNum = 0;
    int indentLevel = 0;
    std::string currentSection;
    bool inCppReceiver = false;

    while (std::getline(file, line)) {
        lineNum++;
        
        std::string originalLine = line;
        Trim(line);
        
        if (line.empty() || line[0] == '#') {
            continue;
        }

        int currentIndent = 0;
        for (char ch : originalLine) {
            if (std::isspace(ch)) {
                currentIndent++;
            } else {
                break;
            }
        }

        if (currentIndent == 0) {
            if (line.back() == ':') {
                std::string section = line.substr(0, line.size() - 1);
                Trim(section);
                inCppReceiver = (section == "cpp_receiver");
            } else {
                PrintError(absPath, lineNum, "YAML syntax error", 
                          "expected a section header ending with ':'");
                std::exit(EXIT_FAILURE);
            }
        } else if (inCppReceiver && currentIndent == 2) {
            size_t colonPos = line.find(':');
            if (colonPos == std::string::npos) {
                PrintError(absPath, lineNum, "YAML syntax error", 
                          "expected a key-value pair");
                std::exit(EXIT_FAILURE);
            }

            std::string key = line.substr(0, colonPos);
            std::string value = line.substr(colonPos + 1);
            Trim(key);
            Trim(value);

            if (value.empty() && !key.empty()) {
                continue;
            }

            try {
                if (key == "udp_port") {
                    config.cpp_receiver.udp_port = std::stoi(value);
                } else if (key == "tcp_port") {
                    config.cpp_receiver.tcp_port = std::stoi(value);
                } else if (key == "queue_size") {
                    config.cpp_receiver.queue_size = std::stoi(value);
                } else if (key == "redis_host") {
                    if (value.size() >= 2 && value.front() == '"' && value.back() == '"') {
                        config.cpp_receiver.redis_host = value.substr(1, value.size() - 2);
                    } else {
                        config.cpp_receiver.redis_host = value;
                    }
                } else if (key == "redis_port") {
                    config.cpp_receiver.redis_port = std::stoi(value);
                } else if (key == "redis_stream") {
                    if (value.size() >= 2 && value.front() == '"' && value.back() == '"') {
                        config.cpp_receiver.redis_stream = value.substr(1, value.size() - 2);
                    } else {
                        config.cpp_receiver.redis_stream = value;
                    }
                } else if (!key.empty()) {
                    PrintError(absPath, lineNum, "unknown configuration field", 
                              "unexpected key: " + key);
                    std::exit(EXIT_FAILURE);
                }
            } catch (const std::invalid_argument& e) {
                PrintError(absPath, lineNum, "type mismatch error", 
                          "invalid numeric value for key '" + key + "'");
                std::exit(EXIT_FAILURE);
            } catch (const std::out_of_range& e) {
                PrintError(absPath, lineNum, "type mismatch error", 
                          "numeric value out of range for key '" + key + "'");
                std::exit(EXIT_FAILURE);
            }
        } else if (inCppReceiver && currentIndent != 0) {
            PrintError(absPath, lineNum, "YAML syntax error", 
                      "unexpected indentation level");
            std::exit(EXIT_FAILURE);
        }
    }

    ValidateConfig(config);
    return config;
}

} // namespace log_receiver