const fs = require('fs');
const path = require('path');
const yaml = require('js-yaml');

function getConfigPath() {
    const envPath = process.env.CONFIG_PATH;
    if (envPath) {
        return envPath;
    }
    return path.join(__dirname, '..', 'config', 'config.yaml');
}

function formatYAMLError(filePath, error) {
    let errType = 'configuration parse error';
    let lineNum = 0;

    if (error.mark) {
        lineNum = error.mark.line + 1;
        if (error.reason) {
            if (error.reason.includes('unknown')) {
                errType = 'unknown configuration field';
            } else if (error.reason.includes('expected')) {
                errType = 'YAML syntax error';
            } else if (error.reason.includes('mismatch')) {
                errType = 'type mismatch error';
            }
        }
    }

    return {
        type: errType,
        file: filePath,
        line: lineNum,
        original: error,
        message: `${errType} in ${filePath}${lineNum > 0 ? ` (line ${lineNum})` : ''}: ${error.message}`
    };
}

function validateConfig(cfg) {
    const errors = [];

    if (!cfg.node_api) {
        errors.push('validation error: node_api section is required');
    } else {
        const api = cfg.node_api;
        if (!api.http_port || typeof api.http_port !== 'number') {
            errors.push('validation error: node_api.http_port is required and must be a number');
        } else if (api.http_port < 1 || api.http_port > 65535) {
            errors.push(`validation error: node_api.http_port must be between 1 and 65535, got ${api.http_port}`);
        }
        if (!api.grpc_host || typeof api.grpc_host !== 'string') {
            errors.push('validation error: node_api.grpc_host is required');
        }
        if (!api.grpc_port || typeof api.grpc_port !== 'number') {
            errors.push('validation error: node_api.grpc_port is required and must be a number');
        } else if (api.grpc_port < 1 || api.grpc_port > 65535) {
            errors.push(`validation error: node_api.grpc_port must be between 1 and 65535, got ${api.grpc_port}`);
        }
        if (!api.redis_host || typeof api.redis_host !== 'string') {
            errors.push('validation error: node_api.redis_host is required');
        }
        if (!api.redis_port || typeof api.redis_port !== 'number') {
            errors.push('validation error: node_api.redis_port is required and must be a number');
        } else if (api.redis_port < 1 || api.redis_port > 65535) {
            errors.push(`validation error: node_api.redis_port must be between 1 and 65535, got ${api.redis_port}`);
        }
    }

    return errors;
}

function loadConfig() {
    const configPath = getConfigPath();

    try {
        const stat = fs.statSync(configPath);
        if (!stat.isFile()) {
            console.error(`ERROR: configuration path is not a file: ${configPath}`);
            process.exit(1);
        }
    } catch (error) {
        if (error.code === 'ENOENT') {
            console.error(`ERROR: configuration file not found: ${configPath}`);
        } else {
            console.error(`ERROR: failed to stat config file ${configPath}: ${error.message}`);
        }
        process.exit(1);
    }

    let data;
    try {
        data = fs.readFileSync(configPath, 'utf8');
    } catch (error) {
        console.error(`ERROR: failed to read config file ${configPath}: ${error.message}`);
        process.exit(1);
    }

    let config;
    try {
        config = yaml.load(data, {
            json: true,
            schema: yaml.FAILSAFE_SCHEMA
        });
    } catch (error) {
        const formatted = formatYAMLError(configPath, error);
        console.error(`ERROR: ${formatted.message}`);
        process.exit(1);
    }

    const validationErrors = validateConfig(config);
    if (validationErrors.length > 0) {
        console.error('ERROR: configuration validation failed:');
        validationErrors.forEach(err => console.error(`  - ${err}`));
        process.exit(1);
    }

    return config;
}

module.exports = { loadConfig, getConfigPath };