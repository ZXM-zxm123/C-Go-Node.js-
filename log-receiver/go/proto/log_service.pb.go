package proto

import (
    context "context"
    fmt "fmt"
    proto "github.com/gogo/protobuf/proto"
    grpc "google.golang.org/grpc"
    codes "google.golang.org/grpc/codes"
    status "google.golang.org/grpc/status"
    math "math"
)

type MetricsRequest struct {
}

func (m *MetricsRequest) Reset()         { *m = MetricsRequest{} }
func (m *MetricsRequest) String() string { return proto.CompactTextString(m) }
func (*MetricsRequest) ProtoMessage()    {}
func (*MetricsRequest) Descriptor() ([]byte, []int) {
    return fileDescriptor_LogService_87cfc07916b9c11f, []int{0}
}

type MetricsResponse struct {
    TotalCount   int64
    LevelCounts  map[string]int64
    SourceCounts map[string]int64
    ErrorRate    float64
    AvgLatency   float64
    Timestamp    string
}

func (m *MetricsResponse) Reset()         { *m = MetricsResponse{} }
func (m *MetricsResponse) String() string { return proto.CompactTextString(m) }
func (*MetricsResponse) ProtoMessage()    {}
func (*MetricsResponse) Descriptor() ([]byte, []int) {
    return fileDescriptor_LogService_87cfc07916b9c11f, []int{1}
}

type AlertsRequest struct {
}

func (m *AlertsRequest) Reset()         { *m = AlertsRequest{} }
func (m *AlertsRequest) String() string { return proto.CompactTextString(m) }
func (*AlertsRequest) ProtoMessage()    {}
func (*AlertsRequest) Descriptor() ([]byte, []int) {
    return fileDescriptor_LogService_87cfc07916b9c11f, []int{2}
}

type AlertRule struct {
    Id         string
    Name       string
    Level      string
    Count      int32
    Window     int32
    WebhookUrl string
    Active     bool
}

func (m *AlertRule) Reset()         { *m = AlertRule{} }
func (m *AlertRule) String() string { return proto.CompactTextString(m) }
func (*AlertRule) ProtoMessage()    {}
func (*AlertRule) Descriptor() ([]byte, []int) {
    return fileDescriptor_LogService_87cfc07916b9c11f, []int{3}
}

type AlertsResponse struct {
    Rules []*AlertRule
}

func (m *AlertsResponse) Reset()         { *m = AlertsResponse{} }
func (m *AlertsResponse) String() string { return proto.CompactTextString(m) }
func (*AlertsResponse) ProtoMessage()    {}
func (*AlertsResponse) Descriptor() ([]byte, []int) {
    return fileDescriptor_LogService_87cfc07916b9c11f, []int{4}
}

func init() {
    proto.RegisterType((*MetricsRequest)(nil), "logservice.MetricsRequest")
    proto.RegisterType((*MetricsResponse)(nil), "logservice.MetricsResponse")
    proto.RegisterType((*AlertsRequest)(nil), "logservice.AlertsRequest")
    proto.RegisterType((*AlertRule)(nil), "logservice.AlertRule")
    proto.RegisterType((*AlertsResponse)(nil), "logservice.AlertsResponse")
}

const _LogService_serviceDesc = "\n\027LogService\022\036\n\013GetMetrics\022\031.logservice.MetricsRequest\032\034.logservice.MetricsResponse\"\000\022\036\n\011GetAlerts\022\031.logservice.AlertsRequest\032\034.logservice.AlertsResponse\"\000\032\000"

type LogServiceClient interface {
    GetMetrics(ctx context.Context, in *MetricsRequest, opts ...grpc.CallOption) (*MetricsResponse, error)
    GetAlerts(ctx context.Context, in *AlertsRequest, opts ...grpc.CallOption) (*AlertsResponse, error)
}

type logServiceClient struct {
    cc *grpc.ClientConn
}

func NewLogServiceClient(cc *grpc.ClientConn) LogServiceClient {
    return &logServiceClient{cc}
}

func (c *logServiceClient) GetMetrics(ctx context.Context, in *MetricsRequest, opts ...grpc.CallOption) (*MetricsResponse, error) {
    out := new(MetricsResponse)
    err := c.cc.Invoke(ctx, "/logservice.LogService/GetMetrics", in, out, opts...)
    if err != nil {
        return nil, err
    }
    return out, nil
}

func (c *logServiceClient) GetAlerts(ctx context.Context, in *AlertsRequest, opts ...grpc.CallOption) (*AlertsResponse, error) {
    out := new(AlertsResponse)
    err := c.cc.Invoke(ctx, "/logservice.LogService/GetAlerts", in, out, opts...)
    if err != nil {
        return nil, err
    }
    return out, nil
}

type LogServiceServer interface {
    GetMetrics(context.Context, *MetricsRequest) (*MetricsResponse, error)
    GetAlerts(context.Context, *AlertsRequest) (*AlertsResponse, error)
}

func RegisterLogServiceServer(s *grpc.Server, srv LogServiceServer) {
    s.RegisterService(&_LogService_serviceDesc, srv)
}

func _LogService_GetMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(MetricsRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(LogServiceServer).GetMetrics(ctx, in)
    }
    info := &grpc.UnaryServerInfo{
        Server:     srv,
        FullMethod: "/logservice.LogService/GetMetrics",
    }
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(LogServiceServer).GetMetrics(ctx, req.(*MetricsRequest))
    }
    return interceptor(ctx, in, info, handler)
}

func _LogService_GetAlerts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(AlertsRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(LogServiceServer).GetAlerts(ctx, in)
    }
    info := &grpc.UnaryServerInfo{
        Server:     srv,
        FullMethod: "/logservice.LogService/GetAlerts",
    }
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(LogServiceServer).GetAlerts(ctx, req.(*AlertsRequest))
    }
    return interceptor(ctx, in, info, handler)
}

type UnimplementedLogServiceServer struct {
}

func (*UnimplementedLogServiceServer) GetMetrics(context.Context, *MetricsRequest) (*MetricsResponse, error) {
    return nil, status.Errorf(codes.Unimplemented, "method GetMetrics not implemented")
}

func (*UnimplementedLogServiceServer) GetAlerts(context.Context, *AlertsRequest) (*AlertsResponse, error) {
    return nil, status.Errorf(codes.Unimplemented, "method GetAlerts not implemented")
}

var fileDescriptor_LogService_87cfc07916b9c11f = []byte{
    0x0a, 0x16, 0x6c, 0x6f, 0x67, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
    0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22,
    0x15, 0x0a, 0x0e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
    0x74, 0x22, 0x95, 0x01, 0x0a, 0x0f, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73,
    0x70, 0x6f, 0x6e, 0x73, 0x65, 0x0a, 0x0d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75,
    0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43,
    0x6f, 0x75, 0x6e, 0x74, 0x0a, 0x0e, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e,
    0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
    0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74,
    0x72, 0x79, 0x52, 0x0d, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x0a,
    0x11, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x03,
    0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
    0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x10,
    0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x0a, 0x0a, 0x65, 0x72,
    0x72, 0x6f, 0x72, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09,
    0x65, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x61, 0x74, 0x65, 0x0a, 0x0d, 0x61, 0x76, 0x67, 0x5f, 0x6c,
    0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x61, 0x76,
    0x67, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x0a, 0x0d, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
    0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
    0x74, 0x61, 0x6d, 0x70, 0x22, 0x16, 0x0a, 0x0e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x52, 0x65,
    0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x95, 0x01, 0x0a, 0x09, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x52,
    0x75, 0x6c, 0x65, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
    0x64, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
    0x61, 0x6d, 0x65, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
    0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04,
    0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x0a, 0x06, 0x77, 0x69, 0x6e,
    0x64, 0x6f, 0x77, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x77, 0x69, 0x6e, 0x64, 0x6f,
    0x77, 0x0a, 0x0b, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x06,
    0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x55, 0x72, 0x6c,
    0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06,
    0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x22, 0x36, 0x0a, 0x0e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x73,
    0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18,
    0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69,
    0x63, 0x65, 0x2e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x75,
    0x6c, 0x65, 0x73, 0x32, 0xa5, 0x01, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69,
    0x63, 0x65, 0x12, 0x36, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
    0x12, 0x31, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x65,
    0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x34, 0x2e, 0x6c,
    0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
    0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x0b, 0x47,
    0x65, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x12, 0x31, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65,
    0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
    0x65, 0x73, 0x74, 0x1a, 0x34, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
    0x2e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
    0x00, 0x3a, 0x00, 0x42, 0x10, 0x67, 0x6f, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x3d, 27,
    0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x7d,
}

func (m *MetricsResponse) Marshal() (dAtA []byte, err error) {
    size := m.Size()
    dAtA = make([]byte, size)
    n, err := m.MarshalTo(dAtA)
    if err != nil {
        return nil, err
    }
    return dAtA[:n], nil
}

func (m *MetricsResponse) MarshalTo(dAtA []byte) (int, error) {
    var i int
    _ = i
    var l int
    _ = l
    if m.TotalCount != 0 {
        dAtA[i] = 0x08
        i++
        i = encodeVarintLogService(dAtA, i, uint64(m.TotalCount))
    }
    if len(m.LevelCounts) > 0 {
        for kv := range m.LevelCounts {
            size := len(kv)
            dAtA[i] = 0x12
            i++
            i = encodeVarintLogService(dAtA, i, uint64(size+8))
            dAtA[i] = 0x0a
            i++
            i = encodeVarintLogService(dAtA, i, uint64(len(kv)))
            i += copy(dAtA[i:], kv)
            dAtA[i] = 0x10
            i++
            i = encodeVarintLogService(dAtA, i, uint64(m.LevelCounts[kv]))
        }
    }
    if len(m.SourceCounts) > 0 {
        for kv := range m.SourceCounts {
            size := len(kv)
            dAtA[i] = 0x1a
            i++
            i = encodeVarintLogService(dAtA, i, uint64(size+8))
            dAtA[i] = 0x0a
            i++
            i = encodeVarintLogService(dAtA, i, uint64(len(kv)))
            i += copy(dAtA[i:], kv)
            dAtA[i] = 0x10
            i++
            i = encodeVarintLogService(dAtA, i, uint64(m.SourceCounts[kv]))
        }
    }
    if m.ErrorRate != 0 {
        dAtA[i] = 0x25
        i++
        bits := math.Float64bits(m.ErrorRate)
        dAtA[i] = uint8(bits)
        i++
        dAtA[i] = uint8(bits >> 8)
        i++
        dAtA[i] = uint8(bits >> 16)
        i++
        dAtA[i] = uint8(bits >> 24)
        i++
        dAtA[i] = uint8(bits >> 32)
        i++
        dAtA[i] = uint8(bits >> 40)
        i++
        dAtA[i] = uint8(bits >> 48)
        i++
        dAtA[i] = uint8(bits >> 56)
        i++
    }
    if m.AvgLatency != 0 {
        dAtA[i] = 0x2d
        i++
        bits := math.Float64bits(m.AvgLatency)
        dAtA[i] = uint8(bits)
        i++
        dAtA[i] = uint8(bits >> 8)
        i++
        dAtA[i] = uint8(bits >> 16)
        i++
        dAtA[i] = uint8(bits >> 24)
        i++
        dAtA[i] = uint8(bits >> 32)
        i++
        dAtA[i] = uint8(bits >> 40)
        i++
        dAtA[i] = uint8(bits >> 48)
        i++
        dAtA[i] = uint8(bits >> 56)
        i++
    }
    if len(m.Timestamp) > 0 {
        dAtA[i] = 0x32
        i++
        i = encodeVarintLogService(dAtA, i, uint64(len(m.Timestamp)))
        i += copy(dAtA[i:], m.Timestamp)
    }
    return i, nil
}

func encodeVarintLogService(dAtA []byte, offset int, v uint64) int {
    for v >= 1<<7 {
        dAtA[offset] = uint8(v&0x7f | 0x80)
        v >>= 7
        offset++
    }
    dAtA[offset] = uint8(v)
    return offset + 1
}

func (m *MetricsResponse) Size() (n int) {
    var l int
    _ = l
    if m.TotalCount != 0 {
        n += 1 + sovLogService(uint64(m.TotalCount))
    }
    if len(m.LevelCounts) > 0 {
        for kv := range m.LevelCounts {
            l = len(kv)
            n += 1 + sovLogService(uint64(l+8)) + 1 + sovLogService(uint64(l)) + l + 1 + sovLogService(uint64(m.LevelCounts[kv]))
        }
    }
    if len(m.SourceCounts) > 0 {
        for kv := range m.SourceCounts {
            l = len(kv)
            n += 1 + sovLogService(uint64(l+8)) + 1 + sovLogService(uint64(l)) + l + 1 + sovLogService(uint64(m.SourceCounts[kv]))
        }
    }
    if m.ErrorRate != 0 {
        n += 9
    }
    if m.AvgLatency != 0 {
        n += 9
    }
    l = len(m.Timestamp)
    if l > 0 {
        n += 1 + sovLogService(uint64(l)) + l
    }
    return n
}

func sovLogService(x uint64) (n int) {
    for {
        n++
        x >>= 7
        if x == 0 {
            break
        }
    }
    return n
}

func (m *AlertRule) Marshal() (dAtA []byte, err error) {
    size := m.Size()
    dAtA = make([]byte, size)
    n, err := m.MarshalTo(dAtA)
    if err != nil {
        return nil, err
    }
    return dAtA[:n], nil
}

func (m *AlertRule) MarshalTo(dAtA []byte) (int, error) {
    var i int
    _ = i
    var l int
    _ = l
    if len(m.Id) > 0 {
        dAtA[i] = 0x0a
        i++
        i = encodeVarintLogService(dAtA, i, uint64(len(m.Id)))
        i += copy(dAtA[i:], m.Id)
    }
    if len(m.Name) > 0 {
        dAtA[i] = 0x12
        i++
        i = encodeVarintLogService(dAtA, i, uint64(len(m.Name)))
        i += copy(dAtA[i:], m.Name)
    }
    if len(m.Level) > 0 {
        dAtA[i] = 0x1a
        i++
        i = encodeVarintLogService(dAtA, i, uint64(len(m.Level)))
        i += copy(dAtA[i:], m.Level)
    }
    if m.Count != 0 {
        dAtA[i] = 0x20
        i++
        dAtA[i] = uint8(m.Count)
        i++
        dAtA[i] = uint8(m.Count >> 8)
        i++
        dAtA[i] = uint8(m.Count >> 16)
        i++
        dAtA[i] = uint8(m.Count >> 24)
        i++
    }
    if m.Window != 0 {
        dAtA[i] = 0x28
        i++
        dAtA[i] = uint8(m.Window)
        i++
        dAtA[i] = uint8(m.Window >> 8)
        i++
        dAtA[i] = uint8(m.Window >> 16)
        i++
        dAtA[i] = uint8(m.Window >> 24)
        i++
    }
    if len(m.WebhookUrl) > 0 {
        dAtA[i] = 0x32
        i++
        i = encodeVarintLogService(dAtA, i, uint64(len(m.WebhookUrl)))
        i += copy(dAtA[i:], m.WebhookUrl)
    }
    if m.Active {
        dAtA[i] = 0x38
        i++
        dAtA[i] = 1
        i++
    }
    return i, nil
}

func (m *AlertRule) Size() (n int) {
    var l int
    _ = l
    l = len(m.Id)
    if l > 0 {
        n += 1 + sovLogService(uint64(l)) + l
    }
    l = len(m.Name)
    if l > 0 {
        n += 1 + sovLogService(uint64(l)) + l
    }
    l = len(m.Level)
    if l > 0 {
        n += 1 + sovLogService(uint64(l)) + l
    }
    if m.Count != 0 {
        n += 5
    }
    if m.Window != 0 {
        n += 5
    }
    l = len(m.WebhookUrl)
    if l > 0 {
        n += 1 + sovLogService(uint64(l)) + l
    }
    if m.Active {
        n += 2
    }
    return n
}

func (m *AlertsResponse) Marshal() (dAtA []byte, err error) {
    size := m.Size()
    dAtA = make([]byte, size)
    n, err := m.MarshalTo(dAtA)
    if err != nil {
        return nil, err
    }
    return dAtA[:n], nil
}

func (m *AlertsResponse) MarshalTo(dAtA []byte) (int, error) {
    var i int
    _ = i
    var l int
    _ = l
    for _, msg := range m.Rules {
        size := msg.Size()
        dAtA[i] = 0x0a
        i++
        i = encodeVarintLogService(dAtA, i, uint64(size))
        n, err := msg.MarshalTo(dAtA[i:])
        if err != nil {
            return 0, err
        }
        i += n
    }
    return i, nil
}

func (m *AlertsResponse) Size() (n int) {
    var l int
    _ = l
    for _, msg := range m.Rules {
        l = msg.Size()
        n += 1 + sovLogService(uint64(l)) + l
    }
    return n
}