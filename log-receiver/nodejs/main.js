const express = require('express');
const cors = require('cors');
const redis = require('ioredis');

const metricsRouter = require('./routes/metrics');
const alertsRouter = require('./routes/alerts');
const logsRouter = require('./routes/logs');
const webhookRouter = require('./routes/webhook');

const app = express();
const PORT = process.env.PORT || 3000;

const redisClient = new redis({
  host: 'localhost',
  port: 6379
});

app.use(cors());
app.use(express.json());

app.use('/api/metrics', metricsRouter);
app.use('/api/alerts', alertsRouter);
app.use('/api/logs', logsRouter);
app.use('/webhook', webhookRouter);

app.get('/', (req, res) => {
  res.json({ status: 'ok', service: 'log-api' });
});

redisClient.on('error', (err) => {
  console.error('Redis connection error:', err);
});

redisClient.on('connect', () => {
  console.log('Connected to Redis');
});

app.listen(PORT, () => {
  console.log(`Log API server running on port ${PORT}`);
});

module.exports = { redisClient };