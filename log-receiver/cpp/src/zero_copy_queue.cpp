#include "zero_copy_queue.h"
#include <sys/mman.h>
#include <chrono>

namespace log_receiver {

ZeroCopyQueue::ZeroCopyQueue(size_t capacity) : capacity_(capacity) {
    slots_.resize(capacity);
    for (auto& slot : slots_) {
        slot.buffer = static_cast<char*>(mmap(nullptr, 65536, 
            PROT_READ | PROT_WRITE, MAP_PRIVATE | MAP_ANONYMOUS, -1, 0));
        slot.capacity = 65536;
    }
}

ZeroCopyQueue::~ZeroCopyQueue() {
    for (auto& slot : slots_) {
        if (slot.buffer) {
            munmap(slot.buffer, slot.capacity);
        }
    }
}

bool ZeroCopyQueue::push(const char* data, size_t size, const std::string& source) {
    if (size == 0 || size > 65536) return false;
    
    size_t current_head = head_.load(std::memory_order_relaxed);
    Slot& slot = slots_[current_head];
    
    if (!slot.taken.compare_exchange_strong(slot.taken, true)) {
        return false;
    }
    
    std::memcpy(slot.buffer, data, size);
    slot.data_size = size;
    slot.timestamp = std::chrono::system_clock::now().time_since_epoch().count();
    slot.source = source;
    
    slot.ready.store(true, std::memory_order_release);
    head_.store((current_head + 1) % capacity_, std::memory_order_relaxed);
    
    return true;
}

bool ZeroCopyQueue::pop(Chunk& chunk) {
    size_t current_tail = tail_.load(std::memory_order_relaxed);
    Slot& slot = slots_[current_tail];
    
    if (!slot.ready.load(std::memory_order_acquire)) {
        return false;
    }
    
    chunk.data = slot.buffer;
    chunk.size = slot.data_size;
    chunk.capacity = slot.capacity;
    chunk.timestamp = slot.timestamp;
    chunk.source = slot.source;
    
    slot.ready.store(false, std::memory_order_release);
    slot.taken.store(false, std::memory_order_release);
    tail_.store((current_tail + 1) % capacity_, std::memory_order_relaxed);
    
    return true;
}

size_t ZeroCopyQueue::size() const {
    size_t h = head_.load(std::memory_order_relaxed);
    size_t t = tail_.load(std::memory_order_relaxed);
    return (h >= t) ? (h - t) : (capacity_ - t + h);
}

bool ZeroCopyQueue::empty() const {
    return head_.load(std::memory_order_relaxed) == tail_.load(std::memory_order_relaxed);
}

} // namespace log_receiver