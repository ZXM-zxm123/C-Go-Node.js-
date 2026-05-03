#include "redis_publisher.h"
#include <hiredis/hiredis.h>
#include <iostream>

namespace log_receiver {

struct RedisPublisher::Impl {
    std::string host;
    int port;
    std::string stream;
    redisContext* ctx{nullptr};
};

RedisPublisher::RedisPublisher(const std::string& host, int port, const std::string& stream)
    : impl_(std::make_unique<Impl>()) {
    impl_->host = host;
    impl_->port = port;
    impl_->stream = stream;
}

RedisPublisher::~RedisPublisher() {
    disconnect();
}

bool RedisPublisher::connect() {
    impl_->ctx = redisConnect(impl_->host.c_str(), impl_->port);
    if (!impl_->ctx || impl_->ctx->err) {
        if (impl_->ctx) {
            std::cerr << "Redis connection error: " << impl_->ctx->errstr << std::endl;
            redisFree(impl_->ctx);
            impl_->ctx = nullptr;
        } else {
            std::cerr << "Redis connection error: cannot allocate context" << std::endl;
        }
        return false;
    }
    return true;
}

void RedisPublisher::disconnect() {
    if (impl_->ctx) {
        redisFree(impl_->ctx);
        impl_->ctx = nullptr;
    }
}

bool RedisPublisher::publish(const std::string& data, const std::string& source) {
    if (!impl_->ctx) {
        if (!connect()) return false;
    }

    redisReply* reply = static_cast<redisReply*>(redisCommand(
        impl_->ctx, "XADD %s * data %s source %s",
        impl_->stream.c_str(),
        data.c_str(),
        source.c_str()
    ));

    if (!reply) {
        disconnect();
        return false;
    }

    freeReplyObject(reply);
    return true;
}

} // namespace log_receiver