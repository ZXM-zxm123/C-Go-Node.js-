const express = require('express');
const router = express.Router();
const { addAlert } = require('./alerts');

router.post('/alert', (req, res) => {
  const alert = req.body;
  
  console.log('Received alert:', JSON.stringify(alert, null, 2));
  
  addAlert(alert);
  
  res.status(200).json({ status: 'ok', received: true });
});

module.exports = router;