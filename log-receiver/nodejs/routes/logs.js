const express = require('express');
const router = express.Router();
const { redisClient } = require('../main');

router.get('/', async (req, res) => {
  const count = parseInt(req.query.count) || 100;
  const stream = req.query.stream || 'log_stream';
  
  try {
    const result = await redisClient.xrevrange(stream, '+', '-', 'COUNT', count);
    
    const logs = result.map(([id, fields]) => {
      const obj = { id };
      for (let i = 0; i < fields.length; i += 2) {
        obj[fields[i]] = fields[i + 1];
      }
      return obj;
    });
    
    res.json({ logs, count: logs.length });
  } catch (err) {
    console.error('Redis error:', err);
    res.status(500).json({ error: 'Failed to fetch logs' });
  }
});

router.get('/search', async (req, res) => {
  const query = req.query.q || '';
  const count = parseInt(req.query.count) || 50;
  const stream = req.query.stream || 'log_stream';
  
  try {
    const result = await redisClient.xrevrange(stream, '+', '-', 'COUNT', count * 2);
    
    const logs = result
      .map(([id, fields]) => {
        const obj = { id };
        for (let i = 0; i < fields.length; i += 2) {
          obj[fields[i]] = fields[i + 1];
        }
        return obj;
      })
      .filter(log => {
        if (!query) return true;
        return JSON.stringify(log).toLowerCase().includes(query.toLowerCase());
      })
      .slice(0, count);
    
    res.json({ logs, count: logs.length });
  } catch (err) {
    console.error('Redis error:', err);
    res.status(500).json({ error: 'Failed to search logs' });
  }
});

module.exports = router;