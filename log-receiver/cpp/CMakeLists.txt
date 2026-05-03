cmake_minimum_required(VERSION 3.10)
project(log_receiver)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -O3 -Wall -Wextra -pthread")

find_package(hiredis REQUIRED)
find_package(Threads REQUIRED)

add_executable(log_receiver
    src/main.cpp
    src/log_receiver.cpp
    src/zero_copy_queue.cpp
    src/redis_publisher.cpp
    src/config.cpp
)

target_include_directories(log_receiver PRIVATE include)
target_link_libraries(log_receiver PRIVATE hiredis::hiredis Threads::Threads)