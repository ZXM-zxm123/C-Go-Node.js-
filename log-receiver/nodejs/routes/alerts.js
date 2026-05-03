const express = require('express');
const router = express.Router();
const { getGrpcClient } = require('../grpc/client');

let alertRules = [
  {
    id: 'error_high',
    name: 'High Error Rate',
    level: 'ERROR',
    count: 5,
    windowSeconds: 300,
    consecutive: 2,
    cooldownSeconds: 600,
    webhookUrl: 'http://localhost:3000/webhook/alert',
    active: true
  },
  {
    id: 'warning_high',
    name: 'High Warning Rate',
    level: 'WARNING',
    count: 10,
    windowSeconds: 300,
    consecutive: 3,
    cooldownSeconds: 900,
    webhookUrl: 'http://localhost:3000/webhook/alert',
    active: true
  }
];

let alertHistory = [];

router.get('/', async (req, res) => {
  try {
    const client = await getGrpcClient();
    client.GetAlerts({}, (err, response) => {
      if (err) {
        console.error('gRPC error:', err);
        return res.json({ rules: alertRules });
      }
      res.json({ rules: response.rules });
    });
  } catch (err) {
    res.json({ rules: alertRules });
  }
});

router.get('/history', (req, res) => {
  res.json({ history: alertHistory });
});

router.post('/', (req, res) => {
  const newRule = {
    id: req.body.id || `rule_${Date.now()}`,
    name: req.body.name,
    level: req.body.level || 'ERROR',
    count: req.body.count || 5,
    windowSeconds: req.body.windowSeconds || 300,
    consecutive: req.body.consecutive || 2,
    cooldownSeconds: req.body.cooldownSeconds || 600,
    webhookUrl: req.body.webhookUrl || 'http://localhost:3000/webhook/alert',
    active: req.body.active !== undefined ? req.body.active : true
  };
  
  alertRules.push(newRule);
  res.status(201).json(newRule);
});

router.put('/:id', (req, res) => {
  const index = alertRules.findIndex(r => r.id === req.params.id);
  if (index === -1) {
    return res.status(404).json({ error: 'Rule not found' });
  }
  
  alertRules[index] = { ...alertRules[index], ...req.body };
  res.json(alertRules[index]);
});

router.delete('/:id', (req, res) => {
  const index = alertRules.findIndex(r => r.id === req.params.id);
  if (index === -1) {
    return res.status(404).json({ error: 'Rule not found' });
  }
  
  alertRules.splice(index, 1);
  res.status(204).send();
});

function addAlert(alert) {
  alertHistory.unshift({
    ...alert,
    receivedAt: new Date().toISOString()
  });
  
  if (alertHistory.length > 100) {
    alertHistory.pop();
  }
}

module.exports = { router, addAlert };
