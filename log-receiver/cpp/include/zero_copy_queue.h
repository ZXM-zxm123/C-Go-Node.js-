#pragma once

#include <atomic>
#include <cstring>
#include <string>
#include <vector>

namespace log_receiver {

class ZeroCopyQueue {
public:
    struct Chunk {
        char* data;
        size_t size;
        size_t capacity;
        uint64_t timestamp;
        std::string source;
    };

    explicit ZeroCopyQueue(size_t capacity);
    ~ZeroCopyQueue();

    bool push(const char* data, size_t size, const std::string& source);
    bool pop(Chunk& chunk);
    size_t size() const;
    bool empty() const;

private:
    struct Slot {
        std::atomic<bool> ready{false};
        std::atomic<bool> taken{false};
        char* buffer;
        size_t capacity;
        size_t data_size;
        uint64_t timestamp;
        std::string source;
    };

    std::vector<Slot> slots_;
    std::atomic<size_t> head_{0};
    std::atomic<size_t> tail_{0};
    const size_t capacity_;
};

} // namespace log_receiver