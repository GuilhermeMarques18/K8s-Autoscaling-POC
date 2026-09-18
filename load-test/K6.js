import http from 'k6/http';
import { check } from 'k6';
import { Trend, Rate } from 'k6/metrics';

const BASE_URL = __ENV.BASE_URL || 'http://kitchen:9000';

const getOrdersDuration = new Trend('get_orders_duration', true);
const createOrderDuration = new Trend('create_order_duration', true);
const errorRate = new Rate('errors');

export const options = {
  scenarios: {
    carga_progressiva: {
      executor: 'ramping-arrival-rate',
      startRate: 50,
      timeUnit: '1s',
      preAllocatedVUs: 300,
      maxVUs: 800,
      stages: [
        { duration: '1m', target: 100 },
        { duration: '2m', target: 200 },
        { duration: '1m', target: 400 },
        { duration: '3m', target: 600 },
        { duration: '1m', target: 0 },
      ],
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<1000'],
    errors: ['rate<0.05'],
  },
};

export default function () {
  const customerId = Math.floor(Math.random() * 100) + 1;

  if (Math.random() < 0.7) {
    const res = http.get(`${BASE_URL}/?customerId=${customerId}`);
    getOrdersDuration.add(res.timings.duration);
    const ok = check(res, {
      'GET /: status é 200': (r) => r.status === 200,
    });
    errorRate.add(!ok);
  } else {
    const payload = JSON.stringify({
      customerId: customerId,
      productId: Math.floor(Math.random() * 20) + 1,
      quantity: Math.floor(Math.random() * 5) + 1,
    });

    const res = http.post(`${BASE_URL}/orders`, payload, {
      headers: { 'Content-Type': 'application/json' },
    });
    createOrderDuration.add(res.timings.duration);
    const ok = check(res, {
      'POST /orders: status é 201': (r) => r.status === 201,
    });
    errorRate.add(!ok);
  }
}