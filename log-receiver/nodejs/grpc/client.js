const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');

const PROTO_PATH = '../proto/log_service.proto';

let client = null;

async function getGrpcClient() {
  if (client) {
    return client;
  }

  const packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true
    }
  );

  const logService = grpc.loadPackageDefinition(packageDefinition).logservice;
  
  client = new logService.LogService(
    'localhost:50051',
    grpc.credentials.createInsecure()
  );

  return client;
}

module.exports = { getGrpcClient };