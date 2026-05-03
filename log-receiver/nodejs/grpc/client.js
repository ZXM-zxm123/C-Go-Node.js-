const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

let client = null;
let config = null;

function setConfig(cfg) {
    config = cfg;
}

async function getGrpcClient() {
    if (client) {
        return client;
    }

    if (!config) {
        const { loadConfig } = require('../config');
        config = loadConfig();
    }

    const PROTO_PATH = path.join(__dirname, 'log_service.proto');

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

    const grpcAddr = `${config.node_api.grpc_host}:${config.node_api.grpc_port}`;
    client = new logService.LogService(
        grpcAddr,
        grpc.credentials.createInsecure()
    );

    return client;
}

module.exports = { getGrpcClient, setConfig };