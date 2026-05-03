const express = require('express');
const router = express.Router();
const { getGrpcClient } = require('../grpc/client');

router.get('/', async (req, res) => {
  try {
    const client = await getGrpcClient();
    client.GetMetrics({}, (err, response) => {
      if (err) {
        console.error('gRPC error:', err);
        return res.status(500).json({ error: 'Failed to fetch metrics' });
      }
      res.json({
        totalCount: response.totalCount,
        levelCounts: response.levelCounts,
        sourceCounts: response.sourceCounts,
        errorRate: response.errorRate,
        avgLatency: response.avgLatency,
        timestamp: response.timestamp
      });
    });
  } catch (err) {
    console.error('Error:', err);
    res.status(500).json({ error: 'Internal server error' });
  }
});

module.exports = router;